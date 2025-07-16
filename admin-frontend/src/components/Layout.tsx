import { Link, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useEffect, useState } from "react";

export default function Layout() {
  const { logout } = useAuth();
  const location = useLocation();
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  useEffect(() => {
    console.log("Route changed to:", location.pathname);
    // Close mobile menu when route changes
    setIsMenuOpen(false);
  }, [location.pathname]);

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-terracotta text-parchment p-4 shadow-md">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold">Toasted Coffee Admin</h1>

          {/* Mobile Menu Button */}
          <button
            className="md:hidden text-parchment focus:outline-none"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
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
                d={
                  isMenuOpen
                    ? "M6 18L18 6M6 6l12 12"
                    : "M4 6h16M4 12h16M4 18h16"
                }
              />
            </svg>
          </button>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex space-x-4">
            <Link
              to="/"
              className={`px-3 py-2 rounded transition ${
                location.pathname === "/" ||
                location.pathname.startsWith("/booking/")
                  ? "bg-mocha text-parchment"
                  : "hover:text-latte"
              }`}
            >
              Bookings
            </Link>
            <Link
              to="/menu"
              className={`px-3 py-2 rounded transition ${
                location.pathname === "/menu"
                  ? "bg-mocha text-parchment"
                  : "hover:text-latte"
              }`}
            >
              Menu
            </Link>
            <Link
              to="/packages"
              className={`px-3 py-2 rounded transition ${
                location.pathname === "/packages"
                  ? "bg-mocha text-parchment"
                  : "hover:text-latte"
              }`}
            >
              Packages
            </Link>
            <button
              onClick={handleLogout}
              className="ml-4 bg-mocha hover:bg-espresso text-parchment px-3 py-2 rounded transition"
            >
              Logout
            </button>
          </nav>
        </div>

        {/* Mobile Menu */}
        {isMenuOpen && (
          <div className="md:hidden mt-4 py-2 bg-mocha/95 rounded-md">
            <div className="flex flex-col space-y-2 px-2">
              <Link
                to="/"
                className={`text-center px-4 py-2 ${
                  location.pathname === "/" ||
                  location.pathname.startsWith("/booking/")
                    ? "bg-terracotta text-parchment"
                    : "bg-parchment text-mocha hover:bg-latte"
                } rounded-full font-medium transition`}
              >
                Bookings
              </Link>
              <Link
                to="/menu"
                className={`text-center px-4 py-2 ${
                  location.pathname === "/menu"
                    ? "bg-terracotta text-parchment"
                    : "bg-parchment text-mocha hover:bg-latte"
                } rounded-full font-medium transition`}
              >
                Menu
              </Link>
              <Link
                to="/packages"
                className={`text-center px-4 py-2 ${
                  location.pathname === "/packages"
                    ? "bg-terracotta text-parchment"
                    : "bg-parchment text-mocha hover:bg-latte"
                } rounded-full font-medium transition`}
              >
                Packages
              </Link>
              <button
                onClick={handleLogout}
                className="text-center px-4 py-2 bg-parchment text-mocha hover:bg-latte rounded-full font-medium transition"
              >
                Logout
              </button>
            </div>
          </div>
        )}
      </header>

      <main className="container mx-auto p-4 mt-8">
        <Outlet />
      </main>
    </div>
  );
}
