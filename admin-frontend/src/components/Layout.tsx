import { Link, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useEffect } from "react";

export default function Layout() {
  const { logout } = useAuth();
  const location = useLocation();

  useEffect(() => {
    console.log("Route changed to:", location.pathname);
  }, [location.pathname]);

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-terracotta text-parchment p-4 shadow-md">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold">Toasted Coffee Admin</h1>
          <nav className="flex space-x-4">
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
      </header>

      <main className="container mx-auto p-4 mt-8">
        <Outlet />
      </main>
    </div>
  );
}
