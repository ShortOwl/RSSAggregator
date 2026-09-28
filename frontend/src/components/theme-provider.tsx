"use client";

import {
  createContext,
  useContext,
  useLayoutEffect,
  useSyncExternalStore,
} from "react";
import { parseTheme, themeKey, type Theme } from "@/lib/theme";

let preference: Theme | undefined;
const listeners = new Set<() => void>();

function getSnapshot(): Theme {
  if (preference === undefined) {
    try {
      preference = parseTheme(localStorage.getItem(themeKey));
    } catch {
      preference = "system";
    }
  }
  return preference;
}

function applyTheme(theme: Theme) {
  document.documentElement.dataset.theme =
    theme === "system"
      ? window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light"
      : theme;
}

function notify() {
  listeners.forEach((listener) => listener());
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  function onStorage(event: StorageEvent) {
    if (event.key === themeKey || event.key === null) {
      preference = parseTheme(event.key === null ? null : event.newValue);
      applyTheme(preference);
      notify();
    }
  }
  window.addEventListener("storage", onStorage);
  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", onStorage);
  };
}

function setTheme(theme: Theme) {
  preference = theme;
  applyTheme(theme);
  try {
    localStorage.setItem(themeKey, theme);
  } catch {
    // Keep the preference in memory for this visit when storage is blocked.
  }
  notify();
}

const ThemeContext = createContext({ theme: "system" as Theme, setTheme });

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const theme = useSyncExternalStore(
    subscribe,
    getSnapshot,
    () => "system" as Theme,
  );
  useLayoutEffect(() => {
    // Read the client snapshot, including during hydration and Strict Mode remounts.
    applyTheme(getSnapshot());
    const media = window.matchMedia("(prefers-color-scheme: dark)");
    const update = () => applyTheme(getSnapshot());
    media.addEventListener("change", update);
    return () => media.removeEventListener("change", update);
  }, []);
  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  return useContext(ThemeContext);
}
