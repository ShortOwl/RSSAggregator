"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ArrowRight, Eye, EyeOff } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { AuthResponse } from "@/lib/types";
import { saveToken } from "@/lib/session";
import { useAuth } from "./auth-provider";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
export function AuthForm({ mode }: { mode: "login" | "register" }) {
  const register = mode === "register";
  const router = useRouter();
  const { authenticated, ready } = useAuth();
  const [visible, setVisible] = useState(false);
  const [validation, setValidation] = useState("");
  useEffect(() => {
    if (ready && authenticated) router.replace("/");
  }, [ready, authenticated, router]);
  const mutation = useMutation({
    mutationFn: (body: Record<string, string>) =>
      api<AuthResponse>(
        `/${mode}`,
        { method: "POST", body: JSON.stringify(body) },
        false,
      ),
    onSuccess: (data) => {
      saveToken(data.token);
      router.replace("/");
    },
  });
  function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setValidation("");
    const data = new FormData(event.currentTarget);
    const password = String(data.get("password"));
    if (register) {
      const length = new TextEncoder().encode(password).length;
      if (length < 8 || length > 72) {
        setValidation(
          "Use a password between 8 and 72 bytes (some characters use more than one byte).",
        );
        return;
      }
    }
    mutation.mutate({
      email: String(data.get("email")).trim(),
      password,
      ...(register ? { name: String(data.get("name")).trim() } : {}),
    });
  }
  return (
    <Card className="mx-auto w-full max-w-md p-8">
      <p className="mb-2 text-xs font-medium uppercase tracking-[0.12em] text-ink-40">
        Your reading space
      </p>
      <h2 className="text-[22px] font-semibold tracking-[-0.242px]">
        {register ? "A fresh page starts here." : "Welcome back."}
      </h2>
      <p className="mt-2 text-sm text-muted">
        {register
          ? "Create an account and follow your curiosity."
          : "Good stories are waiting for you."}
      </p>
      <form className="mt-8 space-y-5" onSubmit={submit}>
        {register && (
          <div className="space-y-2">
            <Label htmlFor="name">Your name</Label>
            <Input
              id="name"
              name="name"
              autoComplete="name"
              placeholder="Alex Morgan"
              required
              maxLength={100}
            />
          </div>
        )}
        <div className="space-y-2">
          <Label htmlFor="email">Email address</Label>
          <Input
            id="email"
            type="email"
            name="email"
            autoComplete="email"
            placeholder="you@example.com"
            required
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="password">Password</Label>
          <div className="relative">
            <Input
              id="password"
              name="password"
              type={visible ? "text" : "password"}
              autoComplete={register ? "new-password" : "current-password"}
              required
              className="pr-12"
              aria-describedby={register ? "password-hint" : undefined}
            />
            <button
              type="button"
              onClick={() => setVisible(!visible)}
              aria-label={visible ? "Hide password" : "Show password"}
              className="absolute right-1 top-1 rounded-lg p-2 text-ink-50"
            >
              {visible ? <EyeOff size={18} /> : <Eye size={18} />}
            </button>
          </div>
          {register && (
            <p id="password-hint" className="text-xs text-ink-50">
              At least 8 characters for a fresh start.
            </p>
          )}
        </div>
        {!register && (
          <p className="text-right text-sm">
            <Link
              className="font-medium text-notion-blue hover:underline"
              href="/forgot-password"
            >
              Forgot password?
            </Link>
          </p>
        )}
        {(validation || mutation.error) && (
          <p
            role="alert"
            className="rounded-lg bg-canvas p-3 text-sm text-vermillion"
          >
            {validation || mutation.error?.message}
          </p>
        )}
        <Button
          type="submit"
          className="w-full py-3"
          disabled={mutation.isPending}
        >
          {mutation.isPending
            ? "One moment…"
            : register
              ? "Create account"
              : "Sign in"}
          <ArrowRight size={16} />
        </Button>
      </form>
      <p className="mt-6 text-center text-sm text-muted">
        {register ? "Already have an account? " : "New to Margin? "}
        <Link
          className="font-medium text-notion-blue hover:underline"
          href={register ? "/login" : "/register"}
        >
          {register ? "Sign in" : "Create an account"}
        </Link>
      </p>
    </Card>
  );
}
