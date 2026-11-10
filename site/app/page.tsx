import Nav from "@/components/sections/Nav";
import Hero from "@/components/sections/Hero";
import Fallback from "@/components/sections/Fallback";
import Install from "@/components/sections/Install";
import Footer from "@/components/sections/Footer";

export default function Page() {
  return (
    <>
      <Nav />
      <main id="top">
        <Hero />
        <Fallback />
        <Install />
      </main>
      <Footer />
    </>
  );
}
