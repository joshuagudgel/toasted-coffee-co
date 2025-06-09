import React from "react";
import { useMenu } from "../../context/MenuContext";

const Menu: React.FC = () => {
  const { coffeeOptions, milkOptions, loading } = useMenu();
  
  // Only show active items
  const activeCoffeeOptions = coffeeOptions.filter(option => option.active);
  const activeMilkOptions = milkOptions.filter(option => option.active);

  return (
    <section id="menu" className="py-20 bg-caramel relative overflow-hidden">
      {/* title area above */}
      <div className="relative z-10 container mx-auto px-4 text-center">
        {/* Menu content on parchment */}
        <div className="mt-12 mx-auto max-w-4xl px-4 bg-parchment bg-opacity-90 rounded-xl p-8 shadow-lg">
          <h2 className="text-3xl md:text-4xl font-bold text-mocha mb-4 tracking-tight">
            TOASTED COFFEE CO
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <h2 className="text-2xl font-bold text-terracotta mb-4 tracking-tight">
                COLD BREW SPECIALTIES
              </h2>
              {loading ? (
                <p>Loading...</p>
              ) : (
                <ul className="list-none list-inside text-lg text-espresso">
                  {activeCoffeeOptions.map((option) => (
                    <li key={option.id}>{option.label}</li>
                  ))}
                </ul>
              )}
            </div>
            <div>
              <h2 className="text-2xl font-bold text-terracotta mb-4 tracking-tight">
                MILK OPTIONS
              </h2>
              {loading ? (
                <p>Loading...</p>
              ) : (
                <ul className="list-none list-inside text-lg text-espresso">
                  {activeMilkOptions.map((option) => (
                    <li key={option.id}>{option.label}</li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Menu;
