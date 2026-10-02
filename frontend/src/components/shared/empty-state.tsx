import { BookOpen, type LucideIcon } from "lucide-react";
export function EmptyState({
  title,
  description,
  children,
  icon: Icon = BookOpen,
}: {
  title: string;
  description: string;
  children?: React.ReactNode;
  icon?: LucideIcon;
}) {
  return (
    <div className="rounded-xl border border-foreground/[0.08] bg-surface px-6 py-16 text-center">
      <span className="mx-auto mb-6 flex size-16 items-center justify-center rounded-xl bg-sky-tint text-black">
        <Icon size={28} strokeWidth={1.5} />
      </span>
      <h2 className="text-[22px] font-semibold tracking-tight">{title}</h2>
      <p className="mx-auto mt-2 max-w-sm text-sm leading-relaxed text-muted">
        {description}
      </p>
      <div className="mt-6">{children}</div>
    </div>
  );
}
