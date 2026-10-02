import { SourceMark } from "@/components/shared/source-mark";
import Link from "next/link";
import { ArrowRight, Rss, Sparkles } from "lucide-react";
import type { Feed } from "@/lib/types";
import { useFollows } from "@/hooks/use-feeds";
export function ReadingSidebar({
  feeds,
  onSelectFeed,
  selectedFeed,
}: {
  feeds: Feed[];
  onSelectFeed?: (id: string) => void;
  selectedFeed?: string;
}) {
  const follows = useFollows();
  return (
    <aside className="space-y-6 lg:sticky lg:top-28">
      <div className="rounded-xl border border-foreground/[0.08] bg-surface p-6">
        <div className="flex items-center justify-between">
          <h2 className="text-sm font-semibold">Your sources</h2>
          <Rss size={16} className="text-ink-40" />
        </div>
        {follows.isPending ? (
          <p className="mt-4 text-sm text-muted">Loading your sources…</p>
        ) : follows.isError ? (
          <p className="mt-4 text-sm text-muted">Sources are unavailable.</p>
        ) : (
          <>
            <p className="mt-4 text-[40px] font-semibold leading-none tracking-tight">
              {follows.data?.length || 0}
              <span className="ml-2 text-sm font-normal tracking-normal text-ink-50">
                feeds followed
              </span>
            </p>
            <div className="mt-5 space-y-3">
              {feeds.slice(0, 5).map((feed) => (
                <div key={feed.id} className="flex items-center gap-2 text-sm">
                  <SourceMark name={feed.name} />
                  {onSelectFeed ? (
                    <button
                      onClick={() =>
                        onSelectFeed(selectedFeed === feed.id ? "" : feed.id)
                      }
                      aria-pressed={selectedFeed === feed.id}
                      className={
                        "truncate text-left hover:text-notion-blue " +
                        (selectedFeed === feed.id
                          ? "font-medium text-notion-blue"
                          : "text-muted")
                      }
                    >
                      {feed.name}
                    </button>
                  ) : (
                    <Link
                      href="/feeds"
                      className="truncate text-muted hover:text-notion-blue"
                    >
                      {feed.name}
                    </Link>
                  )}
                </div>
              ))}
            </div>
          </>
        )}
        <Link
          href="/feeds"
          className="mt-6 inline-flex items-center gap-2 text-sm font-medium text-notion-blue"
        >
          Manage feeds
          <ArrowRight size={14} />
        </Link>
      </div>
      <div className="rounded-xl bg-marigold p-6 accent-surface">
        <Sparkles size={24} strokeWidth={1.5} />
        <h2 className="mt-8 text-[22px] font-semibold leading-[1.27] tracking-tight">
          Good ideas need
          <br />a little room.
        </h2>
        <p className="mt-3 font-serif text-lg leading-[1.56]">
          A quiet corner of the internet. Curated by you.
        </p>
        <div className="mt-8 border-t border-foreground/15 pt-4 text-xs">
          Less scrolling. More discovering.
        </div>
      </div>
      <p className="px-2 text-xs leading-relaxed text-ink-40">
        A small ritual. A wider world.
        <br />
        Make yourself at home.
      </p>
    </aside>
  );
}
