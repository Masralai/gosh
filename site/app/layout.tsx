import type { Metadata, Viewport } from "next";
import type { ReactNode } from "react";
import Script from "next/script";

import "./kit/motion.css";
import "./globals.css";

import KitScripts from "@/components/KitScripts";

const basePath = process.env.NEXT_PUBLIC_BASE_PATH || "";
const siteUrl = process.env.NEXT_PUBLIC_SITE_URL || "https://masralai.github.io";

const title = "GoSh: a shell written in Go";
const description =
  "GoSh is a shell written in Go with 27 built-in commands for files, processes, memory, search and archives.";

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title,
  description,
  openGraph: {
    type: "website",
    title,
    description:
      "27 built-in commands for files, processes, memory, search and archives. One binary for Windows, macOS and Linux.",
    images: [`${basePath}/share.png`],
  },
  twitter: { card: "summary_large_image" },
  icons: { icon: `${basePath}/favicon.png` },
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  viewportFit: "cover",
  themeColor: "#202224",
};

const guards = `document.documentElement.classList.add('m-enter');setTimeout(()=>document.documentElement.classList.add('m-enter-done'),700);document.documentElement.classList.add('mo-js');setTimeout(()=>document.documentElement.classList.add('mo-safe'),2500)`;

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html
      lang="en"
      data-type="software"
      data-radius="sharp"
      data-buttons="pill"
      data-tone="dark"
      data-motion=""
    >
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="" />
        <link
          href="https://fonts.googleapis.com/css2?family=Manrope:wght@400..700&family=JetBrains+Mono:wght@400..500&display=swap"
          rel="stylesheet"
          precedence="default"
        />
        <link
          rel="stylesheet"
          href={`${basePath}/assets/moments/word-highlight/moment.css`}
          precedence="default"
        />
        <link rel="modulepreload" href={`${basePath}/assets/moments/moments.js`} />
        <link rel="modulepreload" href={`${basePath}/assets/moments/lib.js`} />
        <link
          rel="modulepreload"
          href={`${basePath}/assets/moments/word-highlight/moment.js`}
        />
        <Script
          id="kit-guards"
          strategy="beforeInteractive"
          dangerouslySetInnerHTML={{ __html: guards }}
        />
      </head>
      <body className="min-h-screen bg-page text-ink antialiased">
        {children}
        <KitScripts />
      </body>
    </html>
  );
}
