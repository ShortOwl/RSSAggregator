"use client";
import { useCallback } from "react";
import { useInfiniteScroll } from "@/hooks/use-infinite-scroll";
import { Button } from "@/components/ui/button";
export function LoadMore({
  hasMore,
  loading,
  error,
  load,
}: {
  hasMore: boolean;
  loading: boolean;
  error: boolean;
  load: () => void;
}) {
  const callback = useCallback(() => load(), [load]);
  const ref = useInfiniteScroll(callback, hasMore && !loading && !error);
  return (
    <div ref={ref} className="flex justify-center py-8" aria-live="polite">
      {hasMore ? (
        <Button variant="ghost" disabled={loading} onClick={load}>
          {loading
            ? "Gathering more stories…"
            : error
              ? "Retry loading more"
              : "Load more stories"}
        </Button>
      ) : (
        <p className="text-xs text-black/40">
          You’re all caught up. A little space to think.
        </p>
      )}
    </div>
  );
}
