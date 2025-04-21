import React, { useState } from "react";

type BookingFormData = {
  name: string;
  date: string;
  time: string;
  people: string;
  coffeeFlavors: string[];
  milkOptions: string[];
  location: string;
  notes: string;
  package?: string;
};

type BookingModalProps = {
  isOpen: boolean;
  onClose: () => void;
  selectedPackage?: string;
};

const BookingModal: React.FC<BookingModalProps> = ({
  isOpen,
  onClose,
  selectedPackage,
}) => {
  const [formData, setFormData] = useState<BookingFormData>({
    name: "",
    date: "",
    time: "",
    people: "",
    coffeeFlavors: [],
    milkOptions: [],
    location: "",
    notes: "",
    package: selectedPackage || "",
  });

  console.log("Initial coffeeFlavor state:", formData.coffeeFlavors);

  const coffeeOptions = [
    { value: "french_toast", label: "French Toast" },
    { value: "dirty_vanilla_chai", label: "Dirty Vanilla Chai" },
    { value: "mexican_mocha", label: "Mexican Mocha" },
    { value: "cinnamon_brown_sugar", label: "Cinnamon Brown Sugar" },
    { value: "horchata", label: "Horchata (made w/ rice milk)" },
  ];

  const milkOptions = [
    { value: "whole", label: "Whole Milk" },
    { value: "half_and_half", label: "Half & Half" },
    { value: "oat", label: "Oat Milk" },
    { value: "almond", label: "Almond Milk" },
    { value: "rice", label: "Rice Milk" },
  ];

  const handleCheckBoxChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    field: "coffeeFlavors" | "milkOptions"
  ) => {
    const { value, checked } = e.target;

    setFormData((prev) => {
      if (checked) {
        // add value to array
        console.log("Adding value:", value, "to field:", field);
        console.log("Previous state:", prev[field]);
        return {
          ...prev,
          [field]: [...prev[field], value],
        };
      } else {
        // revove value from array
        console.log("Removing value:", value, "from field:", field);
        console.log("Previous state:", prev[field]);
        return {
          ...prev,
          [field]: prev[field].filter((item) => item !== value),
        };
      }
    });
  };

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Check if at least one option is selected for each
    if (
      formData.coffeeFlavors.length === 0 ||
      formData.milkOptions.length === 0
    ) {
      alert("Please select at least one coffee flavor and milk option");
      return;
    }

    // Before submitting:
    const dataToSubmit = {
      ...formData,
      people: parseInt(formData.people, 10), // Convert string to number
    };

    try {
      const response = await fetch("http://localhost:8080/api/v1/bookings", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(dataToSubmit),
      });

      if (!response.ok) throw new Error("Booking submission failed");

      const data = await response.json();
      console.log("Booking submitted:", data);
      alert("Thank you for your booking request! We'll be in touch shortly.");
      onClose();
    } catch (error) {
      console.error("Error:", error);
      alert("There was a problem submitting your booking. Please try again.");
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-parchment rounded-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-terracotta text-parchment p-4 flex justify-between items-center">
          <h2 className="text-xl font-bold">Book Your Coffee Experience</h2>
          <button onClick={onClose} className="text-parchment hover:text-latte">
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
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div>
              <label className="block text-espresso font-medium mb-1">
                Name
              </label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
                required
              />
            </div>

            <div>
              <label className="block text-espresso font-medium mb-1">
                Date
              </label>
              <input
                type="date"
                name="date"
                value={formData.date}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
                required
              />
            </div>

            <div>
              <label className="block text-espresso font-medium mb-1">
                Time
              </label>
              <input
                type="time"
                name="time"
                value={formData.time}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
                required
              />
            </div>

            <div>
              <label className="block text-espresso font-medium mb-1">
                Number of People
              </label>
              <input
                type="number"
                name="people"
                min="1"
                value={formData.people}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
                required
              />
            </div>

            <div className="md:col-span-2">
              <label className="block text-espresso font-medium mb-2">
                Coffee Flavors (select all that apply)
              </label>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 bg-white p-3 rounded-md border border-gray-300">
                {coffeeOptions.map((option) => (
                  <div key={option.value} className="flex items-center mb-2">
                    <input
                      type="checkbox"
                      id={`coffee-${option.value}`}
                      value={option.value}
                      checked={formData.coffeeFlavors.includes(option.value)}
                      onChange={(e) => handleCheckBoxChange(e, "coffeeFlavors")}
                      className="h-4 w-4 text-terracotta border-gray-300 rounded focus:ring-terracotta mr-2"
                    />
                    <label
                      htmlFor={`coffee-${option.value}`}
                      className="text-espresso"
                    >
                      {option.label}
                    </label>
                  </div>
                ))}
              </div>
            </div>

            <div className="md:col-span-2">
              <label className="block text-espresso font-medium mb-2">
                Milk Options (select all that apply)
              </label>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 bg-white p-3 rounded-md border border-gray-300">
                {milkOptions.map((option) => (
                  <div key={option.value} className="flex items-center mb-2">
                    <input
                      type="checkbox"
                      id={`milk-${option.value}`}
                      value={option.value}
                      checked={formData.milkOptions.includes(option.value)}
                      onChange={(e) => handleCheckBoxChange(e, "milkOptions")}
                      className="h-4 w-4 text-terracotta border-gray-300 rounded focus:ring-terracotta mr-2"
                    />
                    <label
                      htmlFor={`milk-${option.value}`}
                      className="text-espresso"
                    >
                      {option.label}
                    </label>
                  </div>
                ))}
              </div>
            </div>

            <div className="md:col-span-2">
              <label className="block text-espresso font-medium mb-1">
                Event Location
              </label>
              <input
                type="text"
                name="location"
                value={formData.location}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
                required
              />
            </div>

            <div className="md:col-span-2">
              <label className="block text-espresso font-medium mb-1">
                Additional Notes
              </label>
              <textarea
                name="notes"
                value={formData.notes}
                onChange={handleChange}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
              ></textarea>
            </div>

            {selectedPackage && (
              <div className="md:col-span-2">
                <div className="bg-latte/50 p-3 rounded-md">
                  <p className="text-espresso">
                    Selected Package:{" "}
                    <span className="font-semibold">{selectedPackage}</span>
                  </p>
                </div>
              </div>
            )}
          </div>

          <div className="flex justify-end">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-espresso mr-2"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-terracotta text-parchment rounded-md hover:bg-latte hover:text-mocha transition"
            >
              Submit Booking
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default BookingModal;
