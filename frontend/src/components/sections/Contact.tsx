import React from "react";
import BeanShape from "../ui/BeanShape";

const Contact: React.FC = () => {
  return (
    <footer className="py-16 relative overflow-hidden">
      {/* Background Elements z-index 1-9 */}
      <div className="absolute inset-0 z-[1] bg-espresso bg-opacity-60"></div>

      {/* Background Decorative Elements z-index 10-19 */}
      <BeanShape
        position="z-[10]"
        color="#DD9D79"
        style={{
          left: "15%",
          bottom: "20%",
          transform: "translate(-50%, 50%) scale(1.2) rotate(25deg)",
        }}
      />

      <BeanShape
        position="z-[10]"
        color="#BF7454"
        style={{
          right: "10%",
          top: "30%",
          transform: "translate(50%, -50%) scale(1.1) rotate(-15deg)",
        }}
      />

      {/* Main Content z-index 20+ */}
      <div className="container mx-auto px-4 relative z-[20] text-center">
        <div className="mx-auto max-w-4xl bg-white bg-opacity-90 rounded-xl p-8 shadow-lg">
          <div className="flex flex-col md:flex-row items-center justify-center gap-8 mb-8">
            <div>
              <h3 className="text-xl font-bold text-mocha mb-2">Contact Us</h3>
              <p className="text-espresso">toastedcoffeeco@gmail.com</p>
              <p className="text-espresso">(805)858-8171</p>
            </div>
          </div>

          <div className="border-t border-latte pt-6 mt-8">
            <p className="text-sm text-espresso">
              © {new Date().getFullYear()} Toasted Coffee Co. All rights
              reserved.
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Contact;
