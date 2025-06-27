import Navbar from "./components/ui/Navbar";
import InquiryModal from "./components/ui/InquiryModal";
import Hero from "./components/sections/Hero";
import Packages from "./components/sections/Packages";
import Menu from "./components/sections/Menu";
import Contact from "./components/sections/Contact";
import BookingModal from "./components/ui/BookingModal";
import { BookingProvider, useBooking } from "./context/BookingContext";
import { MenuProvider } from "./context/MenuContext";
import { InquiryProvider } from "./context/InquiryContext";

const AppContent = () => {
  const { isModalOpen, selectedPackage, closeModal } = useBooking();

  return (
    <div className="min-h-screen">
      <Navbar />
      <InquiryModal />
      <Hero />
      <Packages />
      <BookingModal
        isOpen={isModalOpen}
        onClose={closeModal}
        selectedPackage={selectedPackage}
      />
      <Menu />
      <Contact />
    </div>
  );
};

function App() {
  return (
    <MenuProvider>
      <BookingProvider>
        <InquiryProvider>
          <AppContent />
        </InquiryProvider>
      </BookingProvider>
    </MenuProvider>
  );
}

export default App;
