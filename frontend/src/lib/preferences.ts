export const preferencesKey = "margin.preferences";
export interface Preferences {
  compact: boolean;
}
export function readPreferences(): Preferences {
  try {
    return {
      compact:
        JSON.parse(localStorage.getItem(preferencesKey) || "{}").compact ===
        true,
    };
  } catch {
    return { compact: false };
  }
}
