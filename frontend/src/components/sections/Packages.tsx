import React from "react";

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
    description: "Small gatherings up to 25 people",
    features: [
      "25 People",
      "1 hour service",
      "Cold brew based drinks",
      "Paper goods",
      "Setup/tear down (please allow 30 min.)",
    ],
  },
  {
    id: 2,
    name: "Crowd",
    price: "$275",
    description: "Ideal for gatherings of up to 50 people",
    features: [
      "50 People",
      "1.5 hour service",
      "Cold brew based drinks",
      "Paper goods",
      "Setup/tear down (please allow 30 min.)",
    ],
  },
  {
    id: 3,
    name: "Party",
    price: "$410",
    description: "Guarenteeing coffee for up to 75 people",
    features: [
      "75 People",
      "2 hour service",
      "Cold brew based drinks",
      "Paper goods",
      "Setup/tear down (please allow 30 min.)",
    ],
  },
];

const Packages: React.FC = () => {
  return (
    <section className="py-20 bg-amber-50">
      <div className="container mx-auto px-4">
        <div className="text-center mb-16">
          <h2 className="text-4xl font-bold text-amber-900 mb-4">
            Our Packages
          </h2>
          <p className="text-lg text-amber-700 max-w-2xl mx-auto">
            At Toasted Coffee Co, our goal is to elevate your event with an
            exceptional cold brew cart.
          </p>
          <p className="text-lg text-amber-700 max-w-2xl mx-auto">
            Whether you're hosting a corporate gathering or celebrating a
            special occasion, we strive to deliver the perfect cup every time.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {coffeePackages.map((pkg) => (
            <div
              key={pkg.id}
              className="bg-white rounded-xl shadow-lg overflow-hidden transition-all hover:shadow-xl hover:-translate-y-1"
            >
              <div className="bg-amber-700 py-4">
                <h3 className="text-2xl font-bold text-center text-white">
                  {pkg.name}
                </h3>
              </div>
              <div className="p-8">
                <p className="text-3xl font-bold text-center mb-4 text-amber-900">
                  {pkg.price}
                </p>
                <p className="text-center mb-6 text-gray-600">
                  {pkg.description}
                </p>
                <ul className="space-y-3 mb-8">
                  {pkg.features.map((feature, index) => (
                    <li key={index} className="flex items-center">
                      <svg
                        className="w-5 h-5 text-amber-500 mr-2"
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
                <button className="w-full py-3 bg-amber-700 hover:bg-amber-800 text-white font-bold rounded-lg transition">
                  Select Package
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Packages;
