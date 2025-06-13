import React, { useState, useEffect } from "react";
import { useBooking } from "../../context/BookingContext";
import { scrollToSection } from "../../utils/scrollUtils";
import heroSpillTop from "../../assets/hero-spill-top.png";
import heroSpillBottom from "../../assets/hero-spill-bottom.png";
import heroSpillLeft from "../../assets/hero-spill-left.png";
import heroCircleBottom from "../../assets/hero-circle-bottom.png";
import heroCircleLeft from "../../assets/hero-circle-left.png";
import heroRibbonsTop from "../../assets/hero-ribbons-top.png";
import heroRibbonsBottom from "../../assets/hero-ribbons-bottom.png";

const Hero: React.FC = () => {
  const { openModal } = useBooking();
  const [animated, setAnimated] = useState(false);

  // Trigger animation when component mounts
  useEffect(() => {
    // Short delay to ensure everything is loaded
    const timer = setTimeout(() => {
      setAnimated(true);
    }, 300);

    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="relative h-screen">
      {/* Background Elements z-index 1-9*/}
      <div className="absolute inset-0 z-[1] bg-peach"></div>

      {/* Animations and Decorative Elements z-index 10-19*/}

      {/* Circle animations */}
      <div
        className={`absolute bottom-0 right-0 h-auto z-[10] transition-all duration-1000 ease-out delay-150 ${
          animated
            ? "translate-x-0 translate-y-0"
            : "translate-x-full translate-y-full"
        }`}
      >
        <img
          src={heroCircleBottom}
          alt=""
          className="w-full h-auto"
          aria-hidden="true"
        />
      </div>
      <div
        className={`absolute left-0 h-auto z-[10] transition-all duration-1000 ease-out delay-150 ${
          animated ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <img
          src={heroCircleLeft}
          alt=""
          className="w-full h-auto"
          aria-hidden="true"
        />
      </div>
      {/* Spill animations */}
      <div
        className={`absolute top-0 left-0 right-0 w-full z-[15] transition-transform duration-1000 ease-out ${
          animated ? "translate-y-0" : "-translate-y-full"
        }`}
      >
        <img
          src={heroSpillTop}
          alt=""
          className="w-full h-auto object-cover"
          aria-hidden="true"
        />
      </div>
      <div
        className={`absolute bottom-0 left-0 right-0 w-full z-[15] transition-transform duration-1000 ease-out ${
          animated ? "translate-y-0" : "translate-y-full"
        }`}
      >
        <img
          src={heroSpillBottom}
          alt=""
          className="w-full h-auto object-cover"
          aria-hidden="true"
        />
      </div>

      <div
        className={`absolute top-0 left-0 h-auto z-[15] transition-all duration-1000 ease-out delay-150 ${
          animated
            ? "translate-x-0 translate-y-0"
            : "-translate-x-full -translate-y-full"
        }`}
      >
        <img
          src={heroSpillLeft}
          alt=""
          className="w-full h-auto"
          aria-hidden="true"
        />
      </div>
      {/* Ribbon animations */}
      <div
        className={`absolute top-0 left-0 right-0 w-full z-[19] transition-transform duration-1000 ease-out delay-300 ${
          animated ? "translate-y-0" : "-translate-y-full"
        }`}
      >
        <img
          src={heroRibbonsTop}
          alt=""
          className="w-full h-auto object-cover"
          aria-hidden="true"
        />
      </div>
      <div
        className={`absolute bottom-0 left-0 right-0 w-full z-[19] transition-transform duration-1000 ease-out delay-300 ${
          animated ? "translate-y-0" : "translate-y-full"
        }`}
      >
        <img
          src={heroRibbonsBottom}
          alt=""
          className="w-full h-auto object-cover"
          aria-hidden="true"
        />
      </div>
      {/* END: Animations and Decorative Elements z-index 10-29+*/}

      {/* Content container z-index=50+*/}
      <div className="relative z-[50] h-full flex flex-col items-center justify-center text-center px-4 md:px-8">
        <h1
          className={`text-5xl md:text-7xl font-bold text-mocha mb-4 tracking-tight transition-opacity duration-1000 delay-500 ${
            animated ? "opacity-100" : "opacity-0"
          }`}
        >
          TOASTED COFFEE CO
        </h1>
        <p
          className={`text-espresso text-4xl md:text-5xl max-w-2xl mb-8 transition-opacity duration-1000 delay-700 ${
            animated ? "opacity-100" : "opacity-0"
          }`}
        >
          COLD BREW BAR
        </p>
        <div
          className={`flex flex-col sm:flex-row gap-4 transition-opacity duration-1000 delay-900 ${
            animated ? "opacity-100" : "opacity-0"
          }`}
        >
          <button
            className="px-8 py-3 bg-parchment text-mocha font-semibold rounded-full hover:bg-latte transition-all shadow-lg"
            onClick={() => scrollToSection("packages")}
          >
            View Packages
          </button>
          <button
            className="px-8 py-3 bg-transparent text-parchment font-semibold border-2 border-parchment rounded-full hover:bg-latte hover:text-mocha hover:border-latte transition-all"
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
