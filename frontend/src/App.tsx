import { useRef } from "react";
import Navbar from "./components/ui/Navbar";
import Hero from "./components/sections/Hero";
import Packages from "./components/sections/Packages";
import Menu from "./components/sections/Menu";
import BookingModal from "./components/ui/BookingModal";
import { BookingProvider, useBooking } from "./context/BookingContext";
import { MenuProvider } from "./context/MenuContext";
import { Parallax, ParallaxLayer, IParallax } from "@react-spring/parallax";

const AppContent = () => {
  const { isModalOpen, selectedPackage, closeModal } = useBooking();
  const parallax = useRef<IParallax>(null!);
  return (
    <Parallax ref={parallax} pages={3}>
      <div className="min-h-screen">
        <ParallaxLayer offset={0} speed={0}>
          <Navbar />
          <Hero />
        </ParallaxLayer>
        <ParallaxLayer offset={1} speed={0}>
          <Packages />
          <BookingModal
            isOpen={isModalOpen}
            onClose={closeModal}
            selectedPackage={selectedPackage}
          />
        </ParallaxLayer>

        <ParallaxLayer offset={2} speed={0}>
          <Menu />
        </ParallaxLayer>
        {/* Other sections will be added here later */}
      </div>
    </Parallax>
  );
};

function App() {
  return (
    <MenuProvider>
      <BookingProvider>
        <AppContent />
      </BookingProvider>
    </MenuProvider>
  );
}

export default App;
