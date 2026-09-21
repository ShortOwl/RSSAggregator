"use client";
import { useInfiniteQuery } from "@tanstack/react-query";
import { api, queryString } from "@/lib/api";
import type { Page, Post, PostFilters } from "@/lib/types";
export function usePosts(mode: "posts" | "bookmarks", filters: PostFilters) {
  return useInfiniteQuery({
    queryKey: [mode, filters],
    initialPageParam: "",
    queryFn: ({ pageParam, signal }) =>
      api<Page<Post>>(
        `/${mode}?${queryString({ ...(mode === "posts" ? filters : {}), limit: 20, cursor: pageParam })}`,
        { signal },
      ),
    getNextPageParam: (page) => page.next_cursor || undefined,
  });
}
