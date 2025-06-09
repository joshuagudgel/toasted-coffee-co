import { useState } from "react";
import { useMenu, MenuItem } from "../context/MenuContext";
import { MenuItemTable } from "./MenuItemTable";

interface FormData {
  value: string;
  label: string;
  type: "coffee_flavor" | "milk_option";
  active: boolean;
}

export default function MenuManagement() {
  const {
    coffeeItems,
    milkItems,
    loading,
    error,
    addMenuItem,
    updateMenuItem,
    deleteMenuItem,
  } = useMenu();

  const [isAddingItem, setIsAddingItem] = useState(false);
  const [editingItem, setEditingItem] = useState<MenuItem | null>(null);
  const [formData, setFormData] = useState<FormData>({
    value: "",
    label: "",
    type: "coffee_flavor",
    active: true,
  });

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
      alert(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm("Are you sure you want to delete this menu item?")) {
      return;
    }

    try {
      await deleteMenuItem(id);
      alert("Menu item deleted successfully");
    } catch (err) {
      alert(err instanceof Error ? err.message : "An error occurred");
    }
  };

  if (loading) return <div className="p-4">Loading menu items...</div>;
  if (error) return <div className="p-4 text-red-500">Error: {error}</div>;

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
                  Value (identifier)
                </label>
                <input
                  type="text"
                  name="value"
                  value={formData.value}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  required
                  placeholder="e.g., french_toast"
                />
                <p className="text-xs text-gray-500 mt-1">
                  Unique identifier (no spaces, use underscores)
                </p>
              </div>

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
                />
                <p className="text-xs text-gray-500 mt-1">
                  Human-readable label shown to users
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
