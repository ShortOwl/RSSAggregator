"use client";
import { ErrorState } from "@/components/shared/error-state";
export default function ErrorPage({ reset }: { reset: () => void }) {
  return (
    <main className="mx-auto max-w-xl px-4 py-20">
      <ErrorState
        message="Please try again to reopen your reading space."
        retry={reset}
      />
    </main>
  );
}
