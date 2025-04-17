import React from "react";
import { useBooking } from "../../context/BookingContext";
import { scrollToSection } from "../../utils/scrollUtils";

const Hero: React.FC = () => {
  const { openModal } = useBooking();

  return (
    <div className="relative h-screen w-full bg-gradient-to-r from-caramel to-terracotta overflow-hidden">
      {/* Content container */}
      <div className="relative z-10 h-full flex flex-col items-center justify-center text-center px-4 md:px-8">
        <h1 className="text-5xl md:text-7xl font-bold text-mocha mb-4 tracking-tight">
          TOASTED COFFEE CO
        </h1>
        <p className="text-espresso text-4xl md:text-5xl max-w-2xl mb-8">
          COLD BREW BAR
        </p>
        <div className="flex flex-col sm:flex-row gap-4">
          <button
            className="px-8 py-3 bg-parchment text-mocha font-semibold rounded-full hover:bg-latte transition-all shadow-lg"
            onClick={() => scrollToSection("packages")}
          >
            View Packages
          </button>
          <button
            className="px-8 py-3 bg-transparent text-parchment font-semibold border-2 border-parchment rounded-full hover:bg-latte hover:text-mocha transition-all"
            onClick={() => openModal()}
          >
            Book Now
          </button>
        </div>
      </div>
    </div>
  );
};

export default Hero;
