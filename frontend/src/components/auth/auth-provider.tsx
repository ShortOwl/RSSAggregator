"use client";
import { createContext, useContext, useEffect, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  clearSession,
  expiresAt,
  getToken,
  SESSION_EVENT,
  sessionKey,
} from "@/lib/session";
const Context = createContext({
  ready: false,
  authenticated: false,
  signOut: () => {},
});
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const client = useQueryClient();
  const [state, setState] = useState({ ready: false, authenticated: false });
  useEffect(() => {
    let timer: ReturnType<typeof setTimeout>;
    let previous: string | null = null;
    const sync = () => {
      clearTimeout(timer);
      const token = getToken();
      if (token !== previous) {
        void client.cancelQueries();
        client.clear();
        previous = token;
      }
      setState({ ready: true, authenticated: !!token });
      if (token)
        timer = setTimeout(
          clearSession,
          Math.max(0, expiresAt(token) - Date.now()),
        );
    };
    sync();
    const storage = (e: StorageEvent) => {
      if (e.key === sessionKey || e.key === null) sync();
    };
    window.addEventListener(SESSION_EVENT, sync);
    window.addEventListener("storage", storage);
    return () => {
      clearTimeout(timer);
      window.removeEventListener(SESSION_EVENT, sync);
      window.removeEventListener("storage", storage);
    };
  }, [client]);
  return (
    <Context.Provider value={{ ...state, signOut: clearSession }}>
      {children}
    </Context.Provider>
  );
}
export function useAuth() {
  return useContext(Context);
}
