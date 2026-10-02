import { test, expect, type Page } from "@playwright/test";

const feed = {
  id: "feed-1",
  name: "The Curious Engineer",
  url: "https://example.com/feed.xml",
  user_id: "user-1",
  created_at: "2026-09-01T12:00:00Z",
  updated_at: "2026-09-01T12:00:00Z",
};
const first = {
  id: "post-1",
  title: "The quiet art of building things that last",
  description:
    "<p>What happens when we slow down and make room for thoughtful work? A reflection on craft, care, and the small decisions that shape our days.</p>",
  url: "https://example.com/craft",
  feed_id: feed.id,
  published_at: "2026-09-21T08:00:00Z",
  created_at: "2026-09-21T08:00:00Z",
  updated_at: "2026-09-21T08:00:00Z",
};
const second = {
  ...first,
  id: "post-2",
  title: "A field guide to following your curiosity",
  url: "https://example.com/curiosity",
};
const third = {
  ...first,
  id: "post-3",
  title: "Small rituals, wider perspectives",
  url: "https://example.com/rituals",
};
const jwt = `header.${Buffer.from(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 86400 })).toString("base64url")}.signature`;

async function setup(page: Page, authenticated = true) {
  const state = {
    bookmarked: false,
    read: false,
    followed: true,
    failBookmark: false,
    failFollow: false,
    created: 0,
    requests: [] as string[],
  };
  if (authenticated)
    await page.addInitScript(
      (token) => localStorage.setItem("margin.session", token),
      jwt,
    );
  await page.route("http://localhost:8080/v1/**", async (route) => {
    const request = route.request();
    const url = new URL(request.url());
    const path = url.pathname;
    const method = request.method();
    const json = (body: unknown, status = 200) =>
      route.fulfill({
        status,
        json: body,
        headers: { "Access-Control-Allow-Origin": "*" },
      });
    state.requests.push(method + " " + path + url.search);
    if (method === "OPTIONS")
      return route.fulfill({
        status: 204,
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Access-Control-Allow-Headers": "*",
          "Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
        },
      });
    if (path.endsWith("/login") || path.endsWith("/register"))
      return json({ token: jwt, token_type: "Bearer" });
    if (path.endsWith("/users"))
      return json({
        id: "user-1",
        name: "Alex Morgan",
        email: "alex@example.com",
        created_at: feed.created_at,
        updated_at: feed.updated_at,
      });
    if (path.endsWith("/feeds")) {
      if (method === "POST") {
        state.created++;
        return json({ ...feed, id: "feed-new" }, 201);
      }
      return json({ data: [feed], next_cursor: "" });
    }
    if (path.includes("/feed_follows")) {
      if (method === "DELETE") {
        state.followed = false;
        return route.fulfill({ status: 204 });
      }
      if (method === "POST") {
        if (state.failFollow) {
          state.failFollow = false;
          return json({ error: "Following failed" }, 500);
        }
        state.followed = true;
        return json({ id: "follow-1", feed_id: feed.id }, 201);
      }
      return json(state.followed ? [{ id: "follow-1", feed_id: feed.id }] : []);
    }
    if (path.endsWith("/bookmark")) {
      if (state.failBookmark)
        return json({ error: "Couldn’t save bookmark" }, 500);
      state.bookmarked = method === "PUT";
      return route.fulfill({ status: 204 });
    }
    if (path.endsWith("/read")) {
      state.read = method === "PUT";
      return route.fulfill({ status: 204 });
    }
    if (path.endsWith("/bookmarks"))
      return json({ data: state.bookmarked ? [first] : [], next_cursor: "" });
    if (path.endsWith("/posts")) {
      if (url.searchParams.get("unread") === "true")
        return json({
          data: state.read ? [second, third] : [first, second, third],
          next_cursor: "",
        });
      if (url.searchParams.get("search"))
        return json({ data: [second], next_cursor: "" });
      if (url.searchParams.get("cursor"))
        return json({ data: [third], next_cursor: "" });
      return json({ data: [first, second], next_cursor: "page-two" });
    }
    return json({ error: "Unexpected API request" }, 404);
  });
  return state;
}

