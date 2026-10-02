"use client";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Check, Plus, ArrowUpRight } from "lucide-react";
import { toast } from "sonner";
import { api } from "@/lib/api";
import { keys } from "@/lib/query-keys";
import { hostname, safeUrl } from "@/lib/format";
import type { Feed, FeedFollow } from "@/lib/types";
import { SourceMark } from "@/components/shared/source-mark";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
export function FeedCard({
  feed,
  follow,
  ready,
}: {
  feed: Feed;
  follow?: FeedFollow;
  ready: boolean;
}) {
  const client = useQueryClient();
  const mutation = useMutation({
    mutationFn: async () => {
      await (follow
        ? api<void>(`/feed_follows/${follow.id}`, { method: "DELETE" })
        : api<FeedFollow>("/feed_follows", {
            method: "POST",
            body: JSON.stringify({ feed_id: feed.id }),
          }));
    },
    onSuccess: () => {
      [keys.follows, keys.posts, keys.state].forEach(
        (queryKey) => void client.invalidateQueries({ queryKey }),
      );
      toast.success(follow ? "Feed unfollowed" : "Feed followed");
    },
    onError: (error) => toast.error(error.message),
  });
  const url = safeUrl(feed.url);
  return (
    <Card className="flex h-full flex-col">
      <div className="mb-6 flex items-start justify-between">
        <SourceMark name={feed.name} large />
        <Button
          variant={follow ? "ghost" : "secondary"}
          disabled={!ready || mutation.isPending}
          onClick={() => mutation.mutate()}
          aria-label={`${follow ? "Unfollow" : "Follow"} ${feed.name}`}
        >
          {mutation.isPending ? (
            "Updating…"
          ) : follow ? (
            <>
              <Check />
              Following
            </>
          ) : (
            <>
              <Plus />
              Follow
            </>
          )}
        </Button>
      </div>
      <h2 className="break-words text-[22px] font-semibold leading-[1.27] tracking-tight">
        {feed.name}
      </h2>
      <p className="mt-2 break-all text-sm text-muted">{hostname(feed.url)}</p>
      <div className="mt-auto pt-6">
        {url && (
          <a
            href={url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 text-xs text-ink-50 hover:text-notion-blue"
          >
            Visit source
            <ArrowUpRight size={14} />
            <span className="sr-only"> (opens in a new tab)</span>
          </a>
        )}
      </div>
    </Card>
  );
}
