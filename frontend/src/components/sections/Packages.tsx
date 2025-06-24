import React from "react";
import { useBooking } from "../../context/BookingContext";
import BeanShape from "../ui/BeanShape";
//import packagesBeanTop from "../../assets/packages-bean-top.png";
//import packagesBeanLeft from "../../assets/packages-bean-left.png";
//import packagesBeanBottom from "../../assets/packages-bean-bottom.png";

type Package = {
  id: number;
  name: string;
  price: string;
  description: string;
  features: string[];
};

const coffeePackages: Package[] = [
  {
    id: 1,
    name: "Group",
    price: "$135",
    description: "Small gatherings of up to 25 people",
    features: ["25 People", "1 hour service"],
  },
  {
    id: 2,
    name: "Crowd",
    price: "$275",
    description: "Ideal for gatherings of up to 50 people",
    features: ["50 People", "1.5 hour service"],
  },
  {
    id: 3,
    name: "Party",
    price: "$410",
    description: "Guarenteeing coffee for up to 75 people",
    features: ["75 People", "2 hour service"],
  },
];

const Packages: React.FC = () => {
  const { openModal } = useBooking();

  return (
    <section id="packages" className="relative py-20 overflow-hidden">
      {/* Background Elements z-index 1-9 */}
      <div className="absolute inset-0 z-[1] bg-parchment"></div>

      {/* Background Decorative Elements z-index 10-19 */}
      {/* Top right bean */}
      <BeanShape
        position="z-[10]"
        color="#DD9D79"
        style={{
          left: "75%",
          top: "25%",
          transform: "translate(-50%, -50%) scale(1.5) rotate(15deg)",
        }}
      />

      {/* Bottom left bean */}
      <BeanShape
        position="z-[10]"
        color="#FAE7D3"
        style={{
          left: "25%",
          top: "66.67%",
          transform: "translate(-50%, -50%) scale(1.2) rotate(-20deg)",
        }}
      />

      {/* Bottom right bean */}
      <BeanShape
        position="z-[10]"
        color={"#DD9D79"}
        style={{
          left: "80%",
          top: "80%",
          transform: "translate(-50%, -50%) scale(1.8) rotate(45deg)",
        }}
      />

      {/* Main Content z-index 20+ */}
      <div className="container mx-auto px-4 relative z-[20]">
        <div className="text-center mb-16">
          <div className="flex items-center justify-center mb-8">
            <div className="hidden md:block w-24 h-1 bg-terracotta rounded-full"></div>
            <h2 className="text-4xl font-bold text-espresso mx-4 tracking-tight">
              Hand Crafted Service
            </h2>
            <div className="hidden md:block w-24 h-1 bg-terracotta rounded-full"></div>
          </div>
          <p className="text-lg text-espresso max-w-xl mx-auto">
            We offer 3 basic packages but we understand that not every event is
            the same. We can tailor to your event, email us with specifics and
            we'll provide you a quote!
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {coffeePackages.map((pkg) => (
            <div
              key={pkg.id}
              className="bg-parchment rounded-xl shadow-lg overflow-hidden transition-all hover:shadow-xl hover:-translate-y-1"
            >
              <div className="bg-terracotta py-4">
                <h3 className="text-2xl font-bold text-center text-parchment">
                  {pkg.name}
                </h3>
              </div>
              <div className="p-8">
                <p className="text-3xl font-bold text-center mb-4 text-espresso">
                  {pkg.price}
                </p>
                <p className="text-center mb-6 text-espresso">
                  {pkg.description}
                </p>
                <ul className="space-y-3 mb-8">
                  {pkg.features.map((feature, index) => (
                    <li key={index} className="flex items-center">
                      <svg
                        className="w-5 h-5 text-peach mr-2"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      {feature}
                    </li>
                  ))}
                </ul>
                <button
                  className="w-full py-3 bg-terracotta hover:bg-latte text-parchment hover:text-mocha font-bold rounded-lg transition"
                  onClick={() => openModal(pkg.name)}
                >
                  Select Package
                </button>
              </div>
            </div>
          ))}
        </div>
        <div className="text-center my-16">
          <p className="text-lg text-espresso max-w-xl mx-auto">
            Every service includes cold brew based drinks, paper goods, and
            setup/tear down (please allow 30 min.)
          </p>
          <p className="text-lg text-espresso max-w-xl mx-auto">
            *Bookings must be made 2 weeks in advance
          </p>
        </div>
      </div>
    </section>
  );
};

export default Packages;
