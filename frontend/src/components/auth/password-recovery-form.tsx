"use client";

import { useRef, useState } from "react";
import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type Step = "email" | "code" | "done";

export function PasswordRecoveryForm() {
  const [step, setStep] = useState<Step>("email");
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [busy, setBusy] = useState(false);
  const otpInput = useRef<HTMLInputElement>(null);

  async function resendCode() {
    if (busy) return;
    setError("");
    setNotice("");
    setBusy(true);
    try {
      await api<{ message: string }>(
        "/forgot-password",
        { method: "POST", body: JSON.stringify({ email }) },
        false,
      );
      if (otpInput.current) otpInput.current.value = "";
      setNotice("If this email is registered, a new code is on its way.");
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Please try again.");
    } finally {
      setBusy(false);
    }
  }

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (busy) return;
    const data = new FormData(event.currentTarget);
    const address = step === "email" ? String(data.get("email")).trim() : email;
    const otp = String(data.get("otp") ?? "");
    const password = String(data.get("password") ?? "");

    setError("");
    setNotice("");
    if (step === "code") {
      const length = new TextEncoder().encode(password).length;
      if (length < 8 || length > 72) {
        setError("Use a password between 8 and 72 bytes.");
        return;
      }
      if (password !== String(data.get("confirm-password"))) {
        setError("The passwords do not match.");
        return;
      }
    }

    setBusy(true);
    try {
      if (step === "email") {
        await api<{ message: string }>(
          "/forgot-password",
          { method: "POST", body: JSON.stringify({ email: address }) },
          false,
        );
        setEmail(address);
        setStep("code");
      } else {
        await api<{ message: string }>(
          "/reset-password",
          { method: "POST", body: JSON.stringify({ email, otp, password }) },
          false,
        );
        setStep("done");
      }
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Please try again.");
    } finally {
      setBusy(false);
    }
  }

  return (
    <Card className="mx-auto w-full max-w-md p-8">
      <p className="mb-2 text-xs font-medium uppercase tracking-[0.12em] text-ink-40">
        Your reading space
      </p>
      <h2 className="text-[22px] font-semibold tracking-[-0.242px]">
        {step === "email"
          ? "Forgot your password?"
          : step === "code"
            ? "Check your inbox."
            : "Your password is updated."}
      </h2>
      <p className="mt-2 text-sm text-muted">
        {step === "email"
          ? "Enter your email and we’ll send a six-digit recovery code."
          : step === "code"
            ? "If this email is registered, a code is on its way. It expires in 10 minutes."
            : "Sign in with your new password to return to your reading space."}
      </p>

      {step !== "done" && (
        <form className="mt-8 space-y-5" onSubmit={submit}>
          {step === "email" ? (
            <div className="space-y-2">
              <Label htmlFor="recovery-email">Email address</Label>
              <Input
                id="recovery-email"
                name="email"
                type="email"
                autoComplete="email"
                placeholder="you@example.com"
                defaultValue={email}
                required
              />
            </div>
          ) : (
            <>
              <p className="text-sm text-muted">Code sent to {email}</p>
              <div className="space-y-2">
                <Label htmlFor="recovery-otp">Six-digit code</Label>
                <Input
                  id="recovery-otp"
                  ref={otpInput}
                  name="otp"
                  type="text"
                  inputMode="numeric"
                  autoComplete="one-time-code"
                  pattern="[0-9]{6}"
                  minLength={6}
                  maxLength={6}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="recovery-password">New password</Label>
                <Input
                  id="recovery-password"
                  name="password"
                  type="password"
                  autoComplete="new-password"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="recovery-confirm-password">
                  Confirm new password
                </Label>
                <Input
                  id="recovery-confirm-password"
                  name="confirm-password"
                  type="password"
                  autoComplete="new-password"
                  required
                />
              </div>
            </>
          )}
          {error && (
            <p
              role="alert"
              className="rounded-lg bg-canvas p-3 text-sm text-vermillion"
            >
              {error}
            </p>
          )}
          {notice && (
            <p role="status" className="text-sm text-muted">
              {notice}
            </p>
          )}
          <Button type="submit" className="w-full py-3" disabled={busy}>
            {busy
              ? "One moment…"
              : step === "email"
                ? "Send code"
                : "Reset password"}
            <ArrowRight size={16} />
          </Button>
          {step === "code" && (
            <div className="flex justify-center gap-6 text-sm font-medium text-notion-blue">
              <button
                type="button"
                className="hover:underline"
                onClick={resendCode}
                disabled={busy}
              >
                Resend code
              </button>
              <button
                type="button"
                className="hover:underline"
                onClick={() => {
                  setStep("email");
                  setError("");
                  setNotice("");
                }}
                disabled={busy}
              >
                Change email
              </button>
            </div>
          )}
        </form>
      )}
      <p className="mt-6 text-center text-sm text-muted">
        <Link
          className="font-medium text-notion-blue hover:underline"
          href="/login"
        >
          Back to sign in
        </Link>
      </p>
    </Card>
  );
}
