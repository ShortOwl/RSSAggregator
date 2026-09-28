import { expect, test } from "@playwright/test";

test("recovers a password with an emailed code", async ({ page }, testInfo) => {
  await page.addInitScript(() => localStorage.setItem("margin.theme", "dark"));
  const requests: { path: string; body: Record<string, string> }[] = [];
  await page.route("http://localhost:8080/v1/**", async (route) => {
    if (route.request().method() === "OPTIONS") {
      return route.fulfill({
        status: 204,
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Access-Control-Allow-Headers": "*",
          "Access-Control-Allow-Methods": "POST,OPTIONS",
        },
      });
    }
    const path = new URL(route.request().url()).pathname;
    const body = route.request().postDataJSON() as Record<string, string>;
    requests.push({ path, body });
    const status =
      path.endsWith("/reset-password") && body.otp !== "482913" ? 400 : 200;
    await route.fulfill({
      status,
      json:
        status === 400
          ? { error: "Invalid or expired code" }
          : {
              message:
                "If that email is registered, a recovery code will be sent.",
            },
      headers: { "Access-Control-Allow-Origin": "*" },
    });
  });

  await page.goto("/login");
  await page.getByRole("link", { name: "Forgot password?" }).click();
  await expect(page).toHaveURL(/forgot-password/);
  await page.getByLabel("Email address").fill("reader@example.com");
  await page.getByRole("button", { name: "Send code" }).click();
  await expect(
    page.getByRole("heading", { name: "Check your inbox." }),
  ).toBeVisible();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
  await expect(page.locator("[data-slot='card']")).toHaveCSS(
    "background-color",
    "rgb(36, 35, 33)",
  );
  await expect
    .poll(() =>
      page.evaluate(
        () => document.documentElement.scrollWidth <= window.innerWidth,
      ),
    )
    .toBe(true);
  await page.screenshot({
    path: `test-results/${testInfo.project.name}-password-recovery-dark.png`,
    fullPage: true,
  });
  expect(requests[0]).toEqual({
    path: "/v1/forgot-password",
    body: { email: "reader@example.com" },
  });

  await page.getByLabel("Six-digit code").fill("123456");
  await page.getByLabel("New password", { exact: true }).fill("NewPassword123");
  await page.getByLabel("Confirm new password").fill("NewPassword123");
  await page.getByRole("button", { name: "Reset password" }).click();
  await expect(page.locator("form [role='alert']")).toHaveText(
    "Invalid or expired code",
  );
  await expect(
    page.getByRole("heading", { name: "Check your inbox." }),
  ).toBeVisible();

  await page.getByLabel("Six-digit code").fill("482913");
  await page.getByRole("button", { name: "Reset password" }).click();
  await expect(
    page.getByRole("heading", { name: "Your password is updated." }),
  ).toBeVisible();
  expect(requests[2]).toEqual({
    path: "/v1/reset-password",
    body: {
      email: "reader@example.com",
      otp: "482913",
      password: "NewPassword123",
    },
  });
  await page.getByRole("link", { name: "Back to sign in" }).click();
  await expect(page).toHaveURL(/login/);
});

test("checks password confirmation before submitting and supports resending", async ({
  page,
}) => {
  const requests: string[] = [];
  let failNextResend = false;
  await page.route("http://localhost:8080/v1/**", async (route) => {
    if (route.request().method() === "OPTIONS") {
      return route.fulfill({
        status: 204,
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Access-Control-Allow-Headers": "*",
          "Access-Control-Allow-Methods": "POST,OPTIONS",
        },
      });
    }
    requests.push(new URL(route.request().url()).pathname);
    const failed = failNextResend;
    failNextResend = false;
    await route.fulfill({
      status: failed ? 503 : 200,
      json: failed
        ? { error: "Code delivery is unavailable. Please try again." }
        : { message: "Check your inbox." },
      headers: { "Access-Control-Allow-Origin": "*" },
    });
  });

  await page.goto("/forgot-password");
  await page.getByLabel("Email address").fill("reader@example.com");
  await page.getByRole("button", { name: "Send code" }).click();
  await page.getByLabel("Six-digit code").fill("482913");
  await page.getByLabel("New password", { exact: true }).fill("NewPassword123");
  await page.getByLabel("Confirm new password").fill("DifferentPassword");
  await page.getByRole("button", { name: "Reset password" }).click();
  await expect(page.locator("form [role='alert']")).toHaveText(
    "The passwords do not match.",
  );
  expect(requests).toEqual(["/v1/forgot-password"]);

  failNextResend = true;
  await page.getByRole("button", { name: "Resend code" }).click();
  await expect(page.locator("form [role='alert']")).toHaveText(
    "Code delivery is unavailable. Please try again.",
  );
  await expect(page.getByLabel("Six-digit code")).toHaveValue("482913");
  await page.getByRole("button", { name: "Resend code" }).click();
  await expect(page.getByRole("status")).toHaveText(
    "If this email is registered, a new code is on its way.",
  );
  await expect(page.getByLabel("Six-digit code")).toHaveValue("");
  await expect(page.getByLabel("Email address")).toHaveCount(0);
  expect(requests).toEqual([
    "/v1/forgot-password",
    "/v1/forgot-password",
    "/v1/forgot-password",
  ]);

  await page.getByRole("button", { name: "Change email" }).click();
  await expect(page.getByLabel("Email address")).toHaveValue(
    "reader@example.com",
  );
});
