import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function SignIn() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  // Get the intended destination, or default to the home page
  const from = location.state?.from?.pathname || "/";

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsSubmitting(true);

    console.log("Attempting login with API URL:", import.meta.env.VITE_API_URL);

    try {
      const success = await login(username, password);
      if (success) {
        console.log("Login successful, navigating to:", from);
        navigate(from, { replace: true });
      } else {
        console.error("Login failed - server returned unsuccessful response");
        setError("Invalid username or password");
      }
    } catch (err) {
      console.error("Login error:", err);

      // Check for specific error types
      if (err instanceof Error) {
        if (err.message.includes("CORS")) {
          setError(
            "Network error: CORS issue detected. Please check API configuration."
          );
        } else {
          setError(`Error: ${err.message}`);
        }
      } else {
        setError("An unexpected error occurred during sign in");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-parchment">
      <div className="max-w-md w-full bg-white rounded-lg shadow-lg overflow-hidden">
        <div className="bg-terracotta py-6 px-6 text-center">
          <h2 className="text-3xl font-bold text-parchment">
            Toasted Coffee Admin
          </h2>
          <p className="text-latte mt-2">Sign in to access the admin panel</p>
        </div>

        <form onSubmit={handleSubmit} className="py-8 px-6 space-y-6">
          {error && (
            <div className="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
              <p>{error}</p>
            </div>
          )}

          <div>
            <label
              htmlFor="user"
              className="block text-sm font-medium text-espresso mb-1"
            >
              User
            </label>
            <input
              id="user"
              type="user"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
              required
            />
          </div>

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-espresso mb-1"
            >
              Password
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-terracotta"
              required
            />
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="w-full py-3 bg-terracotta hover:bg-latte text-parchment hover:text-mocha font-bold rounded-lg transition disabled:opacity-50"
          >
            {isSubmitting ? "Signing in..." : "Sign In"}
          </button>

          <div className="text-center text-sm text-gray-500 mt-4">
            <p>For demo: admin / admin</p>
          </div>
        </form>
      </div>
    </div>
  );
}
