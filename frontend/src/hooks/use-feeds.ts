"use client";
import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { api, queryString } from "@/lib/api";
import { keys } from "@/lib/query-keys";
import type { Feed, FeedFollow, Page } from "@/lib/types";
export function useFeeds(search = "") {
  return useInfiniteQuery({
    queryKey: [...keys.feeds, search],
    initialPageParam: "",
    queryFn: ({ pageParam, signal }) =>
      api<Page<Feed>>(
        `/feeds?${queryString({ search, cursor: pageParam, limit: 100 })}`,
        { signal },
      ),
    getNextPageParam: (page) => page.next_cursor || undefined,
  });
}
export function useFollows() {
  return useQuery({
    queryKey: keys.follows,
    queryFn: ({ signal }) => api<FeedFollow[]>("/feed_follows", { signal }),
  });
}
