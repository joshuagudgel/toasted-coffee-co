import Navbar from "./components/ui/Navbar";
import Hero from "./components/sections/Hero";
import Packages from "./components/sections/Packages";

function App() {
  return (
    <div className="min-h-screen">
      <Navbar />
      <Hero />
      <Packages />
      {/* Other sections will be added here later */}
    </div>
  );
}

export default App;
