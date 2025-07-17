import React from "react";
import { useMenu } from "../../context/MenuContext";
import toastedLogo from "../../assets/toasted-logo.png";

const Menu: React.FC = () => {
  const { coffeeOptions, milkOptions, loading } = useMenu();

  // Only show active items
  const activeCoffeeOptions = coffeeOptions.filter((option) => option.active);
  const activeMilkOptions = milkOptions.filter((option) => option.active);

  return (
    <section id="menu" className="py-20 relative overflow-hidden">
      {/* Background Elements z-index 1-9 */}
      <div className="absolute inset-0 z-[1] bg-caramel"></div>

      {/* Background Decorative Elements z-index 10-19 */}

      {/* Main Content z-index 20+ */}
      <div className="container mx-auto px-4 relative z-[20] text-center">
        {/* Menu content card */}
        <div className="mt-12 mx-auto max-w-4xl px-4 bg-parchment rounded-xl p-8 shadow-lg relative z-[30]">
          <div className="flex justify-center mb-4">
            <img
              src={toastedLogo}
              alt="TOASTED COFFEE CO"
              className="h-40 md:h-60"
            />
          </div>
          <p className="text-base text-espresso max-w-xl mx-auto mb-10">
            We source only the highest quality beans for our signature cold
            brew. Choose from our specialty flavors and milk options for a
            perfectly crafted experience.
          </p>
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
                    <li key={option.id} className="mb-2 pb-2">
                      {option.label}
                    </li>
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
                    <li key={option.id} className="mb-2 pb-2">
                      {option.label}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>

          {/* Additional menu information */}
          <div className="mt-8 pt-6 text-center">
            <p className="text-lg text-espresso italic">
              All drinks served over ice. Ask about our seasonal specials!
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Menu;
