import Container from "@/components/ui/Container";

const banner = String.raw`       ^  ^  ^   ^      ___I_      ^  ^   ^  ^  ^   ^  ^
      /|\/|\/|\ /|\    /\-_--\    /|\/|\ /|\/|\/|\ /|\/|\
      /|\/|\/|\ /|\   /  \_-__\   /|\/|\ /|\/|\/|\ /|\/|\
      /|\/|\/|\ /|\   |[]| [] |   /|\/|\ /|\/|\/|\ /|\/|\ `;

export default function Hero() {
  return (
    <section
      className="hero relative overflow-hidden pb-20 pt-28 sm:pb-24 sm:pt-36"
      data-enter=""
      aria-labelledby="hero-h"
    >
      <div className="hero-glow pointer-events-none absolute inset-x-0 top-0 h-[560px]" aria-hidden="true" />

      <Container className="relative">
        <div className="mx-auto max-w-3xl text-center">
          <h1
            id="hero-h"
            className="text-[clamp(2.6rem,6.2vw,4.6rem)] font-medium leading-[1.02] tracking-[-0.045em] text-ink text-balance"
          >
            A shell written in Go
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-base leading-relaxed text-muted sm:text-lg">
            27 built-in commands for files, processes, memory, search and archives. Built for
            people who already live in a terminal.
          </p>
          <div className="mt-8 flex justify-center">
            <a href="#install" className="btn btn--primary">
              Install GoSh
            </a>
          </div>
        </div>

        <div
          className="hero__window mx-auto mt-14 max-w-5xl overflow-hidden rounded-lg border border-line bg-term shadow-2xl shadow-black/50"
          data-visual=""
          role="img"
          aria-label="A GoSh session: ls -a lists the project, mu shows memory use, and grep -r finds Readline in cmd/gosh/main.go"
        >
          <div
            className="flex items-center gap-3 border-b border-white/10 bg-term-bar px-4 py-3 text-xs font-medium text-term-ink"
            aria-hidden="true"
          >
            <span className="flex gap-1.5">
              <i className="h-2.5 w-2.5 rounded-full bg-accent/70" />
              <i className="h-2.5 w-2.5 rounded-full bg-accent/45" />
              <i className="h-2.5 w-2.5 rounded-full bg-accent/25" />
            </span>
            <span>gosh — ~/code/gosh</span>
          </div>

          <div className="px-4 py-5 font-mono text-[12px] leading-[1.7] text-term-ink sm:px-6 sm:py-6 sm:text-[13px]">
            <pre className="mb-3 hidden whitespace-pre text-[11px] leading-[1.12] text-term-prompt sm:block">
              {banner}
            </pre>
            <div className="mb-2 text-term-prompt sm:hidden">GoSh · ~/code/gosh</div>

            <div className="space-y-1 whitespace-pre-wrap break-words" aria-hidden="true">
              <div className="text-term-out">Welcome to GoSh</div>
              <div>
                <span className="text-term-prompt">~/code/gosh ❯</span>{" "}
                <span className="text-ink">ls -a</span>
              </div>
              <div className="text-term-out">
                .github/  cmd/  internal/  site/  go.mod  Makefile  Readme.md  install.sh
              </div>
              <div>
                <span className="text-term-prompt">~/code/gosh ❯</span>{" "}
                <span className="text-ink">mu</span>
              </div>
              <div className="text-term-out">
                Memory  15.7 GB total  1.8 GB free  88% used  █████████░
              </div>
              <div>
                <span className="text-term-prompt">~/code/gosh ❯</span>{" "}
                <span className="text-ink">grep -r &quot;Readline&quot; cmd</span>
              </div>
              <div className="text-term-out">cmd/gosh/main.go:45: line, err := rl.Readline()</div>
            </div>
          </div>
        </div>
      </Container>
    </section>
  );
}
