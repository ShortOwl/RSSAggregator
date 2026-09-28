"use client";
import { Bookmark, BookmarkCheck, Check, MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { usePostAction } from "@/hooks/use-post-state";
export function PostActions({
  id,
  bookmarked,
  read,
}: {
  id: string;
  bookmarked?: boolean;
  read?: boolean;
}) {
  const action = usePostAction(id);
  return (
    <div className="flex shrink-0 items-center gap-1">
      {bookmarked !== undefined ? (
        <Button
          variant="ghost"
          size="icon"
          aria-label={bookmarked ? "Remove bookmark" : "Bookmark post"}
          aria-pressed={bookmarked}
          disabled={action.isPending}
          onClick={() =>
            action.mutate({ kind: "bookmark", active: !bookmarked })
          }
        >
          {bookmarked ? (
            <BookmarkCheck className="text-notion-blue" />
          ) : (
            <Bookmark />
          )}
        </Button>
      ) : null}
      {read !== undefined ? (
        <Button
          variant="ghost"
          size="icon"
          aria-label={read ? "Mark as unread" : "Mark as read"}
          aria-pressed={read}
          disabled={action.isPending}
          onClick={() => action.mutate({ kind: "read", active: !read })}
        >
          <Check className={read ? "text-notion-blue" : "text-ink-40"} />
        </Button>
      ) : null}
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            aria-label="Post actions"
            disabled={action.isPending}
          >
            <MoreHorizontal />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem
            onSelect={() => action.mutate({ kind: "read", active: true })}
          >
            Mark as read
          </DropdownMenuItem>
          <DropdownMenuItem
            onSelect={() => action.mutate({ kind: "read", active: false })}
          >
            Mark as unread
          </DropdownMenuItem>
          {bookmarked === undefined && (
            <>
              <DropdownMenuItem
                onSelect={() =>
                  action.mutate({ kind: "bookmark", active: true })
                }
              >
                Save bookmark
              </DropdownMenuItem>
              <DropdownMenuItem
                onSelect={() =>
                  action.mutate({ kind: "bookmark", active: false })
                }
              >
                Remove bookmark
              </DropdownMenuItem>
            </>
          )}
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
