"use client";

import { useEffect, useRef, useState } from "react";

const COMMANDS = {
  curl: ["curl -fsSL https://raw.githubusercontent.com/Masralai/gosh/main/install.sh | sh"],
  source: ["git clone https://github.com/Masralai/gosh", "cd gosh && make build"],
} as const;

type Method = keyof typeof COMMANDS;

const LABELS: Record<Method, string> = {
  curl: "Curl",
  source: "Source",
};

function pillClass(active: boolean) {
  return [
    "inline-flex min-h-11 cursor-pointer appearance-none items-center justify-center rounded-[4px] px-3.5 text-[13px] font-medium transition-colors sm:min-h-8",
    active ? "bg-raised text-ink" : "text-muted hover:text-ink",
  ].join(" ");
}

async function copyText(text: string) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text);
    return;
  }

  const area = document.createElement("textarea");
  area.value = text;
  area.setAttribute("readonly", "");
  area.style.position = "fixed";
  area.style.opacity = "0";
  document.body.appendChild(area);
  area.select();
  const ok = document.execCommand("copy");
  document.body.removeChild(area);
  if (!ok) throw new Error("copy failed");
}

export default function InstallCommand() {
  const [method, setMethod] = useState<Method>("curl");
  const [copied, setCopied] = useState<"idle" | "copied" | "error">("idle");
  const timer = useRef<number | null>(null);

  useEffect(
    () => () => {
      if (timer.current) window.clearTimeout(timer.current);
    },
    [],
  );

  const onCopy = async () => {
    try {
      await copyText(COMMANDS[method].join("\n"));
      setCopied("copied");
    } catch {
      setCopied("error");
    }

    if (timer.current) window.clearTimeout(timer.current);
    timer.current = window.setTimeout(() => setCopied("idle"), 2000);
  };

  return (
    <div className="mx-auto max-w-3xl">
      <div className="flex gap-1.5" role="group" aria-label="Install method">
        {(Object.keys(COMMANDS) as Method[]).map((key) => (
          <button
            key={key}
            type="button"
            aria-pressed={method === key}
            onClick={() => setMethod(key)}
            className={pillClass(method === key)}
          >
            {LABELS[key]}
          </button>
        ))}
      </div>

      <div className="mt-2.5 flex items-center gap-3 rounded-md bg-surface px-4 py-3 sm:px-5">
        <code className="min-w-0 flex-1 break-all font-mono text-[13px] leading-relaxed text-ink">
          {COMMANDS[method].map((line) => (
            <span key={line} className="block">
              {line}
            </span>
          ))}
        </code>

        <button
          type="button"
          onClick={onCopy}
          aria-label={
            method === "curl" ? "Copy curl install command" : "Copy source build commands"
          }
          className="inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-[4px] border border-line text-muted transition-colors hover:border-line-strong hover:text-ink active:translate-y-px sm:h-9 sm:w-9"
        >
          <svg width="15" height="15" viewBox="0 0 16 16" fill="none" aria-hidden="true">
            {copied === "copied" ? (
              <path
                d="M3.5 8.5l3 3 6-7"
                stroke="currentColor"
                strokeWidth="1.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            ) : (
              <>
                <rect
                  x="5"
                  y="5"
                  width="8"
                  height="9"
                  rx="1.5"
                  stroke="currentColor"
                  strokeWidth="1.3"
                />
                <path
                  d="M3.5 10.5V3.5A1.5 1.5 0 015 2h6.5"
                  stroke="currentColor"
                  strokeWidth="1.3"
                  strokeLinecap="round"
                />
              </>
            )}
          </svg>
          <span className="sr-only" aria-live="polite">
            {copied === "copied" ? "Copied" : copied === "error" ? "Copy failed" : "Copy"}
          </span>
        </button>
      </div>
    </div>
  );
}
