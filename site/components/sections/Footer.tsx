import Container from "@/components/ui/Container";

const links = [
  { href: "https://github.com/Masralai/gosh", label: "GitHub" },
  { href: "https://github.com/Masralai/gosh#command-reference", label: "Docs" },
  { href: "https://github.com/Masralai/gosh/releases", label: "Releases" },
];

export default function Footer() {
  return (
    <footer
      className="border-t border-line bg-surface text-muted"
      id="contact"
      aria-label="Site footer"
    >
      <Container className="flex flex-wrap items-center justify-between gap-x-6 gap-y-2 py-6">
        <nav aria-label="Footer" className="flex flex-wrap items-center gap-x-6">
          {links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="inline-flex items-center py-3 text-sm font-medium transition-colors hover:text-ink"
            >
              {link.label}
            </a>
          ))}
        </nav>

        <p className="text-xs">© 2026 GoSh</p>
      </Container>
    </footer>
  );
}
