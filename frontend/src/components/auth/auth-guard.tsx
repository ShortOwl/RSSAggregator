"use client";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "./auth-provider";
export function AuthGuard({ children }: { children: React.ReactNode }) {
  const { ready, authenticated } = useAuth();
  const router = useRouter();
  useEffect(() => {
    if (ready && !authenticated) router.replace("/login");
  }, [ready, authenticated, router]);
  if (!ready || !authenticated)
    return (
      <div
        className="flex min-h-screen items-center justify-center text-sm text-graphite"
        role="status"
      >
        Opening your reading space…
      </div>
    );
  return children;
}
