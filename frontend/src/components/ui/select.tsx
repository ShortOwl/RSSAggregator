import { cn } from "@/lib/utils";
export function Select({
  className,
  ...props
}: React.ComponentProps<"select">) {
  return (
    <select
      className={cn(
        "min-h-11 rounded-lg border border-black/15 bg-white px-3 py-2 text-sm focus-visible:outline-notion-blue",
        className,
      )}
      {...props}
    />
  );
}
