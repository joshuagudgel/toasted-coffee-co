import React, { createContext, useState, useContext, useCallback } from "react";
import { useAuth } from "./AuthContext";

export interface MenuItem {
  id: number;
  value: string;
  label: string;
  type: "coffee_flavor" | "milk_option";
  active: boolean;
}

interface MenuContextType {
  coffeeItems: MenuItem[];
  milkItems: MenuItem[];
  loading: boolean;
  error: string | null;
  fetchMenuItems: () => Promise<boolean>;
  addMenuItem: (item: Omit<MenuItem, "id">) => Promise<void>;
  updateMenuItem: (id: number, item: Partial<MenuItem>) => Promise<void>;
  deleteMenuItem: (id: number) => Promise<void>;
}

const MenuContext = createContext<MenuContextType | undefined>(undefined);

export const MenuProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [coffeeItems, setCoffeeItems] = useState<MenuItem[]>([]);
  const [milkItems, setMilkItems] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Use apiRequest from AuthContext
  const { apiRequest } = useAuth();

  const fetchMenuItems = useCallback(async (): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      // Use apiRequest instead of fetch
      const data = await apiRequest<MenuItem[]>("/api/v1/menu");

      console.log("Data received:", data.length, "items");

      // Separate items by type
      setCoffeeItems(data.filter((item) => item.type === "coffee_flavor"));
      setMilkItems(data.filter((item) => item.type === "milk_option"));
      return true;
    } catch (err) {
      console.error("Error in fetchMenuItems:", err);
      setError(err instanceof Error ? err.message : "Unknown error");
      return false;
    } finally {
      setLoading(false);
    }
  }, [apiRequest]);

  const addMenuItem = async (item: Omit<MenuItem, "id">) => {
    try {
      // Use apiRequest instead of fetch
      await apiRequest<MenuItem>("/api/v1/menu", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(item),
      });

      // Refresh menu items
      await fetchMenuItems();
    } catch (err) {
      throw err;
    }
  };

  const updateMenuItem = async (id: number, item: Partial<MenuItem>) => {
    try {
      // Get current item first to maintain type and other fields
      const currentItem = [...coffeeItems, ...milkItems].find(
        (i) => i.id === id
      );
      if (!currentItem) {
        throw new Error("Item not found");
      }

      const updatedItem = { ...currentItem, ...item };

      // Use apiRequest instead of fetch
      await apiRequest<MenuItem>(`/api/v1/menu/${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedItem),
      });

      // Refresh menu items
      await fetchMenuItems();
    } catch (err) {
      throw err;
    }
  };

  const deleteMenuItem = async (id: number) => {
    try {
      // Use apiRequest instead of fetch
      await apiRequest(`/api/v1/menu/${id}`, {
        method: "DELETE",
      });

      // Refresh menu items
      await fetchMenuItems();
    } catch (err) {
      throw err;
    }
  };

  return (
    <MenuContext.Provider
      value={{
        coffeeItems,
        milkItems,
        loading,
        error,
        fetchMenuItems,
        addMenuItem,
        updateMenuItem,
        deleteMenuItem,
      }}
    >
      {children}
    </MenuContext.Provider>
  );
};

export const useMenu = () => {
  const context = useContext(MenuContext);
  if (context === undefined) {
    throw new Error("useMenu must be used within a MenuProvider");
  }
  return context;
};
