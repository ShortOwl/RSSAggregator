"use client";
import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, ArrowRight } from "lucide-react";
import { toast } from "sonner";
import { api } from "@/lib/api";
import { keys } from "@/lib/query-keys";
import { safeUrl } from "@/lib/format";
import type { Feed, FeedFollow } from "@/lib/types";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
export function AddFeedDialog() {
  const [open, setOpen] = useState(false);
  const [created, setCreated] = useState<Feed | null>(null);
  const [validation, setValidation] = useState("");
  const client = useQueryClient();
  const mutation = useMutation({
    mutationFn: async (body: { name: string; url: string }) => {
      const feed =
        created ||
        (await api<Feed>("/feeds", {
          method: "POST",
          body: JSON.stringify(body),
        }));
      setCreated(feed);
      void client.invalidateQueries({ queryKey: keys.feeds });
      await api<FeedFollow>("/feed_follows", {
        method: "POST",
        body: JSON.stringify({ feed_id: feed.id }),
      });
    },
    onSuccess: () => {
      [keys.feeds, keys.follows, keys.posts, keys.state].forEach(
        (queryKey) => void client.invalidateQueries({ queryKey }),
      );
      toast.success("Feed added to your reading space");
      setOpen(false);
      setCreated(null);
    },
    onError: () => {},
  });
  function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setValidation("");
    const data = new FormData(event.currentTarget);
    const name = String(data.get("name") || created?.name || "").trim();
    const url = String(data.get("url") || created?.url || "").trim();
    if (!name || !safeUrl(url)) {
      setValidation("Enter a name and a valid HTTP or HTTPS feed URL.");
      return;
    }
    mutation.mutate({ name, url });
  }
  return (
    <Dialog
      open={open}
      onOpenChange={(value) => {
        if (!mutation.isPending) {
          setOpen(value);
          setValidation("");
          mutation.reset();
        }
      }}
    >
      <DialogTrigger asChild>
        <Button>
          <Plus />
          Add a feed
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogTitle className="pr-8 text-[22px] font-semibold tracking-tight">
          A new voice to follow.
        </DialogTitle>
        <DialogDescription className="mt-2 text-sm text-muted">
          Add an RSS or Atom feed to your reading space.
        </DialogDescription>
        <form onSubmit={submit} className="mt-6 space-y-5">
          <div className="space-y-2">
            <Label htmlFor="feed-name">Feed name</Label>
            <Input
              id="feed-name"
              name="name"
              placeholder="A favorite publication"
              defaultValue={created?.name}
              disabled={!!created || mutation.isPending}
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="feed-url">Feed URL</Label>
            <Input
              id="feed-url"
              name="url"
              type="url"
              placeholder="https://example.com/feed.xml"
              defaultValue={created?.url}
              disabled={!!created || mutation.isPending}
              required
            />
          </div>
          {created && mutation.isError && (
            <p className="text-sm text-muted">
              Your feed was created. Following it didn’t finish; retry below
              without creating it again.
            </p>
          )}
          {(validation || mutation.error) && (
            <p role="alert" className="text-sm text-vermillion">
              {validation || mutation.error?.message}
            </p>
          )}
          <p className="text-xs leading-relaxed text-ink-50">
            Stories arrive after the next feed refresh, usually within 10
            minutes.
          </p>
          <Button
            type="submit"
            disabled={mutation.isPending}
            className="w-full"
          >
            {mutation.isPending
              ? "Adding your feed…"
              : created
                ? "Retry following"
                : "Add and follow"}
            <ArrowRight />
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
