import React, { createContext, useContext, useState, useEffect } from "react";

interface AuthContextType {
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  isLoading: boolean;
  user?: { id: number; role: string };
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [user, setUser] = useState<{ id: number; role: string } | undefined>();
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  useEffect(() => {
    // Check if user is authenticated by validating credentials with backend
    validateToken();
  }, []);

  const validateToken = async () => {
    try {
      // With cookies, we just need to call the endpoint
      // The cookie will be sent automatically
      const response = await fetch(`${API_URL}/api/v1/auth/validate`, {
        credentials: "include", // Important: include cookies in request
      });

      if (response.ok) {
        const userData = await response.json();
        setUser({ id: userData.userId, role: userData.role });
        setIsAuthenticated(true);
      } else {
        setIsAuthenticated(false);
        setUser(undefined);
      }
    } catch (error) {
      console.error("Token validation error:", error);
      setIsAuthenticated(false);
      setUser(undefined);
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (
    username: string,
    password: string
  ): Promise<boolean> => {
    try {
      const response = await fetch(`${API_URL}/api/v1/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include",
      });

      if (!response.ok) {
        return false;
      }

      const { user: userData } = await response.json();
      setUser({ id: userData.id, role: userData.role });
      setIsAuthenticated(true);

      // Use a setTimeout to ensure state updates before redirect
      setTimeout(() => {
        // Force a validation check after login
        checkAuth();
      }, 100);

      return true;
    } catch (error) {
      console.error("Login error:", error);
      return false;
    }
  };

  // Add this function to explicitly check auth status
  const checkAuth = async () => {
    try {
      const response = await fetch(`${API_URL}/api/v1/auth/validate`, {
        credentials: "include",
      });

      const isValid = response.ok;
      console.log("Auth validation result:", isValid);
      setIsAuthenticated(isValid);

      // If validation fails, clear any user data
      if (!isValid) {
        setUser(undefined);
      }

      return isValid;
    } catch (error) {
      console.error("Auth validation error:", error);
      setIsAuthenticated(false);
      setUser(undefined);
      return false;
    }
  };

  const logout = async () => {
    try {
      // Call logout endpoint to clear cookies on server
      await fetch(`${API_URL}/api/v1/auth/logout`, {
        method: "POST",
        credentials: "include",
      });
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      // Always clear local state regardless of server response
      setUser(undefined);
      setIsAuthenticated(false);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        login,
        logout,
        isLoading,
        user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
