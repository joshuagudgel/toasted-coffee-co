import React, { useState } from "react";

const Navbar: React.FC = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  return (
    <nav className="absolute top-0 left-0 z-20 w-full py-4 px-6">
      <div className="container mx-auto flex items-center justify-between">
        {/* Logo */}
        <div className="text-2xl font-bold text-white">Toasted Coffee Co</div>

        {/* Desktop Menu */}
        <div className="hidden md:flex items-center space-x-8">
          <a href="#" className="text-white hover:text-amber-200 transition">
            Home
          </a>
          <a href="#" className="text-white hover:text-amber-200 transition">
            About
          </a>
          <a href="#" className="text-white hover:text-amber-200 transition">
            Packages
          </a>
          <a href="#" className="text-white hover:text-amber-200 transition">
            Contact
          </a>
          <button className="ml-4 px-4 py-2 bg-white text-amber-800 rounded-full font-medium hover:bg-amber-100 transition">
            Book Now
          </button>
        </div>

        {/* Mobile Menu Button */}
        <button
          className="md:hidden text-white"
          onClick={() => setIsMenuOpen(!isMenuOpen)}
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d={
                isMenuOpen ? "M6 18L18 6M6 6l12 12" : "M4 6h16M4 12h16M4 18h16"
              }
            />
          </svg>
        </button>
      </div>

      {/* Mobile Menu */}
      {isMenuOpen && (
        <div className="md:hidden absolute top-16 left-0 w-full bg-amber-900/95 py-4">
          <div className="container mx-auto flex flex-col space-y-3 px-6">
            <a
              href="#"
              className="text-white hover:text-amber-200 transition py-2"
            >
              Home
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-200 transition py-2"
            >
              About
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-200 transition py-2"
            >
              Packages
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-200 transition py-2"
            >
              Gallery
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-200 transition py-2"
            >
              Contact
            </a>
            <button className="mt-2 px-4 py-2 bg-white text-amber-800 rounded-full font-medium hover:bg-amber-100 transition w-full">
              Book Now
            </button>
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;
