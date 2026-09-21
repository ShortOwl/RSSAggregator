export function PageHeader({
  eyebrow,
  title,
  description,
  children,
}: {
  eyebrow: string;
  title: string;
  description: string;
  children?: React.ReactNode;
}) {
  return (
    <div className="mb-9 flex flex-wrap items-end justify-between gap-6">
      <div>
        <p className="mb-3 text-xs font-medium uppercase tracking-[0.12em] text-black/40">
          {eyebrow}
        </p>
        <h1 className="text-[40px] font-semibold leading-[1.14] tracking-[-0.036em] sm:text-[54px] sm:leading-[1.04]">
          {title}
        </h1>
        <p className="mt-4 font-serif text-lg leading-[1.56] text-graphite">
          {description}
        </p>
      </div>
      {children}
    </div>
  );
}
