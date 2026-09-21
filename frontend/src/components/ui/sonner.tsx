"use client";
import { Toaster as Sonner } from "sonner";
export function Toaster() {
  return (
    <Sonner
      theme="light"
      position="bottom-right"
      toastOptions={{
        style: {
          background: "#fff",
          color: "#111",
          border: "1px solid rgba(0,0,0,0.08)",
          borderRadius: 12,
          boxShadow: "none",
        },
      }}
    />
  );
}
