import { describe, it, expect, vi, afterEach } from "vitest";
import { api, ApiError, queryString } from "../src/lib/api";
vi.mock("../src/lib/session", () => ({
  getToken: () => "test-jwt",
  clearSession: vi.fn(),
}));
afterEach(() => vi.unstubAllGlobals());
describe("Go API contract", () => {
  it("encodes search and forwards an opaque cursor", () =>
    expect(
      queryString({
        search: "a & b",
        cursor: "abc-_",
        unread: true,
        feed_id: "",
      }),
    ).toBe("search=a+%26+b&cursor=abc-_&unread=true"));
  it("sends Bearer authorization and accepts a 204 mutation response", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValue(new Response(null, { status: 204 }));
    vi.stubGlobal("fetch", fetch);
    await expect(
      api("/posts/1/read", { method: "PUT" }),
    ).resolves.toBeUndefined();
    expect(fetch.mock.calls[0][1].headers.get("Authorization")).toBe(
      "Bearer test-jwt",
    );
  });
  it("preserves backend errors and rate-limit delay", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ error: "rate limit exceeded" }), {
          status: 429,
          headers: { "Retry-After": "3" },
        }),
      ),
    );
    await expect(api("/posts")).rejects.toMatchObject({
      status: 429,
      retryAfter: 3,
      message: "rate limit exceeded",
    });
  });
  it("keeps public login free of stale credentials", async () => {
    const fetch = vi.fn().mockResolvedValue(Response.json({ token: "new" }));
    vi.stubGlobal("fetch", fetch);
    await api("/login", { method: "POST", body: "{}" }, false);
    expect(fetch.mock.calls[0][1].headers.has("Authorization")).toBe(false);
  });
  it("reports network errors", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockRejectedValue(new TypeError("Failed to fetch")),
    );
    await expect(api("/posts")).rejects.toBeInstanceOf(ApiError);
  });
});
