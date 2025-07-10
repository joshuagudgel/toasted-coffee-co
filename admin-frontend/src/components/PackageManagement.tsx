import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { usePackages, Package, PackageInput } from "../context/PackageContext";
import { useAuth } from "../context/AuthContext";

// Form data for adding/editing packages
interface FormData {
  name: string;
  price: string;
  description: string;
  points: string[];
  displayOrder: number;
  active: boolean;
}

export default function PackageManagement() {
  const {
    packages,
    loading: packagesLoading,
    error: packagesError,
    fetchPackages,
    addPackage,
    updatePackage,
    deletePackage,
  } = usePackages();

  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  // Track if a fetch operation is in progress
  const isFetchingRef = useRef(false);

  const [isAddingItem, setIsAddingItem] = useState(false);
  const [editingItem, setEditingItem] = useState<Package | null>(null);
  const [retryCount, setRetryCount] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [initialLoadAttempted, setInitialLoadAttempted] = useState(false);
  const [formData, setFormData] = useState<FormData>({
    name: "",
    price: "",
    description: "",
    points: [""],
    displayOrder: 0,
    active: true,
  });

  const MAX_RETRIES = 3;
  const RETRY_DELAY = 2000;

  useEffect(() => {
    const loadData = async () => {
      // Prevent concurrent fetch operations
      if (isFetchingRef.current) return;

      // Only attempt to load if authenticated
      if (isAuthenticated) {
        try {
          isFetchingRef.current = true;
          setError(null);
          const success = await fetchPackages(true); // Include inactive packages

          if (!success && retryCount < MAX_RETRIES) {
            // Limit to fewer retries
            setTimeout(() => {
              setRetryCount((prev) => prev + 1);
            }, RETRY_DELAY);
          }
        } catch (err) {
          console.error("Failed to load packages:", err);
          // Handle auth errors
          if (
            err instanceof Error &&
            (err.message === "Not authenticated" ||
              err.message === "Session expired")
          ) {
            setError("Your session has expired. Please sign in again.");
          } else {
            setError(
              err instanceof Error ? err.message : "Failed to load packages"
            );
          }
        } finally {
          isFetchingRef.current = false;
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
        "Could not load packages after multiple attempts. Please try refreshing the page."
      );
      return;
    }

    // Only load if not already attempted or authentication changed
    if (!initialLoadAttempted || retryCount > 0) {
      loadData();
    }
  }, [isAuthenticated, fetchPackages, retryCount, initialLoadAttempted]);

  // Form handling
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value, type } = e.target;

    if (type === "checkbox") {
      const checkbox = e.target as HTMLInputElement;
      setFormData((prev) => ({ ...prev, [name]: checkbox.checked }));
    } else {
      setFormData((prev) => ({ ...prev, [name]: value }));
    }
  };

  // Handle points array changes
  const handlePointChange = (index: number, value: string) => {
    setFormData((prev) => {
      const newPoints = [...prev.points];
      newPoints[index] = value;
      return { ...prev, points: newPoints };
    });
  };

  const addPoint = () => {
    setFormData((prev) => ({
      ...prev,
      points: [...prev.points, ""],
    }));
  };

  const removePoint = (index: number) => {
    setFormData((prev) => {
      const newPoints = [...prev.points];
      newPoints.splice(index, 1);
      return { ...prev, points: newPoints.length ? newPoints : [""] };
    });
  };

  const resetForm = () => {
    setFormData({
      name: "",
      price: "",
      description: "",
      points: [""],
      displayOrder: 0,
      active: true,
    });
    setEditingItem(null);
    setIsAddingItem(false);
  };

  const startAddingItem = () => {
    resetForm();
    setIsAddingItem(true);
  };

  const startEditingItem = (item: Package) => {
    setFormData({
      name: item.name,
      price: item.price,
      description: item.description,
      points: item.points.length ? item.points : [""],
      displayOrder: item.displayOrder,
      active: item.active,
    });
    setEditingItem(item);
    setIsAddingItem(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isAuthenticated) {
      alert("Your session has expired. Please log in again.");
      navigate("/signin");
      return;
    }

    // Filter out empty points
    const packageData: PackageInput = {
      ...formData,
      displayOrder: Number(formData.displayOrder),
      points: formData.points.filter((point) => point.trim() !== ""),
    };

    try {
      if (editingItem) {
        await updatePackage(editingItem.id, packageData);
        alert("Package updated successfully");
      } else {
        await addPackage(packageData);
        alert("Package added successfully");
      }
      resetForm();
    } catch (err) {
      // Handle auth errors consistently
      if (
        err instanceof Error &&
        (err.message === "Not authenticated" ||
          err.message === "Session expired")
      ) {
        alert("Your session has expired. Please log in again.");
        navigate("/signin");
      } else {
        alert(err instanceof Error ? err.message : "An error occurred");
      }
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm("Are you sure you want to delete this package?")) {
      return;
    }

    if (!isAuthenticated) {
      alert("Your session has expired. Please log in again.");
      navigate("/signin");
      return;
    }

    try {
      await deletePackage(id);
      alert("Package deleted successfully");
    } catch (err) {
      // Handle auth errors consistently
      if (
        err instanceof Error &&
        (err.message === "Not authenticated" ||
          err.message === "Session expired")
      ) {
        alert("Your session has expired. Please log in again.");
        navigate("/signin");
      } else {
        alert(err instanceof Error ? err.message : "An error occurred");
      }
    }
  };

  // Simplified loading state
  if (packagesLoading || !initialLoadAttempted) {
    return (
      <div className="p-8 flex flex-col items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-terracotta mb-4"></div>
        <p className="text-gray-600">Loading packages...</p>
      </div>
    );
  }

  // Handle errors
  const displayError = error || packagesError;
  if (displayError) {
    return (
      <div className="p-4 text-red-500 bg-red-50 rounded-md">
        <h3 className="font-medium">Error</h3>
        <p className="mb-4">{displayError}</p>
        {displayError.includes("Not authenticated") ||
        displayError.includes("session has expired") ? (
          <button
            onClick={() => navigate("/signin")}
            className="px-4 py-2 bg-terracotta text-white rounded hover:bg-peach"
          >
            Go to Login
          </button>
        ) : (
          <button
            onClick={() => {
              setError(null);
              setRetryCount(0);
              fetchPackages(true);
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
        <h1 className="text-2xl font-bold">Package Management</h1>
        <button
          onClick={startAddingItem}
          className="px-4 py-2 bg-terracotta text-white rounded hover:bg-peach"
        >
          Add New Package
        </button>
      </div>

      {/* Form for adding/editing packages */}
      {(isAddingItem || editingItem) && (
        <div className="bg-white shadow rounded-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">
            {editingItem ? "Edit Package" : "Add New Package"}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Package Name
                </label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  required
                  placeholder="e.g., Basic Package"
                  autoFocus
                />
              </div>

              {/* Price */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Price
                </label>
                <input
                  type="text"
                  name="price"
                  value={formData.price}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  required
                  placeholder="e.g., $199"
                />
              </div>

              {/* Description */}
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  name="description"
                  value={formData.description}
                  onChange={handleChange}
                  rows={2}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  required
                  placeholder="Brief description of the package"
                ></textarea>
              </div>

              {/* Points/Features */}
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Features/Points
                </label>
                {formData.points.map((point, index) => (
                  <div key={index} className="flex mb-2">
                    <input
                      type="text"
                      value={point}
                      onChange={(e) => handlePointChange(index, e.target.value)}
                      className="flex-grow px-3 py-2 border border-gray-300 rounded-md"
                      placeholder={`Feature ${index + 1}`}
                    />
                    <button
                      type="button"
                      onClick={() => removePoint(index)}
                      className="ml-2 px-3 py-2 bg-red-100 text-red-700 rounded hover:bg-red-200"
                      disabled={formData.points.length <= 1}
                    >
                      Remove
                    </button>
                  </div>
                ))}
                <button
                  type="button"
                  onClick={addPoint}
                  className="mt-2 px-3 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200"
                >
                  Add Feature
                </button>
              </div>

              {/* Display Order */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Display Order
                </label>
                <input
                  type="number"
                  name="displayOrder"
                  value={formData.displayOrder}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  min="0"
                />
                <p className="text-xs text-gray-500 mt-1">
                  Lower numbers appear first. Packages with the same order are
                  sorted by name.
                </p>
              </div>

              {/* Active Status */}
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
                {editingItem ? "Update Package" : "Add Package"}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Package List */}
      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-4">Packages</h2>
        {packages.length === 0 ? (
          <p className="text-gray-500 text-center py-4">No packages found</p>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {packages.map((pkg) => (
              <div
                key={pkg.id}
                className={`border rounded-lg overflow-hidden ${
                  !pkg.active ? "opacity-60" : ""
                }`}
              >
                <div className="bg-terracotta p-4 text-white">
                  <div className="flex justify-between items-center">
                    <h3 className="font-bold">{pkg.name}</h3>
                    <span className="font-bold">{pkg.price}</span>
                  </div>
                </div>
                <div className="p-4">
                  <p className="text-gray-600 mb-4">{pkg.description}</p>
                  <h4 className="font-medium mb-2">Features:</h4>
                  <ul className="list-disc pl-5 mb-4">
                    {pkg.points.map((point, index) => (
                      <li key={index} className="text-gray-700">
                        {point}
                      </li>
                    ))}
                  </ul>
                  <div className="flex justify-between items-center mt-4">
                    <span
                      className={`px-2 py-1 text-xs rounded-full ${
                        pkg.active
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                      }`}
                    >
                      {pkg.active ? "Active" : "Inactive"}
                    </span>
                    <span className="text-xs text-gray-500">
                      Order: {pkg.displayOrder}
                    </span>
                    <div>
                      <button
                        onClick={() => startEditingItem(pkg)}
                        className="text-blue-600 hover:text-blue-800 mr-3"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(pkg.id)}
                        className="text-red-600 hover:text-red-800"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
