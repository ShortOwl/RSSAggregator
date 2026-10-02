"use client";
import { useEffect, useState } from "react";
import { Search } from "lucide-react";
import { useFeeds, useFollows } from "@/hooks/use-feeds";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { PageHeader } from "@/components/layout/page-header";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { ErrorState } from "@/components/shared/error-state";
import { LoadMore } from "@/components/shared/load-more";
import { FeedCard } from "./feed-card";
import { AddFeedDialog } from "./add-feed-dialog";
export function FeedManager() {
  const [search, setSearch] = useState("");
  const [tab, setTab] = useState<"following" | "discover">("following");
  const query = useFeeds(useDebouncedValue(search));
  const follows = useFollows();
  const all = query.data?.pages.flatMap((page) => page.data) || [];
  const following = new Map(
    follows.data?.map((follow) => [follow.feed_id, follow]),
  );
  const feeds = all.filter(
    (feed) => tab === "discover" || following.has(feed.id),
  );
  useEffect(() => {
    if (
      tab === "following" &&
      query.hasNextPage &&
      !query.isFetching &&
      !query.isError
    )
      void query.fetchNextPage();
  }, [tab, query]);
  return (
    <>
      <PageHeader
        eyebrow="Choose your corner of the internet"
        title="Follow your curiosity."
        description="Bring your favorite publications, blogs, and big ideas together."
      >
        <AddFeedDialog />
      </PageHeader>
      <div className="mb-8 rounded-xl bg-sky-tint p-6 accent-surface sm:flex sm:items-center sm:justify-between sm:gap-6">
        <div>
          <h2 className="text-[22px] font-semibold tracking-tight">
            A feed that feels like you.
          </h2>
          <p className="mt-2 text-sm text-muted">
            Follow a source below, or add one you already love.
          </p>
        </div>
        <p className="mt-4 font-serif text-lg text-midnight-ink sm:mt-0">
          Your interests. Your pace.
        </p>
      </div>
      <div className="mb-6 flex flex-col justify-between gap-4 sm:flex-row">
        <div className="flex gap-2" aria-label="Feed view">
          <Button
            variant={tab === "following" ? "secondary" : "ghost"}
            aria-pressed={tab === "following"}
            onClick={() => setTab("following")}
          >
            Following {follows.data ? `(${follows.data.length})` : ""}
          </Button>
          <Button
            variant={tab === "discover" ? "secondary" : "ghost"}
            aria-pressed={tab === "discover"}
            onClick={() => setTab("discover")}
          >
            Discover feeds
          </Button>
        </div>
        <div className="relative sm:w-80">
          <Search size={16} className="absolute left-3 top-3.5 text-ink-40" />
          <Input
            aria-label="Search feeds"
            placeholder="Search by feed name or URL"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            className="pl-10"
          />
        </div>
      </div>
      {follows.isError && (
        <div className="mb-4">
          <ErrorState
            message={follows.error.message}
            retry={() => void follows.refetch()}
          />
        </div>
      )}
      {query.isPending || follows.isPending ? (
        <div
          className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3"
          role="status"
          aria-label="Loading feeds"
        >
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-56 rounded-xl" />
          ))}
        </div>
      ) : query.isError && !query.data ? (
        <ErrorState
          message={query.error.message}
          retry={() => void query.refetch()}
        />
      ) : feeds.length ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {feeds.map((feed) => (
            <FeedCard
              key={feed.id}
              feed={feed}
              follow={following.get(feed.id)}
              ready={follows.isSuccess}
            />
          ))}
        </div>
      ) : (
        <EmptyState
          title={
            query.isFetching
              ? "Finding your feeds…"
              : search
                ? "No feeds found."
                : "Every good habit starts with one."
          }
          description={
            search
              ? "Try a different name or URL."
              : "Discover a new publication or add your favorite RSS feed."
          }
        >
          {!search && (
            <Button variant="secondary" onClick={() => setTab("discover")}>
              Discover feeds
            </Button>
          )}
        </EmptyState>
      )}
      {all.length > 0 && (
        <LoadMore
          kind="feeds"
          hasMore={!!query.hasNextPage}
          loading={query.isFetchingNextPage}
          error={query.isFetchNextPageError}
          load={() => void query.fetchNextPage()}
        />
      )}
    </>
  );
}
