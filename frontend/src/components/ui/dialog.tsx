"use client";
import { Dialog as Primitive } from "radix-ui";
import { X } from "lucide-react";
export const Dialog = Primitive.Root;
export const DialogTrigger = Primitive.Trigger;
export const DialogTitle = Primitive.Title;
export const DialogDescription = Primitive.Description;
export function DialogContent({ children }: { children: React.ReactNode }) {
  return (
    <Primitive.Portal>
      <Primitive.Overlay className="fixed inset-0 z-50 bg-black/30" />
      <Primitive.Content className="fixed left-1/2 top-1/2 z-50 w-[calc(100%-32px)] max-w-lg -translate-x-1/2 -translate-y-1/2 rounded-xl border border-foreground/10 bg-surface p-6 focus:outline-none">
        {children}
        <Primitive.Close
          aria-label="Close dialog"
          className="absolute right-4 top-4 rounded-lg p-2 hover:bg-foreground/5"
        >
          <X size={18} />
        </Primitive.Close>
      </Primitive.Content>
    </Primitive.Portal>
  );
}
