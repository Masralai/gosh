import type { ReactNode } from "react";

export default function SectionHeading({
  eyebrow,
  title,
  id,
  children,
}: {
  eyebrow: string;
  title: string;
  id?: string;
  children?: ReactNode;
}) {
  return (
    <div className="max-w-2xl">
      <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-accent">
        {eyebrow}
      </p>
      <h2
        id={id}
        className="mt-3 text-[clamp(1.9rem,3.4vw,2.75rem)] font-medium leading-[1.1] tracking-[-0.035em] text-ink text-balance"
      >
        {title}
      </h2>
      {children ? (
        <div className="mt-4 text-base leading-relaxed text-muted sm:text-lg">{children}</div>
      ) : null}
    </div>
  );
}
