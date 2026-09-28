"use client";
import { Search, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Select } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import type { Feed, PostFilters } from "@/lib/types";
export function PostFiltersBar({
  filters,
  onChange,
  feeds,
}: {
  filters: PostFilters;
  onChange: (value: PostFilters) => void;
  feeds: Feed[];
}) {
  return (
    <div className="mb-6 flex flex-col gap-3 sm:flex-row">
      <div className="relative flex-1">
        <Search className="absolute left-3 top-3.5 size-4 text-ink-40" />
        <Input
          aria-label="Search posts"
          placeholder="Find a story, an idea, a little inspiration…"
          value={filters.search}
          onChange={(event) =>
            onChange({ ...filters, search: event.target.value })
          }
          className="pl-10 pr-10"
        />
        {filters.search && (
          <button
            aria-label="Clear search"
            className="absolute right-3 top-3.5"
            onClick={() => onChange({ ...filters, search: "" })}
          >
            <X size={16} />
          </button>
        )}
      </div>
      <Select
        aria-label="Filter by feed"
        value={filters.feed_id}
        onChange={(event) =>
          onChange({ ...filters, feed_id: event.target.value })
        }
      >
        <option value="">All followed feeds</option>
        {feeds.map((feed) => (
          <option key={feed.id} value={feed.id}>
            {feed.name}
          </option>
        ))}
      </Select>
      <Button
        variant={filters.unread ? "secondary" : "ghost"}
        aria-pressed={filters.unread}
        onClick={() => onChange({ ...filters, unread: !filters.unread })}
      >
        Unread only
      </Button>
    </div>
  );
}