test("signs in and restores the persistent session", async ({ page }) => {
  await setup(page, false);
  await page.goto("/");
  await expect(page).toHaveURL(/login/);
  await page.getByLabel("Email address").fill("alex@example.com");
  await page.getByLabel("Password", { exact: true }).fill("password123");
  await page.getByRole("button", { name: "Sign in", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "Your daily reading." }),
  ).toBeVisible();
  await page.reload();
  await expect(
    page.getByRole("heading", { name: "Your daily reading." }),
  ).toBeVisible();
});

test("paginates, searches, filters, and persists reading actions", async ({
  page,
}) => {
  const state = await setup(page);
  await page.goto("/");
  const article = page.getByRole("article").filter({ hasText: first.title });
  await expect(
    article.getByRole("button", { name: "Bookmark post", exact: true }),
  ).toBeVisible();
  await article
    .getByRole("button", { name: "Bookmark post", exact: true })
    .click();
  await expect(
    article.getByRole("button", { name: "Remove bookmark", exact: true }),
  ).toBeVisible();
  await article
    .getByRole("button", { name: "Mark as read", exact: true })
    .click();
  await expect(
    article.getByRole("button", { name: "Mark as unread", exact: true }),
  ).toBeVisible();
  await page.getByRole("article").last().scrollIntoViewIfNeeded();
  await expect(page.getByRole("heading", { name: third.title })).toBeVisible();
  await page
    .getByRole("combobox", { name: "Filter by feed" })
    .selectOption(feed.id);
  await page.getByRole("textbox", { name: "Search posts" }).fill("curiosity");
  await expect(page.getByRole("article")).toHaveCount(1);
  expect(
    state.requests.some(
      (request) =>
        request.includes("search=curiosity") &&
        request.includes("feed_id=feed-1"),
    ),
  ).toBeTruthy();
  await page.goto("/bookmarks");
  await expect(page.getByRole("heading", { name: first.title })).toBeVisible();
  await page
    .getByRole("button", { name: "Remove bookmark", exact: true })
    .click();
  await expect(
    page.getByRole("heading", { name: "Keep a little inspiration." }),
  ).toBeVisible();
});

test("does not show success or change state when a mutation fails", async ({
  page,
}) => {
  const state = await setup(page);
  state.failBookmark = true;
  await page.goto("/");
  const article = page.getByRole("article").filter({ hasText: first.title });
  await article
    .getByRole("button", { name: "Bookmark post", exact: true })
    .click();
  await expect(
    page.getByText("Couldn’t save bookmark", { exact: true }),
  ).toBeVisible();
  await expect(
    article.getByRole("button", { name: "Bookmark post", exact: true }),
  ).toBeVisible();
});

test("unfollows with subscription ID and recovers partial feed creation", async ({
  page,
}) => {
  const state = await setup(page);
  await page.goto("/feeds");
  await page.getByRole("button", { name: `Unfollow ${feed.name}` }).click();
  await expect(
    page.getByRole("heading", { name: "Every good habit starts with one." }),
  ).toBeVisible();
  expect(state.requests).toContain("DELETE /v1/feed_follows/follow-1");
  state.failFollow = true;
  await page.getByRole("button", { name: "Add a feed", exact: true }).click();
  await page.getByLabel("Feed name").fill("A new voice");
  await page.getByLabel("Feed URL").fill("https://example.com/new.xml");
  await page.getByRole("button", { name: "Add and follow" }).click();
  await expect(
    page.getByText("Your feed was created.", { exact: false }),
  ).toBeVisible();
  await page.getByRole("button", { name: "Retry following" }).click();
  await expect(page.getByRole("dialog")).toHaveCount(0);
  expect(state.created).toBe(1);
});

test("saves compact preference and signs out", async ({ page }) => {
  await setup(page);
  await page.goto("/settings");
  await expect(page.getByText("Alex Morgan", { exact: true })).toBeVisible();
  await page.getByRole("checkbox", { name: /Compact stories/ }).check();
  await page.goto("/");
  await expect(page.getByRole("heading", { name: first.title })).toBeVisible();
  await expect(page.getByText(/What happens when we slow down/)).toHaveCount(0);
  await page.goto("/settings");
  await page.getByRole("button", { name: "Sign out", exact: true }).click();
  await expect(page).toHaveURL(/login/);
});

