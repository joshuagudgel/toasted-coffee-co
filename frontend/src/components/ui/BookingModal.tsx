import React, { useState } from "react";

type BookingFormData = {
  name: string;
  date: string;
  time: string;
  people: string;
  coffeeType: string;
  milkOption: string;
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
    coffeeType: "",
    milkOption: "",
    location: "",
    notes: "",
    package: selectedPackage || "",
  });

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

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Booking submitted:", formData);
    // Here you would typically send this data to your backend
    alert("Thank you for your booking request! We'll be in touch shortly.");
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-terracotta text-parchment p-4 flex justify-between items-center">
          <h2 className="text-xl font-bold">Book Your Coffee Experience</h2>
          <button onClick={onClose} className="text-white hover:text-latte">
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

            <div>
              <label className="block text-espresso font-medium mb-1">
                Coffee Type
              </label>
              <select
                name="coffeeType"
                value={formData.coffeeType}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
                required
              >
                <option value="">Select Coffee Type</option>
                {coffeeOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-espresso font-medium mb-1">
                Milk Option
              </label>
              <select
                name="milkOption"
                value={formData.milkOption}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
                required
              >
                <option value="">Select Milk Option</option>
                {milkOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
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
