import type { Metadata } from "next";
import "@fontsource/inter/400.css";
import "@fontsource/inter/500.css";
import "@fontsource/inter/600.css";
import "@fontsource/inter/700.css";
import "@fontsource/source-serif-pro/400.css";
import "./globals.css";
import { Providers } from "./providers";
import { themeScript } from "@/lib/theme";
export const metadata: Metadata = {
  title: {
    default: "Margin — Your daily reading space",
    template: "%s · Margin",
  },
  description:
    "A quieter place for the stories you follow. Read, collect, and make room for what matters.",
};
export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <script dangerouslySetInnerHTML={{ __html: themeScript }} />
      </head>
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
