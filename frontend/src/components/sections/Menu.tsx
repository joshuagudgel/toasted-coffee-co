import React from "react";

const Menu: React.FC = () => {
  return (
    <section id="menu" className="py-20 bg-caramel relative overflow-hidden">
      {/* Content Container - now elevated above the grid */}
      <div className="relative z-10 container mx-auto px-4 text-center">
        <h1 className="text-5xl md:text-7xl font-bold text-mocha mb-4 tracking-tight">
          MENU
        </h1>

        {/* Add your menu content here */}
        <div className="mt-12 mx-auto max-w-4xl px-4 bg-parchment bg-opacity-90 rounded-xl p-8 shadow-lg">
          <p className="text-xl text-mocha mb-8">
            Our signature cold brew drinks
          </p>
          {/* Menu items will go here */}
        </div>
      </div>
    </section>
  );
};

export default Menu;
