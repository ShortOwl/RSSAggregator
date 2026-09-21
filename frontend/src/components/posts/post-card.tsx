import { ArrowUpRight } from "lucide-react";
import type { Post } from "@/lib/types";
import { dateLabel, excerpt, hostname, safeUrl } from "@/lib/format";
import { Card } from "@/components/ui/card";
import { PostActions } from "./post-actions";
export function PostCard({
  post,
  feedName,
  read,
  bookmarked,
  compact,
}: {
  post: Post;
  feedName?: string;
  read?: boolean;
  bookmarked?: boolean;
  compact: boolean;
}) {
  const description = excerpt(post.description);
  const url = safeUrl(post.url);
  return (
    <Card className="group transition-colors duration-200 hover:border-black/20">
      <article>
        <div className="mb-3 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-black/50">
          <span className="font-medium text-black/70">
            {feedName || hostname(post.url)}
          </span>
          <span aria-hidden="true">·</span>
          <time dateTime={post.published_at}>
            {dateLabel(post.published_at)}
          </time>
          {read === false && (
            <span className="inline-flex items-center gap-1.5 text-notion-blue">
              <span className="size-1.5 rounded-full bg-notion-blue" />
              Unread
            </span>
          )}
          {read === true && <span>Read</span>}
        </div>
        <h2 className="text-[22px] font-semibold leading-[1.27] tracking-[-0.242px]">
          {url ? (
            <a
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-notion-blue"
            >
              {post.title || "Untitled story"}
              <ArrowUpRight className="ml-2 inline size-4 text-black/30" />
              <span className="sr-only"> (opens in a new tab)</span>
            </a>
          ) : (
            post.title || "Untitled story"
          )}
        </h2>
        {!compact && description && (
          <p className="mt-3 line-clamp-3 text-sm leading-[1.6] text-graphite">
            {description}
          </p>
        )}
        <div className="mt-4 flex items-center justify-between gap-2">
          <span className="text-xs text-black/40">{hostname(post.url)}</span>
          <PostActions id={post.id} read={read} bookmarked={bookmarked} />
        </div>
      </article>
    </Card>
  );
}
