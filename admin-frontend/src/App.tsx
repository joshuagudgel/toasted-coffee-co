import { BrowserRouter } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import { MenuProvider } from "./context/MenuContext";
import { PackageProvider } from "./context/PackageContext";
import AppRoutes from "./AppRoutes";

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <MenuProvider>
          <PackageProvider>
            <AppRoutes />
          </PackageProvider>
        </MenuProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}
