"use client";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { api, ApiError, queryString } from "@/lib/api";
import { getToken } from "@/lib/session";
import { keys } from "@/lib/query-keys";
import type { Page, Post } from "@/lib/types";
export async function collectPostIds(
  path: string,
  signal?: AbortSignal,
): Promise<string[]> {
  const ids = new Set<string>();
  const cursors = new Set<string>();
  let cursor = "";
  do {
    signal?.throwIfAborted();
    const url = `${path}${path.includes("?") ? "&" : "?"}${queryString({ limit: 100, cursor })}`;
    let page: Page<Post>;
    for (let attempt = 0; ; attempt++) {
      try {
        page = await api<Page<Post>>(url, { signal });
        break;
      } catch (error) {
        if (
          !(error instanceof ApiError) ||
          error.status !== 429 ||
          attempt >= 3
        )
          throw error;
        // Retry this cursor rather than restarting a large collection.
        await new Promise<void>((resolve) =>
          setTimeout(resolve, error.retryAfter * 1000),
        );
        signal?.throwIfAborted();
      }
    }
    page.data.forEach((post) => ids.add(post.id));
    cursor = page.next_cursor;
    if (cursor && cursors.has(cursor))
      throw new Error("The server returned a repeated page. Please retry.");
    cursors.add(cursor);
    if (cursor) await new Promise<void>((resolve) => setTimeout(resolve, 600));
  } while (cursor);
  return [...ids];
}
export function usePostState() {
  const bookmarks = useQuery({
    queryKey: [...keys.state, "bookmarks"],
    queryFn: ({ signal }) => collectPostIds("/bookmarks", signal),
    staleTime: 60_000,
  });
  const unread = useQuery({
    queryKey: [...keys.state, "unread"],
    queryFn: ({ signal }) => collectPostIds("/posts?unread=true", signal),
    staleTime: 60_000,
  });
  return { bookmarks, unread };
}
export function usePostAction(id: string) {
  const client = useQueryClient();
  return useMutation({
    onMutate: () => ({ token: getToken() }),
    mutationFn: ({
      kind,
      active,
    }: {
      kind: "read" | "bookmark";
      active: boolean;
    }) =>
      api<void>(`/posts/${id}/${kind}`, { method: active ? "PUT" : "DELETE" }),
    onSuccess: async (_, { kind, active }, context) => {
      if (context?.token !== getToken()) return;
      const stateKey = [
        ...keys.state,
        kind === "read" ? "unread" : "bookmarks",
      ];
      await client.cancelQueries({ queryKey: stateKey });
      client.setQueryData<string[]>(
        [...keys.state, kind === "read" ? "unread" : "bookmarks"],
        (old) =>
          old
            ? (kind === "read" ? !active : active)
              ? [...new Set([...old, id])]
              : old.filter((value) => value !== id)
            : old,
      );
      void client.invalidateQueries({ queryKey: stateKey });
      void client.invalidateQueries({ queryKey: keys.posts });
      void client.invalidateQueries({ queryKey: keys.bookmarks });
      toast.success(
        kind === "read"
          ? active
            ? "Marked as read"
            : "Marked as unread"
          : active
            ? "Saved to bookmarks"
            : "Bookmark removed",
      );
    },
    onError: (error) => toast.error(error.message),
  });
}
