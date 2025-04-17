import React, { createContext, useState, useContext } from "react";

interface BookingContextType {
  isModalOpen: boolean;
  selectedPackage: string | null;
  openModal: (packageName?: string) => void;
  closeModal: () => void;
}

const BookingContext = createContext<BookingContextType | undefined>(undefined);

export const BookingProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedPackage, setSelectedPackage] = useState<string | null>(null);

  const openModal = (packageName?: string) => {
    setSelectedPackage(packageName || null);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
  };

  return (
    <BookingContext.Provider
      value={{ isModalOpen, selectedPackage, openModal, closeModal }}
    >
      {children}
    </BookingContext.Provider>
  );
};

export const useBooking = () => {
  const context = useContext(BookingContext);
  if (context === undefined) {
    throw new Error("useBooking must be used within a BookingProvider");
  }
  return context;
};
