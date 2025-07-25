import React, { useState } from "react";
import { useMenu } from "../../context/MenuContext";

type BookingFormData = {
  name: string;
  email: string;
  phone: string;
  date: string;
  time: string;
  people: string;
  coffeeFlavors: string[];
  milkOptions: string[];
  location: string;
  notes: string;
  package?: string;
  isOutdoor: boolean;
  hasShade: boolean;
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
  const { coffeeOptions, milkOptions, loading } = useMenu();

  // Convert menu items to usable format for the form
  const coffeeChoices = coffeeOptions
    .filter((item) => item.active) // Only show active items
    .map((item) => ({
      value: item.value,
      label: item.label,
    }));

  const milkChoices = milkOptions
    .filter((item) => item.active)
    .map((item) => ({
      value: item.value,
      label: item.label,
    }));

  const [formData, setFormData] = useState<BookingFormData>({
    name: "",
    email: "",
    phone: "",
    date: "",
    time: "",
    people: "",
    coffeeFlavors: [],
    milkOptions: [],
    location: "",
    notes: "",
    package: selectedPackage || "",
    isOutdoor: false,
    hasShade: false,
  });

  const handleCheckBoxChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    field: "coffeeFlavors" | "milkOptions"
  ) => {
    const { value, checked } = e.target;

    setFormData((prev) => {
      if (checked) {
        // add value to array
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

  const handleCheckboxToggle = (field: string) => {
    setFormData((prev) => ({
      ...prev,
      [field]: !prev[field as keyof BookingFormData],
      // If turning off outdoor, also turn off shade
      ...(field === "isOutdoor" && !prev.isOutdoor === false ? { hasShade: false } : {})
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // validate for data
    if (formData.email === "" && formData.phone === "") {
      alert("Please provide at least one contact method (email or phone).");
      return;
    }

    const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";
    console.log(API_URL);
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
      people: parseInt(formData.people, 10),
      package: selectedPackage || formData.package,
    };

    console.log("Submitting booking with package:", dataToSubmit.package);

    try {
      const response = await fetch(`${API_URL}/api/v1/bookings`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(dataToSubmit),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error(`Server error (${response.status}): ${errorText}`);
        throw new Error(`Booking submission failed: ${response.status}`);
      }
      const data = await response.json();
      console.log("Booking submitted:", data);
      alert("Thank you for your booking request! We'll be in touch shortly.");

      // Reset form data to initial state
      setFormData({
        name: "",
        email: "",
        phone: "",
        date: "",
        time: "",
        people: "",
        coffeeFlavors: [],
        milkOptions: [],
        location: "",
        notes: "",
        package: "",
        isOutdoor: false,
        hasShade: false,
      });

      onClose();
    } catch (error) {
      console.error("Error:", error);
      alert("There was a problem submitting your booking. Please try again.");
    }
  };

  if (!isOpen) return null;

  if (loading) {
    return (
      <div className={`fixed inset-0 z-50 ${isOpen ? "block" : "hidden"}`}>
        <div className="fixed inset-0 bg-black bg-opacity-50"></div>
        <div className="fixed inset-0 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-lg w-full max-h-screen overflow-hidden">
            <div className="p-8 text-center">
              <p>Loading menu options...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

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
                Email
              </label>
              <input
                type="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
              />
            </div>

            <div>
              <label className="block text-espresso font-medium mb-1">
                Phone
              </label>
              <input
                type="tel"
                name="phone"
                value={formData.phone}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
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
                {coffeeChoices.map((option) => (
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
                {milkChoices.map((option) => (
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

            {/* Add this within your form */}
            <div className="md:col-span-2">
              <div className="bg-white p-3 rounded-md border border-gray-300">
                <div className="flex items-center mb-3">
                  <input
                    type="checkbox"
                    id="is-outdoor"
                    checked={formData.isOutdoor}
                    onChange={() => handleCheckboxToggle("isOutdoor")}
                    className="h-4 w-4 text-terracotta border-gray-300 rounded focus:ring-terracotta mr-2"
                  />
                  <label htmlFor="is-outdoor" className="text-espresso font-medium">
                    This event will be outdoors
                  </label>
                </div>
                
                {formData.isOutdoor && (
                  <div className="ml-6 flex items-center">
                    <input
                      type="checkbox"
                      id="has-shade"
                      checked={formData.hasShade}
                      onChange={() => handleCheckboxToggle("hasShade")}
                      disabled={!formData.isOutdoor}
                      className="h-4 w-4 text-terracotta border-gray-300 rounded focus:ring-terracotta mr-2"
                    />
                    <label htmlFor="has-shade" className="text-espresso">
                      The outdoor area will have shade
                    </label>
                  </div>
                )}
              </div>
            </div>

            {selectedPackage && selectedPackage.length > 0 && (
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
