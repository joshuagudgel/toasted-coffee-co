import Navbar from "./components/ui/Navbar";
import Hero from "./components/sections/Hero";
import Packages from "./components/sections/Packages";
import Menu from "./components/sections/Menu";
import BookingModal from "./components/ui/BookingModal";
import { BookingProvider, useBooking } from "./context/BookingContext";
import { MenuProvider } from "./context/MenuContext";

const AppContent = () => {
  const { isModalOpen, selectedPackage, closeModal } = useBooking();

  return (
    <div className="min-h-screen">
      <Navbar />
      <Hero />
      <Packages />
      <BookingModal
        isOpen={isModalOpen}
        onClose={closeModal}
        selectedPackage={selectedPackage}
      />
      <Menu />
      {/* Other sections will be added here later */}
    </div>
  );
};

function App() {
  return (
    <MenuProvider>
      <BookingProvider>
        <AppContent />
      </BookingProvider>
    </MenuProvider>
  );
}

export default App;
