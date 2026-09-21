import { AuthGuard } from "@/components/auth/auth-guard";
import { AppHeader } from "@/components/layout/app-header";
export default function ReaderLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthGuard>
      <AppHeader />
      <main
        id="main"
        className="mx-auto min-h-screen max-w-[1440px] px-4 pb-16 pt-28 sm:px-8 lg:px-16"
      >
        {children}
      </main>
    </AuthGuard>
  );
}
