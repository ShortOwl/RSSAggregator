"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { BookOpen, Bookmark, Rss, Settings, LogOut, Menu } from "lucide-react";
import { BrandMark } from "@/components/shared/brand-mark";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";
import { useAuth } from "@/components/auth/auth-provider";
import { cn } from "@/lib/utils";
const links = [
  { href: "/", label: "Home feed", icon: BookOpen },
  { href: "/feeds", label: "My feeds", icon: Rss },
  { href: "/bookmarks", label: "Bookmarks", icon: Bookmark },
];
export function AppHeader() {
  const pathname = usePathname();
  const { signOut } = useAuth();
  return (
    <>
      <a
        href="#main"
        className="fixed -top-20 left-4 z-50 rounded-lg bg-notion-blue p-3 text-white focus:top-4"
      >
        Skip to content
      </a>
      <header className="fixed inset-x-0 top-0 z-40 h-16 border-b border-black/[0.08] bg-paper-warmth/95">
        <div className="mx-auto flex h-full max-w-[1440px] items-center justify-between px-4 sm:px-8 lg:px-16">
          <Link href="/" aria-label="Margin home">
            <BrandMark />
          </Link>
          <nav
            aria-label="Main navigation"
            className="hidden items-center gap-2 md:flex"
          >
            {links.map(({ href, label, icon: Icon }) => (
              <Link
                key={href}
                href={href}
                aria-current={pathname === href ? "page" : undefined}
                className={cn(
                  "flex items-center gap-2 rounded-lg px-4 py-3 text-sm font-medium text-black/55 hover:text-black",
                  pathname === href && "bg-sky-tint text-notion-blue",
                )}
              >
                <Icon size={16} />
                {label}
              </Link>
            ))}
          </nav>
          <div className="flex items-center gap-2">
            <Link
              href="/settings"
              aria-label="Settings and profile"
              className={cn(
                "rounded-lg p-2 text-black/55 hover:bg-black/5",
                pathname === "/settings" && "text-notion-blue",
              )}
            >
              <Settings size={20} />
            </Link>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  aria-label="Open account menu"
                >
                  <Menu size={20} />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                {links.map(({ href, label, icon: Icon }) => (
                  <DropdownMenuItem key={href} asChild>
                    <Link href={href}>
                      <Icon size={16} />
                      {label}
                    </Link>
                  </DropdownMenuItem>
                ))}
                <DropdownMenuItem onSelect={signOut}>
                  <LogOut size={16} />
                  Sign out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>
    </>
  );
}
