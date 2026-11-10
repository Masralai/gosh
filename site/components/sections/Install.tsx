import InstallCommand from "@/components/InstallCommand";
import Container from "@/components/ui/Container";

export default function Install() {
  return (
    <section
      id="install"
      className="section-y border-t border-line/60"
      aria-labelledby="install-h"
    >
      <Container>
        <div className="mx-auto max-w-3xl text-center">
          <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-accent">
            Install
          </p>
          <h2
            id="install-h"
            className="mt-3 text-[clamp(1.9rem,3.4vw,2.75rem)] font-medium leading-[1.1] tracking-[-0.035em] text-ink text-balance"
          >
            Install GoSh
          </h2>
          <p className="mx-auto mt-4 max-w-md text-base leading-relaxed text-muted">
            One command, checksum checked, no sudo.
          </p>
        </div>

        <div className="mt-8">
          <InstallCommand />
        </div>
      </Container>
    </section>
  );
}
