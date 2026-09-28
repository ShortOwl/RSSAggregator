"use client";
import { DropdownMenu as Primitive } from "radix-ui";
export const DropdownMenu = Primitive.Root;
export const DropdownMenuTrigger = Primitive.Trigger;
export const DropdownMenuRadioGroup = Primitive.RadioGroup;
export function DropdownMenuRadioItem(
  props: React.ComponentProps<typeof Primitive.RadioItem>,
) {
  return (
    <Primitive.RadioItem
      className="flex cursor-pointer items-center gap-2 rounded-sm px-3 py-2 text-sm outline-none focus:bg-sky-tint focus:text-black data-[state=checked]:font-semibold"
      {...props}
    />
  );
}
export function DropdownMenuContent({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Primitive.Portal>
      <Primitive.Content
        align="end"
        sideOffset={8}
        className="z-50 min-w-44 rounded-lg border border-foreground/10 bg-surface p-1"
      >
        {children}
      </Primitive.Content>
    </Primitive.Portal>
  );
}
export function DropdownMenuItem(
  props: React.ComponentProps<typeof Primitive.Item>,
) {
  return (
    <Primitive.Item
      className="flex cursor-pointer items-center gap-2 rounded-sm px-3 py-2 text-sm outline-none focus:bg-sky-tint focus:text-black data-[disabled]:opacity-40"
      {...props}
    />
  );
}
