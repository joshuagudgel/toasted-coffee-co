import React, { createContext, useState, useContext, useCallback } from "react";

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

  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  // Add refreshToken functionality:
  const fetchMenuItems = useCallback(async (): Promise<boolean> => {
    setLoading(true);
    setError(null);

    try {
      // CHANGE: Remove token check and just use credentials: "include"
      const response = await fetch(`${API_URL}/api/v1/menu`, {
        credentials: "include", // Use HttpOnly cookies
      });

      if (response.status === 401) {
        console.log("401 Unauthorized response");
        setError("Your session has expired. Please log in again.");
        return false;
      }

      if (!response.ok) {
        throw new Error(`Failed to fetch menu items: ${response.status}`);
      }

      const data: MenuItem[] = await response.json();
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
  }, [API_URL]);

  // CHANGE: Update all other methods to use credentials instead of token
  const addMenuItem = async (item: Omit<MenuItem, "id">) => {
    try {
      const response = await fetch(`${API_URL}/api/v1/menu`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // Use HttpOnly cookies
        body: JSON.stringify(item),
      });

      if (!response.ok) {
        throw new Error(`Failed to add menu item: ${response.status}`);
      }

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

      const response = await fetch(`${API_URL}/api/v1/menu/${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // Use HttpOnly cookies
        body: JSON.stringify(updatedItem),
      });

      if (!response.ok) {
        throw new Error(`Failed to update menu item: ${response.status}`);
      }

      // Refresh menu items
      await fetchMenuItems();
    } catch (err) {
      throw err;
    }
  };

  const deleteMenuItem = async (id: number) => {
    try {
      const response = await fetch(`${API_URL}/api/v1/menu/${id}`, {
        method: "DELETE",
        credentials: "include", // Use HttpOnly cookies
      });

      if (!response.ok) {
        throw new Error(`Failed to delete menu item: ${response.status}`);
      }

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
