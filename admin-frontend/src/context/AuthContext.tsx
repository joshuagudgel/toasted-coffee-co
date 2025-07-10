import React, { createContext, useContext, useState, useEffect } from "react";

interface AuthContextType {
  isAuthenticated: boolean;
  token: string | null;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  isLoading: boolean;
  user: { id: number; role: string } | undefined;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [user, setUser] = useState<{ id: number; role: string } | undefined>();
  const [token, setToken] = useState<string | null>(null);
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  useEffect(() => {
    // On mount, check if there's a token in localStorage
    const storedToken = localStorage.getItem("access_token");
    if (storedToken) {
      setToken(storedToken);
      validateToken(storedToken);
    } else {
      setIsLoading(false);
    }
  }, []);

  const validateToken = async (currentToken: string) => {
    try {
      const response = await fetch(`${API_URL}/api/v1/auth/validate`, {
        headers: {
          Authorization: `Bearer ${currentToken}`,
        },
      });

      if (response.ok) {
        const userData = await response.json();
        setUser({ id: userData.userId, role: userData.role });
        setIsAuthenticated(true);
      } else {
        // Clear invalid token
        localStorage.removeItem("access_token");
        setToken(null);
        setIsAuthenticated(false);
        setUser(undefined);
      }
    } catch (error) {
      console.error("Token validation error:", error);
      localStorage.removeItem("access_token");
      setToken(null);
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
      });

      if (!response.ok) {
        return false;
      }

      const { token: accessToken, user: userData } = await response.json();

      // Store token in state and localStorage
      setToken(accessToken);
      localStorage.setItem("access_token", accessToken);

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
      // Call logout endpoint (backend should invalidate token)
      await fetch(`${API_URL}/api/v1/auth/logout`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      // Clear local state and storage
      localStorage.removeItem("access_token");
      setToken(null);
      setUser(undefined);
      setIsAuthenticated(false);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        token,
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