test("renders responsive pages with no horizontal overflow", async ({
  page,
}, testInfo) => {
  await setup(page);
  for (const path of ["/", "/feeds", "/bookmarks", "/settings"]) {
    await page.goto(path);
    await expect(page.getByRole("heading", { level: 1 })).toBeVisible();
    await expect
      .poll(() =>
        page.evaluate(
          () => document.documentElement.scrollWidth <= window.innerWidth,
        ),
      )
      .toBe(true);
  }
  await page.goto("/");
  await expect(page.getByRole("heading", { name: first.title })).toBeVisible();
  await page.screenshot({
    path: `test-results/${testInfo.project.name}-home.png`,
    fullPage: true,
  });
});

test("expired sessions return to login", async ({ page }) => {
  await setup(page, false);
  await page.addInitScript(() =>
    localStorage.setItem(
      "margin.session",
      `header.${btoa(JSON.stringify({ exp: 1 }))}.signature`,
    ),
  );
  await page.goto("/bookmarks");
  await expect(page).toHaveURL(/login/);
});

test("registers and renders the authentication screen", async ({
  page,
}, testInfo) => {
  await setup(page, false);
  await page.goto("/register");
  await expect(
    page.getByRole("heading", { name: "A fresh page starts here." }),
  ).toBeVisible();
  await expect
    .poll(() =>
      page.evaluate(
        () => document.documentElement.scrollWidth <= window.innerWidth,
      ),
    )
    .toBe(true);
  await page.screenshot({
    path: `test-results/${testInfo.project.name}-register.png`,
    fullPage: true,
  });
  await page.getByLabel("Your name").fill("Alex Morgan");
  await page.getByLabel("Email address").fill("alex@example.com");
  await page.getByLabel("Password", { exact: true }).fill("password123");
  await page
    .getByRole("button", { name: "Create account", exact: true })
    .click();
  await expect(
    page.getByRole("heading", { name: "Your daily reading." }),
  ).toBeVisible();
});

test("rejects server-invalid sessions", async ({ page }) => {
  await setup(page);
  await page.route("http://localhost:8080/v1/posts**", (route) =>
    route.fulfill({
      status: 401,
      json: { error: "token is expired or invalid" },
    }),
  );
  await page.goto("/");
  await expect(page).toHaveURL(/login/);
});

test("theme choices persist without changing compact reading or notifications", async ({
  page,
}) => {
  await setup(page);
  await page.emulateMedia({ colorScheme: "light" });
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  page.on("console", (message) => {
    if (
      message.type() === "error" &&
      /hydration|did not match/i.test(message.text())
    )
      errors.push(message.text());
  });
  await page.goto("/settings");
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
  await expect(page.locator("body")).toHaveCSS(
    "background-color",
    "rgb(246, 245, 244)",
  );
  await expect(page.getByLabel("Appearance", { exact: true })).toHaveValue(
    "system",
  );
  await page.getByRole("checkbox", { name: /Compact stories/ }).check();
  await page.getByLabel("Appearance", { exact: true }).selectOption("dark");
  await expect(page.locator("body")).toHaveCSS(
    "background-color",
    "rgb(25, 24, 23)",
  );
  await expect(page.locator("[data-sonner-toast]").first()).toHaveCSS(
    "background-color",
    "rgb(255, 255, 255)",
  );
  await page.reload();
  await expect(page.getByLabel("Appearance", { exact: true })).toHaveValue(
    "dark",
  );
  await expect(
    page.getByRole("checkbox", { name: /Compact stories/ }),
  ).toBeChecked();
  await page.getByRole("link", { name: "Margin home", exact: true }).click();
  await expect(page.getByRole("heading", { name: first.title })).toBeVisible();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
  await page.getByRole("button", { name: "Open account menu" }).click();
  await page.getByRole("menuitemradio", { name: "Light", exact: true }).click();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
  await page.emulateMedia({ colorScheme: "dark" });
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
  await page.getByRole("button", { name: "Open account menu" }).click();
  await page
    .getByRole("menuitemradio", { name: "System", exact: true })
    .click();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
  await page.emulateMedia({ colorScheme: "light" });
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
  expect(errors).toEqual([]);
});

