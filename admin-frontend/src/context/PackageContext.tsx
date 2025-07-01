import React, { createContext, useState, useContext, useCallback } from "react";

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

  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  const fetchPackages = useCallback(
    async (includeInactive = false): Promise<boolean> => {
      setLoading(true);
      setError(null);

      try {
        const queryParam = includeInactive ? "?include_inactive=true" : "";
        const response = await fetch(
          `${API_URL}/api/v1/packages${queryParam}`,
          {
            credentials: "include",
          }
        );

        if (response.status === 401) {
          setError("Not authenticated");
          return false;
        }

        if (!response.ok) {
          throw new Error(`Failed to fetch packages: ${response.status}`);
        }

        const data = await response.json();
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
    [API_URL]
  );

  // Add Content-Type header to your API requests
  const addPackage = async (pkg: PackageInput) => {
    try {
      const response = await fetch(`${API_URL}/api/v1/packages`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(pkg),
      });

      if (!response.ok) {
        throw new Error(`Failed to add package: ${response.status}`);
      }

      // Refresh packages
      await fetchPackages();
    } catch (err) {
      throw err;
    }
  };

  const updatePackage = async (id: number, pkg: PackageInput) => {
    try {
      const response = await fetch(`${API_URL}/api/v1/packages/${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(pkg),
      });

      if (!response.ok) {
        throw new Error(`Failed to update package: ${response.status}`);
      }

      // Refresh packages
      await fetchPackages();
    } catch (err) {
      throw err;
    }
  };

  const deletePackage = async (id: number) => {
    try {
      const response = await fetch(`${API_URL}/api/v1/packages/${id}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error(`Failed to delete package: ${response.status}`);
      }

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
