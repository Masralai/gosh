"use client";

import { useState } from "react";

import Container from "@/components/ui/Container";
import Mark from "@/components/ui/Mark";

const links = [
  { href: "https://github.com/Masralai/gosh#command-reference", label: "Docs" },
  { href: "#install", label: "Install" },
];

export default function Nav() {
  const [open, setOpen] = useState(false);

  return (
    <header className="fixed inset-x-0 top-0 z-50 border-b border-line/80 bg-page/85 backdrop-blur-md">
      <Container className="flex h-16 items-center gap-8 lg:h-20">
        <a href="#top" className="flex items-center gap-2.5" aria-label="GoSh home">
          <Mark className="h-7 w-7 shrink-0" />
          <span className="text-lg font-semibold tracking-[-0.04em] text-ink">GoSh</span>
        </a>

        <nav className="hidden items-center gap-7 md:flex" aria-label="Main">
          {links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="text-sm font-medium text-muted transition-colors hover:text-ink"
            >
              {link.label}
            </a>
          ))}
        </nav>

        <div className="ml-auto flex items-center gap-3">
          <a
            href="#install"
            className="btn btn--primary h-10 px-4"
            onClick={() => setOpen(false)}
          >
            Install GoSh
          </a>
          <button
            type="button"
            onClick={() => setOpen((value) => !value)}
            aria-expanded={open}
            aria-controls="mobile-menu"
            aria-label={open ? "Close menu" : "Open menu"}
            className="inline-flex h-10 w-10 items-center justify-center rounded-[4px] border border-line text-ink transition-colors hover:border-line-strong md:hidden"
          >
            <svg width="18" height="18" viewBox="0 0 18 18" aria-hidden="true">
              {open ? (
                <path
                  d="M4 4l10 10M14 4L4 14"
                  stroke="currentColor"
                  strokeWidth="1.6"
                  strokeLinecap="round"
                />
              ) : (
                <path
                  d="M2 5h14M2 13h14"
                  stroke="currentColor"
                  strokeWidth="1.6"
                  strokeLinecap="round"
                />
              )}
            </svg>
          </button>
        </div>
      </Container>

      {open ? (
        <div id="mobile-menu" className="border-t border-line bg-page md:hidden">
          <Container className="flex flex-col py-3">
            {links.map((link) => (
              <a
                key={link.href}
                href={link.href}
                onClick={() => setOpen(false)}
                className="rounded-[4px] px-2 py-3 text-base font-medium text-muted transition-colors hover:bg-surface hover:text-ink"
              >
                {link.label}
              </a>
            ))}
          </Container>
        </div>
      ) : null}
    </header>
  );
}
