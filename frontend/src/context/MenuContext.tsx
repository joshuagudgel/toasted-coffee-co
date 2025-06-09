import React, { createContext, useState, useContext, useEffect } from "react";

interface MenuItem {
  id: number;
  value: string;
  label: string;
  type: "coffee_flavor" | "milk_option";
  active: boolean;
}

interface MenuContextType {
  coffeeOptions: MenuItem[];
  milkOptions: MenuItem[];
  loading: boolean;
  error: string | null;
}

const MenuContext = createContext<MenuContextType>({
  coffeeOptions: [],
  milkOptions: [],
  loading: true,
  error: null,
});

export const useMenu = () => useContext(MenuContext);

export const MenuProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [coffeeOptions, setCoffeeOptions] = useState<MenuItem[]>([]);
  const [milkOptions, setMilkOptions] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  useEffect(() => {
    async function fetchMenuItems() {
      try {
        const response = await fetch(`${API_URL}/api/v1/menu`);

        if (!response.ok) {
          throw new Error(`Failed to fetch menu items: ${response.status}`);
        }

        const data: MenuItem[] = await response.json();

        // Filter active items only for the frontend
        const activeItems = data.filter((item) => item.active);

        // Separate items by type
        setCoffeeOptions(
          activeItems.filter((item) => item.type === "coffee_flavor")
        );
        setMilkOptions(
          activeItems.filter((item) => item.type === "milk_option")
        );
      } catch (err) {
        setError(err instanceof Error ? err.message : "Unknown error");
        // TODO: log error
        // Fallback to default options if API fails
        setCoffeeOptions([
          {
            id: 1,
            value: "french_toast",
            label: "French Toast",
            type: "coffee_flavor",
            active: true,
          },
          {
            id: 2,
            value: "dirty_vanilla_chai",
            label: "Dirty Vanilla Chai",
            type: "coffee_flavor",
            active: true,
          },
          {
            id: 3,
            value: "mexican_mocha",
            label: "Mexican Mocha",
            type: "coffee_flavor",
            active: true,
          },
          {
            id: 4,
            value: "cinnamon_brown_sugar",
            label: "Cinnamon Brown Sugar",
            type: "coffee_flavor",
            active: true,
          },
          {
            id: 5,
            value: "horchata",
            label: "Horchata (made w/ rice milk)",
            type: "coffee_flavor",
            active: true,
          },
        ]);
        setMilkOptions([
          {
            id: 6,
            value: "whole",
            label: "Whole Milk",
            type: "milk_option",
            active: true,
          },
          {
            id: 7,
            value: "half_and_half",
            label: "Half & Half",
            type: "milk_option",
            active: true,
          },
          {
            id: 8,
            value: "oat",
            label: "Oat Milk",
            type: "milk_option",
            active: true,
          },
          {
            id: 9,
            value: "almond",
            label: "Almond Milk",
            type: "milk_option",
            active: true,
          },
          {
            id: 10,
            value: "rice",
            label: "Rice Milk",
            type: "milk_option",
            active: true,
          },
        ]);
      } finally {
        setLoading(false);
      }
    }

    fetchMenuItems();
  }, [API_URL]);

  return (
    <MenuContext.Provider
      value={{ coffeeOptions, milkOptions, loading, error }}
    >
      {children}
    </MenuContext.Provider>
  );
};
