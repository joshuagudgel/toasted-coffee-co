import React, { useState } from "react";
import { useBooking } from "../../context/BookingContext";
import { useInquiry } from "../../context/InquiryContext";
import { scrollToSection } from "../../utils/scrollUtils";

const Navbar: React.FC = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const { openModal } = useBooking();
  const { openInquiryModal } = useInquiry();

  return (
    <nav className="absolute top-0 left-0 z-[100] w-full py-4 px-6">
      <div className="container mx-auto flex items-center justify-end">
        {/* Desktop Menu */}
        <div className="hidden md:flex items-center space-x-8">
          <button
            className="ml-2 px-4 py-2 bg-parchment text-mocha rounded-full font-medium hover:bg-latte transition"
            onClick={() => scrollToSection("packages")}
          >
            Packages
          </button>
          <button
            className="ml-2 px-4 py-2 bg-parchment text-mocha rounded-full font-medium hover:bg-latte transition"
            onClick={() => scrollToSection("menu")}
          >
            Menu
          </button>
          <button
            className="ml-4 px-4 py-2 bg-parchment text-mocha rounded-full font-medium hover:bg-latte transition"
            onClick={() => {
              openInquiryModal();
            }}
          >
            Inquire
          </button>
          <button
            className="ml-4 px-4 py-2 bg-parchment text-mocha rounded-full font-medium hover:bg-latte transition"
            onClick={() => {
              setIsMenuOpen(false);
              openModal();
            }}
          >
            Book Now
          </button>
        </div>

        {/* Mobile Menu Button */}
        <button
          className="md:hidden text-mocha focus:outline-none"
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
        <div className="md:hidden absolute top-16 left-0 w-full bg-mocha/95 py-4">
          <div className="container mx-auto flex flex-col space-y-3 px-6">
            <a
              href="#"
              className="text-center ml-4 px-4 py-2 bg-parchment text-mocha rounded-full font-medium hover:bg-latte transition"
              onClick={() => scrollToSection("packages")}
            >
              Packages
            </a>
            <a
              href="#"
              className="text-center ml-4 px-4 py-2 bg-parchment text-mocha rounded-full font-medium hover:bg-latte transition"
              onClick={() => scrollToSection("menu")}
            >
              Menu
            </a>
            <button
              className="text-center ml-4 px-4 py-2 bg-parchment text-mocha rounded-full font-medium hover:bg-latte transition"
              onClick={() => {
                openInquiryModal();
              }}
            >
              Inquire
            </button>
            <button
              className="text-center ml-4 px-4 py-2 bg-parchment text-mocha rounded-full font-medium hover:bg-latte transition"
              onClick={() => {
                setIsMenuOpen(false);
                openModal();
              }}
            >
              Book Now
            </button>
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;
