import { cn } from "@/lib/utils";
export function Badge({ className, ...props }: React.ComponentProps<"span">) {
  return (
    <span
      data-slot="badge"
      className={cn(
        "inline-flex rounded-full bg-sky-tint px-3 py-1 text-xs font-medium",
        className,
      )}
      {...props}
    />
  );
}
