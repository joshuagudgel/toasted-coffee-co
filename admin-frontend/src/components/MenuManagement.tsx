import { useState, useEffect } from "react";
import { useMenu, MenuItem } from "../context/MenuContext";
import { MenuItemTable } from "./MenuItemTable";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

interface FormData {
  value: string;
  label: string;
  type: "coffee_flavor" | "milk_option";
  active: boolean;
}

const generateValueFromLabel = (label: string): string => {
  return label
    .toLowerCase()
    .replace(/\s+/g, "_") // Replace spaces with underscores
    .replace(/[^a-z0-9_]/g, "") // Remove special characters
    .replace(/_{2,}/g, "_"); // Replace multiple underscores with single one
};

const MAX_RETRIES = 2;
const RETRY_DELAY = 1000; // ms

export default function MenuManagement() {
  const {
    coffeeItems,
    milkItems,
    loading: menuLoading,
    error: menuError,
    addMenuItem,
    updateMenuItem,
    deleteMenuItem,
    fetchMenuItems,
  } = useMenu();

  const { token, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  const [isAddingItem, setIsAddingItem] = useState(false);
  const [editingItem, setEditingItem] = useState<MenuItem | null>(null);
  const [retryCount, setRetryCount] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [initialLoadAttempted, setInitialLoadAttempted] = useState(false);
  const [formData, setFormData] = useState<FormData>({
    value: "",
    label: "",
    type: "coffee_flavor",
    active: true,
  });

  // Update your useEffect to wait for authentication:
  useEffect(() => {
    const loadData = async () => {
      // Only attempt to load if authentication is confirmed
      if (isAuthenticated && token) {
        try {
          setError(null);
          // Add more delay to ensure token is fully processed
          await new Promise((resolve) => setTimeout(resolve, 500));
          const success = await fetchMenuItems();

          if (!success && retryCount < MAX_RETRIES) {
            // Limit to fewer retries
            setTimeout(() => {
              setRetryCount((prev) => prev + 1);
            }, RETRY_DELAY);
          }
        } catch (err) {
          console.error("Failed to load menu items:", err);
          setError(
            err instanceof Error ? err.message : "Failed to load menu items"
          );
        } finally {
          setInitialLoadAttempted(true);
        }
      } else if (!isAuthenticated) {
        // Only set error if we're definitely not authenticated
        setError("Please log in to access this page");
        setInitialLoadAttempted(true);
      }
    };

    // If we've tried too many times, stop trying
    if (retryCount >= MAX_RETRIES) {
      setError(
        "Could not load menu items after multiple attempts. Please try refreshing the page."
      );
      return;
    }

    // Only load if not already attempted or authentication changed
    if (!initialLoadAttempted || retryCount > 0) {
      loadData();
    }
  }, [
    isAuthenticated,
    token,
    fetchMenuItems,
    retryCount,
    initialLoadAttempted,
  ]);

  // Automatically generate value from label
  useEffect(() => {
    if (formData.label && !editingItem) {
      const generatedValue = generateValueFromLabel(formData.label);
      setFormData((prev) => ({
        ...prev,
        value: generatedValue,
      }));
    }
  }, [formData.label, editingItem]);

  // Form handling
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value, type } = e.target;

    if (type === "checkbox") {
      const checkbox = e.target as HTMLInputElement;
      setFormData((prev) => ({ ...prev, [name]: checkbox.checked }));
    } else {
      setFormData((prev) => ({ ...prev, [name]: value }));
    }
  };

  const resetForm = () => {
    setFormData({
      value: "",
      label: "",
      type: "coffee_flavor",
      active: true,
    });
    setEditingItem(null);
    setIsAddingItem(false);
  };

  const startAddingItem = () => {
    resetForm();
    setIsAddingItem(true);
  };

  const startEditingItem = (item: MenuItem) => {
    setFormData({
      value: item.value,
      label: item.label,
      type: item.type,
      active: item.active,
    });
    setEditingItem(item);
    setIsAddingItem(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isAuthenticated) {
      alert("Your session has expired. Please log in again.");
      navigate("/login");
      return;
    }

    console.log("Sending to API:", formData);
    try {
      if (editingItem) {
        await updateMenuItem(editingItem.id, formData);
        alert("Menu item updated successfully");
      } else {
        await addMenuItem(formData);
        alert("Menu item added successfully");
      }
      resetForm();
    } catch (err) {
      // Handle auth errors specially
      if (err instanceof Error && err.message.includes("Not authenticated")) {
        alert("Your session has expired. Please log in again.");
        navigate("/login");
      } else {
        alert(err instanceof Error ? err.message : "An error occurred");
      }
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm("Are you sure you want to delete this menu item?")) {
      return;
    }

    if (!isAuthenticated) {
      alert("Your session has expired. Please log in again.");
      navigate("/login");
      return;
    }

    try {
      await deleteMenuItem(id);
      alert("Menu item deleted successfully");
    } catch (err) {
      // Handle auth errors specially
      if (err instanceof Error && err.message.includes("Not authenticated")) {
        alert("Your session has expired. Please log in again.");
        navigate("/login");
      } else {
        alert(err instanceof Error ? err.message : "An error occurred");
      }
    }
  };

  // Simplified loading state
  if (menuLoading || !initialLoadAttempted) {
    return (
      <div className="p-8 flex flex-col items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-terracotta mb-4"></div>
        <p className="text-gray-600">Loading menu items...</p>
      </div>
    );
  }

  // Handle errors
  const displayError = error || menuError;
  if (displayError) {
    return (
      <div className="p-4 text-red-500 bg-red-50 rounded-md">
        <h3 className="font-medium">Error</h3>
        <p className="mb-4">{displayError}</p>
        {displayError.includes("Not authenticated") ||
        displayError.includes("session has expired") ? (
          <button
            onClick={() => navigate("/login")}
            className="px-4 py-2 bg-terracotta text-white rounded hover:bg-peach"
          >
            Go to Login
          </button>
        ) : (
          <button
            onClick={() => {
              setError(null);
              setRetryCount(0);
              fetchMenuItems();
            }}
            className="px-4 py-2 bg-terracotta text-white rounded hover:bg-peach"
          >
            Try Again
          </button>
        )}
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Menu Management</h1>
        <button
          onClick={startAddingItem}
          className="px-4 py-2 bg-terracotta text-white rounded hover:bg-peach"
        >
          Add New Item
        </button>
      </div>

      {/* Form for adding/editing items */}
      {(isAddingItem || editingItem) && (
        <div className="bg-white shadow rounded-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">
            {editingItem ? "Edit Menu Item" : "Add New Menu Item"}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Label (display name)
                </label>
                <input
                  type="text"
                  name="label"
                  value={formData.label}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  required
                  placeholder="e.g., French Toast"
                  autoFocus
                />
                <p className="text-xs text-gray-500 mt-1">
                  Human-readable label shown to users
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Value (identifier)
                </label>
                <div className="flex">
                  <input
                    type="text"
                    name="value"
                    value={formData.value}
                    onChange={handleChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                    required
                    placeholder="auto-generated from label"
                  />
                  <button
                    type="button"
                    onClick={() => {
                      const generatedValue = generateValueFromLabel(
                        formData.label
                      );
                      setFormData((prev) => ({
                        ...prev,
                        value: generatedValue,
                      }));
                    }}
                    className="ml-2 px-3 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
                    title="Regenerate from label"
                  >
                    ↻
                  </button>
                </div>
                <p className="text-xs text-gray-500 mt-1">
                  Auto-generated identifier (editable if needed)
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Type
                </label>
                <select
                  name="type"
                  value={formData.type}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  required
                >
                  <option value="coffee_flavor">Coffee Flavor</option>
                  <option value="milk_option">Milk Option</option>
                </select>
              </div>

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="active"
                  name="active"
                  checked={formData.active}
                  onChange={handleChange}
                  className="h-4 w-4 text-terracotta border-gray-300 rounded"
                />
                <label
                  htmlFor="active"
                  className="ml-2 block text-sm text-gray-700"
                >
                  Active (visible to users)
                </label>
              </div>
            </div>

            <div className="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                onClick={resetForm}
                className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="px-4 py-2 bg-terracotta text-white rounded-md hover:bg-peach"
              >
                {editingItem ? "Update Item" : "Add Item"}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Tables */}
      <MenuItemTable
        title="Coffee Flavors"
        items={coffeeItems}
        onEdit={startEditingItem}
        onDelete={handleDelete}
      />
      <MenuItemTable
        title="Milk Options"
        items={milkItems}
        onEdit={startEditingItem}
        onDelete={handleDelete}
      />
    </div>
  );
}
