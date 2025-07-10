import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
} from "react";

interface AuthContextType {
  isAuthenticated: boolean;
  token: string | null;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  isLoading: boolean;
  user: { id: number; role: string } | undefined;
  apiRequest: <T>(endpoint: string, options?: RequestInit) => Promise<T>;
  checkAuthStatus: () => Promise<boolean>;
}

interface JwtPayload {
  userId: number;
  role: string;
  exp: number;
  iat: number;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  // State management
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [user, setUser] = useState<{ id: number; role: string } | undefined>();
  const [token, setToken] = useState<string | null>(null);
  const [refreshToken, setRefreshToken] = useState<string | null>(null);
  const [lastActivity, setLastActivity] = useState<number>(Date.now());
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  // Constants for timeouts
  const INACTIVITY_TIMEOUT = 30 * 60 * 1000; // 30 minutes
  const TOKEN_REFRESH_BUFFER = 5 * 60 * 1000; // 5 minutes before expiration

  // Decode JWT token to get expiration time and payload
  const decodeToken = useCallback((jwtToken: string): JwtPayload | null => {
    try {
      const base64Url = jwtToken.split(".")[1];
      const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split("")
          .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
          .join("")
      );
      return JSON.parse(jsonPayload);
    } catch (error) {
      console.error("Error decoding token:", error);
      return null;
    }
  }, []);

  // Check if token is about to expire
  const isTokenExpiringSoon = useCallback(
    (tokenToCheck: string): boolean => {
      const decoded = decodeToken(tokenToCheck);
      if (!decoded) return true;

      const expiryTime = decoded.exp * 1000; // Convert to milliseconds
      return Date.now() > expiryTime - TOKEN_REFRESH_BUFFER;
    },
    [decodeToken]
  );

  // Update activity timestamp
  const updateActivity = useCallback(() => {
    setLastActivity(Date.now());
  }, []);

  /**
   * Core API request function that handles authentication tokens
   * and automatic token refresh.
   *
   * @param endpoint - API endpoint path (will be appended to API_URL)
   * @param options - Fetch options
   * @returns Promise with response data
   * @throws Error if authentication fails or request fails
   */
  const apiRequest = useCallback(
    async <T,>(endpoint: string, options: RequestInit = {}): Promise<T> => {
      // For login/refresh endpoints that don't need auth
      const isAuthEndpoint =
        endpoint.includes("/api/v1/auth/login") ||
        endpoint.includes("/api/v1/auth/refresh");

      if (!token && !isAuthEndpoint) {
        throw new Error("Not authenticated");
      }

      let currentToken = token;

      // If token is about to expire, refresh it first (but not for auth endpoints)
      if (
        currentToken &&
        !isAuthEndpoint &&
        isTokenExpiringSoon(currentToken)
      ) {
        console.log("Token about to expire, refreshing...");
        try {
          await refreshAccessToken();
          // Update current token after refresh
          currentToken = token;
          if (!currentToken) {
            throw new Error("Session expired");
          }
        } catch (error) {
          console.error("Failed to refresh token:", error);
          throw new Error("Session expired");
        }
      }

      // Prepare headers - only add Authorization for non-auth endpoints
      let headers: Record<string, string> = {
        "Content-Type": "application/json",
      };

      if (currentToken && !isAuthEndpoint) {
        headers["Authorization"] = `Bearer ${currentToken}`;
      }

      try {
        // Make the API request
        const response = await fetch(`${API_URL}${endpoint}`, {
          ...options,
          headers: {
            ...headers,
            ...(options.headers as Record<string, string>), // Cast to merge headers properly
          },
        });

        // Handle 401 by trying to refresh once (but not for auth endpoints)
        if (response.status === 401 && !isAuthEndpoint) {
          // Try to refresh the token
          try {
            const refreshed = await refreshAccessToken();
            if (refreshed && token) {
              // Retry the request with new token
              const retryResponse = await fetch(`${API_URL}${endpoint}`, {
                ...options,
                headers: {
                  ...options.headers,
                  Authorization: `Bearer ${token}`,
                },
              });

              if (retryResponse.ok) {
                updateActivity();
                return await retryResponse.json();
              } else {
                await logout();
                throw new Error("Not authenticated");
              }
            } else {
              await logout();
              throw new Error("Session expired");
            }
          } catch (error) {
            await logout();
            throw new Error("Session expired");
          }
        }

        if (!response.ok) {
          throw new Error(`API request failed: ${response.status}`);
        }

        if (!isAuthEndpoint) {
          updateActivity();
        }

        if (response.status === 204) {
          // No content response
          return {} as T;
        }

        return await response.json();
      } catch (error) {
        console.error("API request error:", error);
        throw error;
      }
    },
    [token, isTokenExpiringSoon, updateActivity, API_URL]
  );

  // Refresh access token using refresh token
  const refreshAccessToken = useCallback(
    async (currentRefreshToken?: string): Promise<boolean> => {
      const refreshTokenToUse = currentRefreshToken || refreshToken;
      if (!refreshTokenToUse) return false;

      try {
        // Use direct fetch instead of apiRequest to avoid circular dependency
        const response = await fetch(`${API_URL}/api/v1/auth/refresh`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ refreshToken: refreshTokenToUse }),
        });

        if (!response.ok) {
          console.error(`Refresh failed with status ${response.status}`);
          setIsAuthenticated(false);
          setToken(null);
          setUser(undefined);
          return false;
        }

        const data = await response.json();
        setToken(data.accessToken);

        // Update user data if needed
        if (!user) {
          // Also use direct fetch for validation if needed
          const userResponse = await fetch(`${API_URL}/api/v1/auth/validate`, {
            headers: {
              Authorization: `Bearer ${data.accessToken}`,
            },
          });

          if (userResponse.ok) {
            const userData = await userResponse.json();
            setUser({ id: userData.userId, role: userData.role });
          }
        }

        setIsAuthenticated(true);
        updateActivity();
        return true;
      } catch (error) {
        console.error("Error refreshing access token:", error);
        setIsAuthenticated(false);
        setToken(null);
        setUser(undefined);
        return false;
      }
    },
    [refreshToken, user, updateActivity, API_URL]
  );

  // Validate token and get user info
  // Update validateToken function
  const validateToken = useCallback(
    async (currentToken: string): Promise<boolean> => {
      try {
        // Use direct fetch instead of apiRequest to avoid circular dependency
        const response = await fetch(`${API_URL}/api/v1/auth/validate`, {
          headers: {
            Authorization: `Bearer ${currentToken}`,
          },
        });

        if (!response.ok) {
          setIsAuthenticated(false);
          setToken(null);
          setUser(undefined);
          return false;
        }

        const userData = await response.json();
        setUser({ id: userData.userId, role: userData.role });
        setIsAuthenticated(true);
        return true;
      } catch (error) {
        console.error("Token validation error:", error);
        setIsAuthenticated(false);
        setToken(null);
        setUser(undefined);
        return false;
      }
    },
    [API_URL] // Only depend on API_URL
  );

  // Login function
  const login = useCallback(
    async (username: string, password: string): Promise<boolean> => {
      try {
        const data = await apiRequest<{
          token: string;
          refreshToken: string;
          user: { id: number; role: string };
        }>("/api/v1/auth/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username, password }),
        });

        // Store refresh token in localStorage (more persistent but less sensitive)
        localStorage.setItem("refresh_token", data.refreshToken);

        // Store access token only in memory
        setToken(data.token);
        setRefreshToken(data.refreshToken);
        setUser({ id: data.user.id, role: data.user.role });
        setIsAuthenticated(true);
        updateActivity();
        return true;
      } catch (error) {
        console.error("Login error:", error);
        return false;
      }
    },
    [apiRequest, updateActivity]
  );

  // Logout function
  const logout = useCallback(async () => {
    try {
      if (token) {
        // Call logout endpoint
        await apiRequest("/api/v1/auth/logout", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ refreshToken }),
        }).catch((error) => {
          console.error("Logout API error (non-critical):", error);
        });
      }
    } finally {
      // Always clear local state and storage
      localStorage.removeItem("refresh_token");
      setToken(null);
      setRefreshToken(null);
      setUser(undefined);
      setIsAuthenticated(false);
    }
  }, [token, refreshToken, apiRequest]);

  // Check for inactivity
  useEffect(() => {
    if (!isAuthenticated) return;

    const activityEvents = ["mousedown", "keypress", "scroll", "touchstart"];

    // Add event listeners
    activityEvents.forEach((event) => {
      window.addEventListener(event, updateActivity);
    });

    // Check for inactivity periodically
    const inactivityCheck = setInterval(() => {
      if (Date.now() - lastActivity > INACTIVITY_TIMEOUT) {
        console.log("User inactive, logging out");
        logout();
      }
    }, 60000); // Check every minute

    return () => {
      // Clean up event listeners
      activityEvents.forEach((event) => {
        window.removeEventListener(event, updateActivity);
      });
      clearInterval(inactivityCheck);
    };
  }, [
    isAuthenticated,
    lastActivity,
    updateActivity,
    logout,
    INACTIVITY_TIMEOUT,
  ]);

  // Auto-refresh token before expiration
  useEffect(() => {
    if (!token || !refreshToken) return;

    const tokenRefresh = async () => {
      try {
        if (isTokenExpiringSoon(token)) {
          console.log("Token about to expire, refreshing...");
          await refreshAccessToken();
        }
      } catch (error) {
        console.error("Error in token refresh check:", error);
      }
    };

    // Initial check
    tokenRefresh();

    // Set up interval for periodic checks
    const refreshInterval = setInterval(tokenRefresh, 60000); // Check every minute

    return () => clearInterval(refreshInterval);
  }, [token, refreshToken, isTokenExpiringSoon, refreshAccessToken]);

  // On mount, try to restore session from refresh token
  useEffect(() => {
    const restoreSession = async () => {
      setIsLoading(true);

      // Check if we have a refresh token stored
      const storedRefreshToken = localStorage.getItem("refresh_token");
      if (storedRefreshToken) {
        setRefreshToken(storedRefreshToken);
        try {
          const success = await refreshAccessToken(storedRefreshToken);
          if (!success) {
            // Clear storage if refresh failed
            localStorage.removeItem("refresh_token");
            setRefreshToken(null);
          }
        } catch (error) {
          console.error("Session restoration error:", error);
          localStorage.removeItem("refresh_token");
          setRefreshToken(null);
        }
      }

      setIsLoading(false);
    };

    restoreSession();
  }, [refreshAccessToken]);

  // Public method to check auth status (for protected routes)
  const checkAuthStatus = useCallback(async (): Promise<boolean> => {
    if (token) {
      // If we have a token, validate it
      return validateToken(token);
    } else if (refreshToken) {
      // If we have a refresh token but no access token, try to refresh
      return refreshAccessToken();
    }
    return false;
  }, [token, refreshToken, validateToken, refreshAccessToken]);

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        token,
        login,
        logout,
        isLoading,
        user,
        apiRequest,
        checkAuthStatus,
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
