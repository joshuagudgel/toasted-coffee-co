import { Link, Outlet } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function Layout() {
  const { logout } = useAuth();

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-terracotta text-parchment p-4 shadow-md">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold">Toasted Coffee Admin</h1>
          <nav>
            <Link to="/" className="hover:text-latte px-3 py-2">
              Bookings
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
