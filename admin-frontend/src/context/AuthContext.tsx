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
      return true;
    } catch (error) {
      console.error("Login error:", error);
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
