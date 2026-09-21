import { describe, it, expect, vi, afterEach } from "vitest";
import {
  expiresAt,
  getToken,
  saveToken,
  clearSession,
} from "../src/lib/session";
const token = (exp: number) =>
  `header.${btoa(JSON.stringify({ exp }))}.signature`;
afterEach(() => vi.unstubAllGlobals());
describe("persistent session", () => {
  it("rejects malformed and expired tokens", () => {
    expect(expiresAt("bad")).toBe(0);
    vi.stubGlobal("window", {});
    vi.stubGlobal("localStorage", { getItem: () => token(1) });
    expect(getToken()).toBeNull();
  });
  it("stores and clears sessions, notifying listeners", () => {
    const data = new Map();
    const dispatchEvent = vi.fn();
    vi.stubGlobal("window", { dispatchEvent });
    vi.stubGlobal("localStorage", {
      getItem: (key: string) => data.get(key),
      setItem: (key: string, value: string) => data.set(key, value),
      removeItem: (key: string) => data.delete(key),
    });
    const jwt = token(Date.now() / 1000 + 1000);
    saveToken(jwt);
    expect(getToken()).toBe(jwt);
    clearSession();
    expect(getToken()).toBeNull();
    expect(dispatchEvent).toHaveBeenCalledTimes(2);
  });
});
