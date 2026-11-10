"use client";

import { useEffect } from "react";

const basePath = process.env.NEXT_PUBLIC_BASE_PATH || "";
const SOURCES = [
  `${basePath}/assets/system.js`,
  `${basePath}/assets/motion/motion.js`,
];

let started = false;

export default function KitScripts() {
  useEffect(() => {
    if (started) return;
    started = true;

    void (async () => {
      for (const src of SOURCES) {
        await new Promise<void>((resolve) => {
          const script = document.createElement("script");
          script.type = "module";
          script.src = src;
          script.async = false;
          script.onload = () => resolve();
          script.onerror = () => resolve();
          document.body.appendChild(script);
        });
      }
    })();
  }, []);

  return null;
}
