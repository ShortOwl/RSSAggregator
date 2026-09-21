import { it, expect, vi } from "vitest";
import { collectPostIds } from "../src/hooks/use-post-state";
import { api } from "../src/lib/api";
vi.mock("../src/lib/api", async (importOriginal) => ({
  ...(await importOriginal<typeof import("../src/lib/api")>()),
  api: vi.fn(),
}));
it("reconciles membership across every cursor page, deduplicating posts", async () => {
  vi.mocked(api)
    .mockResolvedValueOnce({ data: [{ id: "one" }], next_cursor: "next-page" })
    .mockResolvedValueOnce({
      data: [{ id: "one" }, { id: "two" }],
      next_cursor: "",
    });
  await expect(collectPostIds("/posts?unread=true")).resolves.toEqual([
    "one",
    "two",
  ]);
  expect(api).toHaveBeenLastCalledWith(
    "/posts?unread=true&limit=100&cursor=next-page",
    { signal: undefined },
  );
});
it("fails instead of silently trusting a repeated cursor", async () => {
  vi.mocked(api).mockResolvedValue({ data: [], next_cursor: "same" });
  await expect(collectPostIds("/bookmarks")).rejects.toThrow("repeated page");
});
