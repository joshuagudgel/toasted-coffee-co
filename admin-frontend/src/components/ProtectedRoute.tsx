import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { JSX } from 'react';


interface ProtectedRouteProps {
  children: JSX.Element;
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const location = useLocation();
  
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-parchment">
        <div className="text-mocha">Loading...</div>
      </div>
    );
  }
  
  if (!isAuthenticated) {
    // Redirect to login page, but save the current location
    return <Navigate to="/signin" state={{ from: location }} replace />;
  }
  
  return children;
}