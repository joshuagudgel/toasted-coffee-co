import { BrowserRouter, Routes, Route } from "react-router-dom";
import Layout from "./components/Layout";
import BookingList from "./components/BookingList";
import BookingDetail from "./components/BookingDetail";
import MenuManagement from "./components/MenuManagement";
import SignIn from "./components/SignIn";
import ProtectedRoute from "./components/ProtectedRoute";
import { AuthProvider } from "./context/AuthContext";
import { MenuProvider } from "./context/MenuContext";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <MenuProvider>
          <Routes>
            {/* Public Route for SignIn page */}
            <Route path="/signin" element={<SignIn />} />

            {/* Protected Route for the rest of the app */}
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Layout />
                </ProtectedRoute>
              }
            >
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
            </Route>
          </Routes>
        </MenuProvider>
      </BrowserRouter>
    </AuthProvider>
  );
}
