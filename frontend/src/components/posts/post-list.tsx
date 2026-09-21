"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { ArrowRight, Rss } from "lucide-react";
import { usePosts } from "@/hooks/use-posts";
import { useFeeds, useFollows } from "@/hooks/use-feeds";
import { usePostState } from "@/hooks/use-post-state";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { readPreferences } from "@/lib/preferences";
import type { PostFilters } from "@/lib/types";
import { PageHeader } from "@/components/layout/page-header";
import { ReadingSidebar } from "./reading-sidebar";
import { PostCard } from "./post-card";
import { PostFiltersBar } from "./post-filters";
import { EmptyState } from "@/components/shared/empty-state";
import { ErrorState } from "@/components/shared/error-state";
import { LoadMore } from "@/components/shared/load-more";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
export function PostList({ mode }: { mode: "posts" | "bookmarks" }) {
  const saved = mode === "bookmarks";
  const [filters, setFilters] = useState<PostFilters>({
    search: "",
    feed_id: "",
    unread: false,
  });
  const search = useDebouncedValue(filters.search);
  const query = usePosts(mode, { ...filters, search });
  const directory = useFeeds();
  const follows = useFollows();
  const state = usePostState();
  const [compact] = useState(() => readPreferences().compact);
  const allFeeds = directory.data?.pages.flatMap((page) => page.data) || [];
  const following = new Set(follows.data?.map((f) => f.feed_id));
  const feeds = allFeeds.filter((feed) => following.has(feed.id));
  const posts = query.data?.pages.flatMap((page) => page.data) || [];
  useEffect(() => {
    if (directory.hasNextPage && !directory.isFetching && !directory.isError)
      void directory.fetchNextPage();
  }, [directory]);
  const bookmarks = state.bookmarks.data
    ? new Set(state.bookmarks.data)
    : undefined;
  const unread = state.unread.data ? new Set(state.unread.data) : undefined;
  return (
    <>
      <PageHeader
        eyebrow={
          saved ? "Your personal collection" : "A daily dose of perspective"
        }
        title={saved ? "Worth keeping." : "Your daily reading."}
        description={
          saved
            ? "The stories you want to come back to."
            : "Fresh ideas from the voices you choose to follow."
        }
      >
        {!saved && (
          <Button asChild>
            <Link href="/feeds">
              <Rss />
              Find your next feed
            </Link>
          </Button>
        )}
      </PageHeader>
      <div className="grid items-start gap-8 lg:grid-cols-[minmax(0,1fr)_280px]">
        <section aria-label={saved ? "Bookmarked posts" : "Latest posts"}>
          {!saved && (
            <PostFiltersBar
              filters={filters}
              onChange={setFilters}
              feeds={feeds}
            />
          )}
          <div className="mb-4 flex items-center justify-between">
            <p className="text-xs font-medium uppercase tracking-[0.1em] text-black/50">
              {saved
                ? "Saved for another moment"
                : filters.search
                  ? "Search results"
                  : "The latest"}
            </p>
            <span className="text-xs text-black/40">Newest first</span>
          </div>
          {query.isPending ? (
            <div
              className="space-y-4"
              role="status"
              aria-label="Loading stories"
            >
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="rounded-xl border border-black/[0.08] bg-white p-6"
                >
                  <Skeleton className="mb-4 h-3 w-32" />
                  <Skeleton className="mb-3 h-6 w-3/4" />
                  <Skeleton className="h-14 w-full" />
                </div>
              ))}
            </div>
          ) : query.isError && !query.data ? (
            <ErrorState
              message={query.error.message}
              retry={() => void query.refetch()}
            />
          ) : posts.length === 0 ? (
            <EmptyState
              title={
                saved
                  ? "Keep a little inspiration."
                  : filters.search || filters.feed_id || filters.unread
                    ? "A little quiet here."
                    : "Your next good read starts here."
              }
              description={
                saved
                  ? "Bookmark a story from your home feed and it will be waiting here."
                  : filters.search || filters.feed_id || filters.unread
                    ? "Try a different search or clear your filters to find more stories."
                    : "Follow a few feeds to fill this space. New feeds may take a few minutes to bring in their first stories."
              }
            >
              <Button asChild variant="secondary">
                <Link href={saved ? "/" : "/feeds"}>
                  {saved ? "Explore your home feed" : "Discover feeds"}
                  <ArrowRight />
                </Link>
              </Button>
            </EmptyState>
          ) : (
            <>
              <div className="space-y-4">
                {posts.map((post) => (
                  <PostCard
                    key={post.id}
                    post={post}
                    feedName={
                      allFeeds.find((feed) => feed.id === post.feed_id)?.name
                    }
                    compact={compact}
                    bookmarked={saved ? true : bookmarks?.has(post.id)}
                    read={
                      following.has(post.feed_id) && unread
                        ? !unread.has(post.id)
                        : undefined
                    }
                  />
                ))}
              </div>
              <LoadMore
                hasMore={!!query.hasNextPage}
                loading={query.isFetchingNextPage}
                error={query.isFetchNextPageError}
                load={() => void query.fetchNextPage()}
              />
            </>
          )}
          {(state.bookmarks.isError || state.unread.isError) && (
            <p role="status" className="mt-3 text-xs text-graphite">
              Some reading indicators are unavailable. You can still use each
              story’s action menu.{" "}
              <button
                className="text-notion-blue"
                onClick={() => {
                  void state.bookmarks.refetch();
                  void state.unread.refetch();
                }}
              >
                Retry indicators
              </button>
            </p>
          )}
        </section>
        <ReadingSidebar feeds={feeds} />
      </div>
    </>
  );
}
