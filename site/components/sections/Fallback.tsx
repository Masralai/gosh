import Container from "@/components/ui/Container";

export default function Fallback() {
  return (
    <section
      className="section-y border-t border-line/60"
      data-moment="word-highlight"
      aria-label="The bash fallback"
    >
      <Container>
        <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-accent">Fallback</p>
        <p className="mo-wh__text mt-4 max-w-[960px] text-ink">
          GoSh answers 27 commands itself. Everything else goes to <strong>bash</strong>, with the
          same arguments, in the same directory, under the same prompt.
        </p>
      </Container>
    </section>
  );
}
