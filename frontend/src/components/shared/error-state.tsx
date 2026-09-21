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
      className="rounded-xl border border-black/10 bg-white p-6"
    >
      <p className="font-medium">We couldn’t load this just yet.</p>
      <p className="mt-2 text-sm text-graphite">{message}</p>
      <Button variant="secondary" className="mt-4" onClick={retry}>
        Try again
      </Button>
    </div>
  );
}
