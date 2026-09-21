import Link from "next/link";
import { EmptyState } from "@/components/shared/empty-state";
import { Button } from "@/components/ui/button";
export default function NotFound() {
  return (
    <main className="mx-auto max-w-xl px-4 py-20">
      <EmptyState
        title="A page out of place."
        description="This page doesn’t exist. Your reading space is just one click away."
      >
        <Button asChild>
          <Link href="/">Back to home</Link>
        </Button>
      </EmptyState>
    </main>
  );
}
