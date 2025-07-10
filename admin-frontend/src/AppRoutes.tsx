import { Routes, Route, Navigate, useLocation } from "react-router-dom";
import Layout from "./components/Layout";
import BookingList from "./components/BookingList";
import BookingDetail from "./components/BookingDetail";
import MenuManagement from "./components/MenuManagement";
import PackageManagement from "./components/PackageManagement";
import SignIn from "./components/SignIn";
import { useAuth } from "./context/AuthContext";
import { useRef, useEffect } from "react";

export default function AppRoutes() {
  const { isAuthenticated, isLoading } = useAuth();
  const location = useLocation();

  // Add this to track previous state and prevent unnecessary redirects
  const prevStateRef = useRef({ isAuthenticated, path: location.pathname });

  // Only log significant changes
  useEffect(() => {
    const prevState = prevStateRef.current;
    if (
      prevState.isAuthenticated !== isAuthenticated ||
      prevState.path !== location.pathname
    ) {
      console.log("AppRoutes state change:", {
        isAuthenticated,
        isLoading,
        path: location.pathname,
        prevAuth: prevState.isAuthenticated,
        prevPath: prevState.path,
      });
      prevStateRef.current = { isAuthenticated, path: location.pathname };
    }
  }, [isAuthenticated, isLoading, location.pathname]);

  // Two separate route trees based on authentication state
  if (!isAuthenticated) {
    return (
      <Routes>
        <Route path="/signin" element={<SignIn />} />
        <Route path="*" element={<Navigate to="/signin" replace />} />
      </Routes>
    );
  }

  // Only render these routes if authenticated
  return (
    <Routes>
      <Route path="/signin" element={<Navigate to="/" replace />} />
      <Route path="/" element={<Layout />}>
        <Route
          index
          element={
            <BookingList
              hiddenColumns={[
                "id",
                "notes",
                "createdAt",
                "contact",
                "time",
                "people",
                "package",
              ]}
            />
          }
        />
        <Route path="booking/:id" element={<BookingDetail />} />
        <Route path="menu" element={<MenuManagement />} />
        <Route path="packages" element={<PackageManagement />} />
      </Route>
      <Route
        path="*"
        element={
          location.pathname === "/" ? (
            <div>Not Found</div>
          ) : (
            <Navigate to="/" replace />
          )
        }
      />
    </Routes>
  );
}
