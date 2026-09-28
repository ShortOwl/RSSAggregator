"use client";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { LogOut, UserRound, BookOpen } from "lucide-react";
import { toast } from "sonner";
import { api } from "@/lib/api";
import { keys } from "@/lib/query-keys";
import { dateLabel } from "@/lib/format";
import { preferencesKey, readPreferences } from "@/lib/preferences";
import type { User } from "@/lib/types";
import { useAuth } from "@/components/auth/auth-provider";
import { PageHeader } from "@/components/layout/page-header";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { ErrorState } from "@/components/shared/error-state";
import { Select } from "@/components/ui/select";
import { useTheme } from "@/components/theme-provider";
import { parseTheme } from "@/lib/theme";
export function ProfileSettings() {
  const { signOut } = useAuth();
  const { theme, setTheme } = useTheme();
  const user = useQuery({
    queryKey: keys.user,
    queryFn: ({ signal }) => api<User>("/users", { signal }),
  });
  const [compact, setCompact] = useState(() => readPreferences().compact);
  function change(value: boolean) {
    try {
      localStorage.setItem(preferencesKey, JSON.stringify({ compact: value }));
      setCompact(value);
      toast.success("Reading preference saved");
    } catch {
      toast.error("Your browser couldn’t save this preference.");
    }
  }
  return (
    <>
      <PageHeader
        eyebrow="Make yourself at home"
        title="Your little corner."
        description="A few details, a few preferences. All yours."
      />
      <div className="grid max-w-5xl items-start gap-6 lg:grid-cols-[minmax(0,1fr)_280px]">
        <div className="space-y-6">
          <Card>
            <div className="mb-6 flex items-center gap-3">
              <UserRound size={20} />
              <h2 className="text-[22px] font-semibold tracking-tight">
                Profile
              </h2>
            </div>
            {user.isPending ? (
              <Skeleton className="h-36" />
            ) : user.isError ? (
              <ErrorState
                message={user.error.message}
                retry={() => void user.refetch()}
              />
            ) : (
              <dl className="space-y-5">
                <div>
                  <dt className="text-xs text-ink-50">Name</dt>
                  <dd className="mt-1 font-medium">{user.data.name}</dd>
                </div>
                <div>
                  <dt className="text-xs text-ink-50">Email address</dt>
                  <dd className="mt-1 break-all">{user.data.email}</dd>
                </div>
                <div>
                  <dt className="text-xs text-ink-50">Reading with us since</dt>
                  <dd className="mt-1 text-sm">
                    {dateLabel(user.data.created_at)}
                  </dd>
                </div>
              </dl>
            )}
          </Card>
          <Card>
            <div className="mb-5 flex items-center gap-3">
              <BookOpen size={20} />
              <h2 className="text-[22px] font-semibold tracking-tight">
                Reading preferences
              </h2>
            </div>
            <div className="mb-5 flex items-center justify-between gap-6">
              <label htmlFor="theme" className="text-sm font-medium">
                Appearance
              </label>
              <Select
                id="theme"
                value={theme}
                onChange={(event) => setTheme(parseTheme(event.target.value))}
              >
                <option value="light">Light</option>
                <option value="dark">Dark</option>
                <option value="system">System</option>
              </Select>
            </div>
            <label className="flex cursor-pointer items-center justify-between gap-6">
              <span>
                <span className="block text-sm font-medium">
                  Compact stories
                </span>
                <span className="mt-1 block text-sm text-muted">
                  Show headlines without article previews.
                </span>
              </span>
              <input
                type="checkbox"
                checked={compact}
                onChange={(event) => change(event.target.checked)}
                className="size-5 accent-notion-blue"
              />
            </label>
            <p className="mt-5 border-t border-foreground/[0.08] pt-4 text-xs text-ink-40">
              Saved in this browser. System follows your device’s appearance.
            </p>
          </Card>
          <Card>
            <h2 className="text-[22px] font-semibold tracking-tight">
              Until next time.
            </h2>
            <p className="mt-2 text-sm text-muted">
              Your feeds and bookmarks will be here when you return.
            </p>
            <Button variant="secondary" className="mt-5" onClick={signOut}>
              <LogOut />
              Sign out
            </Button>
          </Card>
        </div>
        <aside className="rounded-xl bg-midnight-ink p-6 text-white">
          <p className="text-xs uppercase tracking-[0.1em] text-white/60">
            A note to the reader
          </p>
          <p className="mt-8 font-serif text-[32px] leading-[1.25]">
            Stay curious.
            <br />
            Read widely.
            <br />
            Make room.
          </p>
          <p className="mt-8 text-sm text-white/70">
            The best reading list is the one that’s yours.
          </p>
        </aside>
      </div>
    </>
  );
}
