import React, { createContext, useContext, useState } from "react";

interface InquiryContextType {
  isInquiryModalOpen: boolean;
  openInquiryModal: () => void;
  closeInquiryModal: () => void;
}

const InquiryContext = createContext<InquiryContextType | undefined>(undefined);

export const InquiryProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isInquiryModalOpen, setIsInquiryModalOpen] = useState(false);

  const openInquiryModal = () => {
    setIsInquiryModalOpen(true);
  };

  const closeInquiryModal = () => {
    setIsInquiryModalOpen(false);
  };

  return (
    <InquiryContext.Provider
      value={{
        isInquiryModalOpen,
        openInquiryModal,
        closeInquiryModal,
      }}
    >
      {children}
    </InquiryContext.Provider>
  );
};

export const useInquiry = (): InquiryContextType => {
  const context = useContext(InquiryContext);
  if (context === undefined) {
    throw new Error("useInquiry must be used within an InquiryProvider");
  }
  return context;
};
