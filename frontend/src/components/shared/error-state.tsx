import { Button } from "@/components/ui/button";
export function ErrorState({
  message,
  retry,
}: {
  message: string;
  retry: () => void;
}) {
  return (
    <div
      role="alert"
      className="rounded-xl border border-foreground/10 bg-surface p-6"
    >
      <p className="font-medium">We couldn’t load this just yet.</p>
      <p className="mt-2 text-sm text-muted">{message}</p>
      <Button variant="secondary" className="mt-4" onClick={retry}>
        Try again
      </Button>
    </div>
  );
}
