import React, { useState } from "react";
import { useInquiry } from "../../context/InquiryContext";

type InquiryFormData = {
  name: string;
  email: string;
  phone: string;
  message: string;
};

const InquiryModal: React.FC = () => {
  const { isInquiryModalOpen, closeInquiryModal } = useInquiry();
  const [isSending, setIsSending] = useState(false);
  const [formData, setFormData] = useState<InquiryFormData>({
    name: "",
    email: "",
    phone: "",
    message: "",
  });

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validate form
    if (formData.name.trim() === "") {
      alert("Please enter your name.");
      return;
    }

    if (formData.email.trim() === "" && formData.phone.trim() === "") {
      alert("Please provide either an email or phone number.");
      return;
    }

    if (formData.message.trim() === "") {
      alert("Please enter a message.");
      return;
    }

    setIsSending(true);

    const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

    try {
      const response = await fetch(`${API_URL}/api/v1/contact`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error(`Server error (${response.status}): ${errorText}`);
        throw new Error(`Inquiry submission failed: ${response.status}`);
      }

      // Reset form and close modal
      setFormData({
        name: "",
        email: "",
        phone: "",
        message: "",
      });

      alert("Thank you for your message! We'll get back to you soon.");
      closeInquiryModal();
    } catch (error) {
      console.error("Error:", error);
      alert("There was a problem sending your message. Please try again.");
    } finally {
      setIsSending(false);
    }
  };

  if (!isInquiryModalOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-[200] flex items-center justify-center p-4">
      <div className="bg-parchment rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-terracotta text-parchment p-4 flex justify-between items-center">
          <h2 className="text-xl font-bold">Contact Us</h2>
          <button
            onClick={closeInquiryModal}
            className="text-parchment hover:text-latte"
          >
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
          <div className="space-y-4">
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
                Message
              </label>
              <textarea
                name="message"
                value={formData.message}
                onChange={handleChange}
                rows={5}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
                required
              ></textarea>
            </div>
          </div>

          <div className="flex justify-end mt-6">
            <button
              type="button"
              onClick={closeInquiryModal}
              className="px-4 py-2 text-espresso mr-2"
              disabled={isSending}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-terracotta text-parchment rounded-md hover:bg-latte hover:text-mocha transition disabled:opacity-50"
              disabled={isSending}
            >
              {isSending ? "Sending..." : "Send Message"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default InquiryModal;
