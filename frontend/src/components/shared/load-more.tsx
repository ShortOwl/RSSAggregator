"use client";
import { useCallback } from "react";
import { useInfiniteScroll } from "@/hooks/use-infinite-scroll";
import { Button } from "@/components/ui/button";
export function LoadMore({
  hasMore,
  loading,
  error,
  load,
  kind = "stories",
}: {
  hasMore: boolean;
  loading: boolean;
  error: boolean;
  load: () => void;
  kind?: "stories" | "feeds";
}) {
  const callback = useCallback(() => load(), [load]);
  const ref = useInfiniteScroll(callback, hasMore && !loading && !error);
  return (
    <div ref={ref} className="flex justify-center py-8" aria-live="polite">
      {hasMore ? (
        <Button variant="ghost" disabled={loading} onClick={load}>
          {loading
            ? `Gathering more ${kind}…`
            : error
              ? "Retry loading more"
              : `Load more ${kind}`}
        </Button>
      ) : (
        <p className="text-xs text-ink-40">
          {kind === "feeds"
            ? "All feeds shown."
            : "You’re all caught up. A little space to think."}
        </p>
      )}
    </div>
  );
}
