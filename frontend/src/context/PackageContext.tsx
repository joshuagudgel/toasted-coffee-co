import React, { createContext, useState, useContext, useEffect } from "react";

export interface Package {
  id: number;
  name: string;
  price: string;
  description: string;
  points: string[];
  displayOrder: number;
}

interface PackageContextType {
  packages: Package[];
  loading: boolean;
  error: string | null;
}

const PackageContext = createContext<PackageContextType | undefined>(undefined);

export const PackageProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [packages, setPackages] = useState<Package[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  useEffect(() => {
    const fetchPackages = async () => {
      try {
        const response = await fetch(`${API_URL}/api/v1/packages`);

        if (!response.ok) {
          throw new Error(`Failed to fetch packages: ${response.status}`);
        }

        const data = await response.json();
        setPackages(data);
      } catch (err) {
        console.error("Error fetching packages:", err);
        setError(err instanceof Error ? err.message : "Unknown error");
      } finally {
        setLoading(false);
      }
    };

    fetchPackages();
  }, [API_URL]);

  return (
    <PackageContext.Provider value={{ packages, loading, error }}>
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
