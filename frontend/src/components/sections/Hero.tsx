import React from "react";

const Hero: React.FC = () => {
  return (
    <div className="relative h-screen w-full bg-gradient-to-r from-espresso to-coffee-800 overflow-hidden">
      {/* Dark overlay for better text contrast */}
      <div className="absolute inset-0 bg-black/40"></div>

      {/* Background pattern - optional */}
      <div className="absolute inset-0 opacity-10 bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI2MCIgaGVpZ2h0PSI2MCIgdmlld0JveD0iMCAwIDYwIDYwIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNmZmZmZmYiIGZpbGwtb3BhY2l0eT0iMSI+PHBhdGggZD0iTTM2IDM0djZoLTZWMzRoLTZ2LTZoNnYtNmg2djZoNnY2aC02eiIvPjwvZz48L2c+PC9zdmc+')]"></div>

      {/* Content container */}
      <div className="relative z-10 h-full flex flex-col items-center justify-center text-center px-4 md:px-8">
        <h1 className="text-5xl md:text-7xl font-bold text-toasted mb-4 tracking-tight">
          Toasted Coffee Co
        </h1>
        <p className="text-cream text-xl md:text-2xl max-w-2xl mb-8">
          Cold Brew Bar
        </p>
        <div className="flex flex-col sm:flex-row gap-4">
          <button className="px-8 py-3 bg-white text-amber-800 font-semibold rounded-full hover:bg-amber-100 transition-all shadow-lg">
            View Packages
          </button>
          <button className="px-8 py-3 bg-transparent text-white border-2 border-white rounded-full hover:bg-white/10 transition-all">
            Book Now
          </button>
        </div>
      </div>
    </div>
  );
};

export default Hero;
