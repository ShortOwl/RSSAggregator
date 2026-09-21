import Link from "next/link";
import { BrandMark } from "@/components/shared/brand-mark";
export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <main className="mx-auto min-h-screen max-w-[1440px] px-4 py-8 sm:px-8 lg:px-16">
      <Link href="/" aria-label="Margin home">
        <BrandMark />
      </Link>
      <div className="grid items-center gap-16 py-16 lg:grid-cols-2 lg:py-20">
        <section className="mx-auto max-w-xl text-center lg:text-left">
          <div className="mb-8 inline-flex items-center gap-2 rounded-full bg-white px-3 py-1 text-xs text-graphite">
            <span className="size-2 rounded-full bg-notion-blue" />A little less
            noise. A lot more perspective.
          </div>
          <h1 className="text-[54px] font-semibold leading-[1.04] tracking-[-1.89px] sm:text-[72px] sm:leading-[1.21] sm:tracking-[-2.016px]">
            Make room
            <br />
            for{" "}
            <span className="inline-block rounded-full bg-marigold px-6 py-2">
              curiosity.
            </span>
          </h1>
          <p className="mx-auto mt-8 max-w-md font-serif text-lg leading-[1.56] text-graphite lg:mx-0">
            Your favorite voices, all in one quiet place. Follow what matters.
            Save what stays with you.
          </p>
          <div
            aria-hidden="true"
            className="mt-9 flex justify-center gap-3 lg:justify-start"
          >
            {["#ffb110", "#62aef0", "#f64932"].map((color, i) => (
              <svg key={color} viewBox="0 0 48 48" width="48" height="48">
                <circle
                  cx="24"
                  cy="24"
                  r="22"
                  fill="white"
                  stroke={color}
                  strokeWidth="2"
                />
                <path
                  d={i === 1 ? "M12 22q12-19 24 0" : "M12 19q12-12 24 0"}
                  fill={color}
                />
                <circle cx="19" cy="24" r="1.5" />
                <circle cx="29" cy="24" r="1.5" />
                <path
                  d="M19 31q5 4 10 0"
                  fill="none"
                  stroke="black"
                  strokeWidth="1.5"
                />
              </svg>
            ))}
          </div>
        </section>
        {children}
      </div>
      <p className="text-center text-xs text-black/40">
        A reading ritual, made your own.
      </p>
    </main>
  );
}
