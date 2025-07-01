import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  useRef,
} from "react";

interface AuthContextType {
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  isLoading: boolean;
  user?: { id: number; role: string };
  validateToken: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [user, setUser] = useState<{ id: number; role: string } | undefined>();
  const [isValidating, setIsValidating] = useState(false);
  const validationTimeRef = useRef<number | null>(null);
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  useEffect(() => {
    const initialize = async () => {
      const apiAvailable = await checkApiConnection();
      if (!apiAvailable) {
        console.error("Cannot reach API server");
        setIsLoading(false);
        return;
      }

      validateToken();
    };

    initialize();
  }, []);

  useEffect(() => {
    console.log("Auth state change:", { isAuthenticated, isLoading });
  }, [isAuthenticated, isLoading]);

  const validateToken = useCallback(async () => {
    // Skip if already validating or validated recently
    if (isValidating) {
      console.log("Validation already in progress, skipping");
      return;
    }

    // Throttle validation calls
    const now = Date.now();
    if (validationTimeRef.current && now - validationTimeRef.current < 2000) {
      console.log("Validated recently, skipping");
      return;
    }

    try {
      setIsValidating(true);
      console.log("Actually validating token at:", new Date().toISOString());
      const response = await fetch(`${API_URL}/api/v1/auth/validate`, {
        credentials: "include",
      });

      console.log("Validation response:", response.status);

      validationTimeRef.current = Date.now();

      if (response.ok) {
        // Handle successful validation
        const userData = await response.json();
        setUser({ id: userData.userId, role: userData.role });
        setIsAuthenticated(true);
      } else {
        // Handle failed validation
        setIsAuthenticated(false);
        setUser(undefined);
      }
    } catch (error) {
      // Error handling...
    } finally {
      setIsLoading(false);
      setIsValidating(false);
      validationTimeRef.current = Date.now();
    }
  }, [API_URL, isValidating]);

  const login = async (
    username: string,
    password: string
  ): Promise<boolean> => {
    try {
      console.log("Login attempt for:", username);
      const response = await fetch(`${API_URL}/api/v1/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include",
      });

      console.log("Login response status:", response.status);

      if (!response.ok) {
        console.error("Login failed with status:", response.status);
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

  const checkApiConnection = async () => {
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 5000);

      const response = await fetch(`${API_URL}/api/v1/health`, {
        signal: controller.signal,
      });

      clearTimeout(timeoutId);
      console.log("API connection check:", response.status);

      return response.ok;
    } catch (error) {
      console.error("API connection error:", error);
      return false;
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
        validateToken,
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
