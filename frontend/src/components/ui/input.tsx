import { cn } from "@/lib/utils";
export function Input({ className, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      data-slot="input"
      className={cn(
        "w-full min-h-11 rounded-lg border border-black/15 bg-white px-3 py-2 text-sm placeholder:text-black/40 focus-visible:outline-2 focus-visible:outline-notion-blue disabled:opacity-50",
        className,
      )}
      {...props}
    />
  );
}
