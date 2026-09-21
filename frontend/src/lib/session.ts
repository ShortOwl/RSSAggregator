const KEY = "margin.session";
export const SESSION_EVENT = "margin-session";
export function expiresAt(token: string): number {
  try {
    const part = token.split(".")[1].replace(/-/g, "+").replace(/_/g, "/");
    const payload = JSON.parse(atob(part));
    return typeof payload.exp === "number" ? payload.exp * 1000 : 0;
  } catch {
    return 0;
  }
}
export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  try {
    const token = localStorage.getItem(KEY);
    return token && expiresAt(token) > Date.now() ? token : null;
  } catch {
    return null;
  }
}
export function saveToken(token: string) {
  localStorage.setItem(KEY, token);
  window.dispatchEvent(new Event(SESSION_EVENT));
}
export function clearSession() {
  try {
    localStorage.removeItem(KEY);
  } catch {
    // Storage can be unavailable in restricted browser contexts.
  } finally {
    window.dispatchEvent(new Event(SESSION_EVENT));
  }
}
export const sessionKey = KEY;
