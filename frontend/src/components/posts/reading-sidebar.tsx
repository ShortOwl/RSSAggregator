import Link from "next/link";
import { ArrowRight, Rss, Sparkles } from "lucide-react";
import type { Feed } from "@/lib/types";
import { useFollows } from "@/hooks/use-feeds";
export function ReadingSidebar({ feeds }: { feeds: Feed[] }) {
  const follows = useFollows();
  return (
    <aside className="space-y-6 lg:sticky lg:top-28">
      <div className="rounded-xl bg-marigold p-6">
        <Sparkles size={24} strokeWidth={1.5} />
        <h2 className="mt-8 text-[22px] font-semibold leading-[1.27] tracking-tight">
          Good ideas need
          <br />a little room.
        </h2>
        <p className="mt-3 font-serif text-lg leading-[1.56]">
          A quiet corner of the internet. Curated by you.
        </p>
        <div className="mt-8 border-t border-black/15 pt-4 text-xs">
          Less scrolling. More discovering.
        </div>
      </div>
      <div className="rounded-xl border border-black/[0.08] bg-white p-6">
        <div className="flex items-center justify-between">
          <h2 className="text-sm font-semibold">Your sources</h2>
          <Rss size={16} className="text-black/40" />
        </div>
        {follows.isPending ? (
          <p className="mt-4 text-sm text-graphite">Loading your sources…</p>
        ) : follows.isError ? (
          <p className="mt-4 text-sm text-graphite">Sources are unavailable.</p>
        ) : (
          <>
            <p className="mt-4 text-[40px] font-semibold leading-none tracking-tight">
              {follows.data?.length || 0}
              <span className="ml-2 text-sm font-normal tracking-normal text-black/50">
                feeds followed
              </span>
            </p>
            <div className="mt-5 space-y-3">
              {feeds.slice(0, 5).map((feed, index) => (
                <p
                  key={feed.id}
                  className="flex items-center gap-2 truncate text-sm text-graphite"
                >
                  <span
                    className={
                      ["bg-sky-tint", "bg-marigold", "bg-paper-warmth"][
                        index % 3
                      ] +
                      " flex size-6 shrink-0 items-center justify-center rounded-sm text-xs text-black"
                    }
                  >
                    {feed.name.slice(0, 1).toUpperCase()}
                  </span>
                  {feed.name}
                </p>
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
      <p className="px-2 text-xs leading-relaxed text-black/40">
        A small ritual. A wider world.
        <br />
        Make yourself at home.
      </p>
    </aside>
  );
}
