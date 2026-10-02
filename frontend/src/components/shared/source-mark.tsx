import { cn } from "@/lib/utils";

export function SourceMark({
  name,
  large = false,
}: {
  name: string;
  large?: boolean;
}) {
  const colors = ["bg-sky-tint", "bg-marigold", "bg-sky-wash"];
  const hash = Array.from(name).reduce(
    (value, letter) => (value + letter.codePointAt(0)!) % colors.length,
    0,
  );
  const initials =
    name
      .trim()
      .split(/\s+/)
      .slice(0, 2)
      .map((part) => Array.from(part)[0])
      .join("")
      .toUpperCase() || "R";
  return (
    <span
      aria-hidden="true"
      className={cn(
        "inline-flex shrink-0 items-center justify-center rounded-sm font-medium text-black",
        colors[hash],
        large ? "size-10 text-sm" : "size-6 text-xs",
      )}
    >
      {initials}
    </span>
  );
}
