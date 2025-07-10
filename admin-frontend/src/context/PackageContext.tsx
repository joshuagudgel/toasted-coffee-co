import React, { createContext, useState, useContext, useCallback } from "react";
import { useAuth } from "./AuthContext";

export interface Package {
  id: number;
  name: string;
  price: string;
  description: string;
  points: string[];
  displayOrder: number;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface PackageInput {
  name: string;
  price: string;
  description: string;
  points: string[];
  displayOrder: number;
  active: boolean;
}

interface PackageContextType {
  packages: Package[];
  loading: boolean;
  error: string | null;
  fetchPackages: (includeInactive?: boolean) => Promise<boolean>;
  addPackage: (pkg: PackageInput) => Promise<void>;
  updatePackage: (id: number, pkg: PackageInput) => Promise<void>;
  deletePackage: (id: number) => Promise<void>;
}

const PackageContext = createContext<PackageContextType | undefined>(undefined);

export const PackageProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [packages, setPackages] = useState<Package[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Use the apiRequest function from AuthContext
  const { apiRequest } = useAuth();

  const fetchPackages = useCallback(
    async (includeInactive = false): Promise<boolean> => {
      setLoading(true);
      setError(null);

      try {
        const queryParam = includeInactive ? "?include_inactive=true" : "";
        
        // Use apiRequest instead of fetch
        const data = await apiRequest<Package[]>(`/api/v1/packages${queryParam}`);
        
        setPackages(data);
        return true;
      } catch (err) {
        console.error("Error in fetchPackages:", err);
        setError(err instanceof Error ? err.message : "Unknown error");
        return false;
      } finally {
        setLoading(false);
      }
    },
    [apiRequest]
  );

  const addPackage = async (pkg: PackageInput) => {
    try {
      // Use apiRequest instead of fetch
      await apiRequest<Package>("/api/v1/packages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(pkg),
      });

      // Refresh packages
      await fetchPackages();
    } catch (err) {
      throw err;
    }
  };

  const updatePackage = async (id: number, pkg: PackageInput) => {
    try {
      // Use apiRequest instead of fetch
      await apiRequest<Package>(`/api/v1/packages/${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(pkg),
      });

      // Refresh packages
      await fetchPackages();
    } catch (err) {
      throw err;
    }
  };

  const deletePackage = async (id: number) => {
    try {
      // Use apiRequest instead of fetch
      await apiRequest(`/api/v1/packages/${id}`, {
        method: "DELETE",
      });

      // Refresh packages
      await fetchPackages();
    } catch (err) {
      throw err;
    }
  };

  return (
    <PackageContext.Provider
      value={{
        packages,
        loading,
        error,
        fetchPackages,
        addPackage,
        updatePackage,
        deletePackage,
      }}
    >
      {children}
    </PackageContext.Provider>
  );
};

export const usePackages = () => {
  const context = useContext(PackageContext);
  if (context === undefined) {
    throw new Error("usePackages must be used within a PackageProvider");
  }
  return context;
};