test("dark mode applies before hydration and follows the device by default", async ({
  page,
}) => {
  await page.emulateMedia({ colorScheme: "dark" });
  // Hold the application scripts: the head script must style server-rendered content on its own.
  await page.route("**/_next/**/*.js*", (route) => route.abort());
  await page.goto("/login");
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
  await expect(page.locator("body")).toHaveCSS(
    "background-color",
    "rgb(25, 24, 23)",
  );
  await page.evaluate(() => localStorage.setItem("margin.theme", "light"));
  await page.reload();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
});

test("theme switching works when theme storage is blocked", async ({
  page,
}) => {
  await setup(page);
  await page.emulateMedia({ colorScheme: "dark" });
  await page.addInitScript(() => {
    const get = Storage.prototype.getItem;
    const set = Storage.prototype.setItem;
    Storage.prototype.getItem = function (key) {
      if (key === "margin.theme")
        throw new DOMException("Blocked", "SecurityError");
      return get.call(this, key);
    };
    Storage.prototype.setItem = function (key, value) {
      if (key === "margin.theme")
        throw new DOMException("Blocked", "SecurityError");
      return set.call(this, key, value);
    };
  });
  await page.goto("/settings");
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
  await page.getByLabel("Appearance", { exact: true }).selectOption("light");
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
  await page.getByRole("link", { name: "Margin home", exact: true }).click();
  await expect(page.getByRole("heading", { name: first.title })).toBeVisible();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
});

test("both themes preserve accent colors and responsive surfaces", async ({
  page,
}, testInfo) => {
  await setup(page);
  for (const theme of ["light", "dark"] as const) {
    await page.emulateMedia({ colorScheme: theme });
    for (const path of ["/", "/feeds", "/bookmarks", "/settings", "/login"]) {
      await page.goto(path);
      await expect(page.locator("html")).toHaveAttribute("data-theme", theme);
      await expect(page.getByRole("heading", { level: 1 })).toBeVisible();
      await expect
        .poll(() =>
          page.evaluate(
            () => document.documentElement.scrollWidth <= innerWidth,
          ),
        )
        .toBe(true);
      const accent = page.locator(".bg-marigold").first();
      await expect(accent).toHaveCSS("background-color", "rgb(255, 177, 16)");
      await expect(accent).toHaveCSS("color", "rgb(0, 0, 0)");
      await page.screenshot({
        path: `test-results/${testInfo.project.name}-${theme}-${path.slice(1) || "home"}.png`,
        fullPage: true,
      });
    }
    await page.goto("/feeds");
    await page.getByRole("button", { name: "Add a feed", exact: true }).click();
    await expect(page.getByRole("dialog")).toHaveCSS(
      "background-color",
      theme === "dark" ? "rgb(36, 35, 33)" : "rgb(255, 255, 255)",
    );
    await page.screenshot({
      path: `test-results/${testInfo.project.name}-${theme}-dialog.png`,
      fullPage: true,
    });
  }
});

test("header theme toggle resolves System and persists across navigation and reload", async ({
  page,
}) => {
  await setup(page);
  await page.emulateMedia({ colorScheme: "dark" });
  await page.goto("/");
  await page
    .getByRole("button", { name: "Switch to light mode", exact: true })
    .click();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
  await page.reload();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");
  await page
    .getByRole("button", { name: "Switch to dark mode", exact: true })
    .click();
  await page.getByRole("link", { name: "Settings and profile" }).click();
  await expect(page.getByLabel("Appearance", { exact: true })).toHaveValue(
    "dark",
  );
  await page.getByLabel("Appearance", { exact: true }).selectOption("system");
  await page.emulateMedia({ colorScheme: "light" });
  await expect(
    page.getByRole("button", { name: "Switch to dark mode", exact: true }),
  ).toBeVisible();
});
