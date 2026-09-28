import { cn } from "@/lib/utils";
export function Card({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card"
      className={cn(
        "rounded-xl border border-foreground/[0.08] bg-surface p-6",
        className,
      )}
      {...props}
    />
  );
}
