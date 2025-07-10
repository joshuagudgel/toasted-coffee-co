import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { JSX, memo, useEffect } from "react";

interface ProtectedRouteProps {
  children: JSX.Element;
}

function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, checkAuthStatus } = useAuth();
  const location = useLocation();

  // Verify authentication on protected route access
  useEffect(() => {
    const verifyAuth = async () => {
      if (!isAuthenticated && !isLoading) {
        await checkAuthStatus();
      }
    };

    verifyAuth();
  }, [isAuthenticated, isLoading, checkAuthStatus]);

  console.log("ProtectedRoute render:", {
    isAuthenticated,
    isLoading,
    path: location.pathname,
  });

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-parchment">
        <div className="text-mocha">Loading...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    // Redirect to login page, but save the current location
    console.log("Not authenticated, redirecting to login");
    return <Navigate to="/signin" state={{ from: location }} replace />;
  }

  console.log("Authenticated, rendering protected content");
  return children;
}

export default memo(ProtectedRoute);
