import { BookOpen } from "lucide-react";
export function BrandMark() {
  return (
    <span className="inline-flex items-center gap-2 text-xl font-semibold tracking-tight">
      <span className="flex size-8 items-center justify-center rounded-lg bg-marigold text-black">
        <BookOpen size={20} strokeWidth={1.8} />
      </span>
      margin<span className="text-notion-blue">.</span>
    </span>
  );
}
