import { useState, useEffect } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { Booking } from "../types/booking";
import { useAuth } from "../context/AuthContext";
import { useMenu } from "../context/MenuContext";

export default function BookingDetail() {
  const { id } = useParams<{ id: string }>();
  const [booking, setBooking] = useState<Booking | null>(null);
  const [editedBooking, setEditedBooking] = useState<Booking | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isArchiving, setIsArchiving] = useState(false);
  const { apiRequest } = useAuth(); // Use apiRequest from AuthContext
  const navigate = useNavigate();

  // Get menu options from context
  const { coffeeItems, milkItems } = useMenu();

  // Convert menu items to options format
  const coffeeOptions = coffeeItems.map((item) => ({
    value: item.value,
    label: item.label,
  }));

  const milkOptions = milkItems.map((item) => ({
    value: item.value,
    label: item.label,
  }));

  // Fetch booking details
  useEffect(() => {
    async function fetchBookingDetails() {
      setLoading(true);
      try {
        const data = await apiRequest<Booking>(`/api/v1/bookings/${id}`);
        setBooking(data);
        // Initialize editedBooking with the same data
        setEditedBooking(data);
      } catch (err) {
        if (
          err instanceof Error &&
          (err.message === "Not authenticated" ||
            err.message === "Session expired")
        ) {
          setError("Your session has expired. Please sign in again.");
        } else {
          setError(err instanceof Error ? err.message : "Unknown error");
        }
      } finally {
        setLoading(false);
      }
    }

    fetchBookingDetails();
  }, [id, apiRequest]);

  // Handle form input changes
  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    if (!editedBooking) return;

    const { name, value } = e.target;
    // Convert people field to a number
    if (name === "people") {
      setEditedBooking((prev) => {
        if (!prev) return null;
        return { ...prev, [name]: parseInt(value, 10) || 0 };
      });
    } else {
      setEditedBooking((prev) => {
        if (!prev) return null;
        return { ...prev, [name]: value };
      });
    }
  };

  // Handle checkbox changes for arrays (coffeeFlavors, milkOptions)
  const handleCheckBoxChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    field: "coffeeFlavors" | "milkOptions"
  ) => {
    if (!editedBooking) return;

    const { value, checked } = e.target;
    setEditedBooking((prev) => {
      if (!prev) return null;

      if (checked) {
        // Add value to array
        return {
          ...prev,
          [field]: [...prev[field], value],
        };
      } else {
        // Remove value from array
        return {
          ...prev,
          [field]: prev[field].filter((item) => item !== value),
        };
      }
    });
  };

  // Start editing mode
  const handleEdit = () => {
    setIsEditing(true);
  };

  // Cancel editing and revert changes
  const handleCancel = () => {
    setEditedBooking(booking);
    setIsEditing(false);
  };

  // Save changes
  const handleSave = async () => {
    if (!editedBooking) return;

    // Validation (same as BookingModal)
    if (editedBooking.email === "" && editedBooking.phone === "") {
      alert("Please provide at least one contact method (email or phone).");
      return;
    }

    if (editedBooking.coffeeFlavors.length === 0) {
      alert("Please select at least one coffee flavor.");
      return;
    }

    if (editedBooking.milkOptions.length === 0) {
      alert("Please select at least one milk option.");
      return;
    }

    setIsSaving(true);
    try {
      // Use apiRequest instead of fetch
      await apiRequest(`/api/v1/bookings/${id}`, {
        method: "PUT",
        body: JSON.stringify(editedBooking),
      });

      // Update booking state with edited values
      setBooking(editedBooking);
      setIsEditing(false);
      alert("Booking updated successfully");
    } catch (err) {
      alert(err instanceof Error ? err.message : "Error updating booking");
    } finally {
      setIsSaving(false);
    }
  };

  // Delete booking function
  const handleDelete = async () => {
    if (
      !window.confirm(
        "Are you sure you want to delete this booking? This action cannot be undone."
      )
    ) {
      return;
    }

    setIsDeleting(true);
    try {
      // Use apiRequest instead of fetch
      await apiRequest(`/api/v1/bookings/${id}`, {
        method: "DELETE",
      });

      alert("Booking deleted successfully");
      navigate("/"); // Redirect to booking list
    } catch (err) {
      alert(err instanceof Error ? err.message : "Error deleting booking");
    } finally {
      setIsDeleting(false);
    }
  };

  // Archive/unarchive booking function
  const handleArchiveToggle = async () => {
    if (
      !window.confirm(
        booking?.archived
          ? "Are you sure you want to unarchive this booking?"
          : "Are you sure you want to archive this booking?"
      )
    ) {
      return;
    }

    setIsArchiving(true);
    try {
      const endpoint = booking?.archived
        ? `/api/v1/bookings/${id}/unarchive`
        : `/api/v1/bookings/${id}/archive`;

      // Use apiRequest instead of fetch
      await apiRequest(endpoint, {
        method: "POST",
      });

      alert(
        `Booking ${booking?.archived ? "unarchived" : "archived"} successfully`
      );
      navigate("/"); // Redirect to booking list
    } catch (err) {
      alert(err instanceof Error ? err.message : "Error updating booking");
    } finally {
      setIsArchiving(false);
    }
  };

  // Show loading state
  if (loading) {
    return (
      <div className="p-8 flex flex-col items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-terracotta mb-4"></div>
        <p className="text-gray-600">Loading booking details...</p>
      </div>
    );
  }
  if (error) {
    return (
      <div className="p-4 text-red-500 bg-red-50 rounded-md">
        <h3 className="font-medium">Error</h3>
        <p className="mb-4">{error}</p>
        <button
          onClick={() => navigate("/")}
          className="px-4 py-2 bg-terracotta text-white rounded hover:bg-peach"
        >
          Return to Bookings
        </button>
      </div>
    );
  }
  if (!booking || !editedBooking) return <p>Booking not found</p>;

  return (
    <div>
      <Link
        to="/"
        className="text-blue-600 hover:text-blue-800 flex items-center mb-6"
      >
        ← Back to Bookings
      </Link>

      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Booking #{booking.id}</h1>

        {/* Edit/Save/Cancel buttons */}
        {!isEditing ? (
          <button
            onClick={handleEdit}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Edit Booking
          </button>
        ) : (
          <div className="space-x-2">
            <button
              onClick={handleCancel}
              className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
              disabled={isSaving}
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
              disabled={isSaving}
            >
              {isSaving ? "Saving..." : "Save Changes"}
            </button>
          </div>
        )}
      </div>

      <div className="bg-white shadow rounded-lg p-6">
        {/* View/Edit Form */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-x-4 gap-y-6">
          {/* Customer Name */}
          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">Customer Name</dt>
            {isEditing ? (
              <input
                type="text"
                name="name"
                value={editedBooking.name}
                onChange={handleChange}
                className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md"
                required
              />
            ) : (
              <dd className="mt-1 text-lg text-gray-900">{booking.name}</dd>
            )}
          </div>

          {/* Contact Information */}
          <div>
            <dt className="text-sm font-medium text-gray-500">
              Contact Information
            </dt>
            {isEditing ? (
              <div className="space-y-2 mt-1">
                <div>
                  <label className="block text-xs text-gray-500">Email</label>
                  <input
                    type="email"
                    name="email"
                    value={editedBooking.email || ""}
                    onChange={handleChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div>
                  <label className="block text-xs text-gray-500">Phone</label>
                  <input
                    type="tel"
                    name="phone"
                    value={editedBooking.phone || ""}
                    onChange={handleChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  />
                </div>
              </div>
            ) : (
              <dd className="mt-1 text-gray-900">
                {booking.email && <div>Email: {booking.email}</div>}
                {booking.phone && <div>Phone: {booking.phone}</div>}
              </dd>
            )}
          </div>

          {/* Date */}
          <div>
            <dt className="text-sm font-medium text-gray-500">Date</dt>
            {isEditing ? (
              <input
                type="date"
                name="date"
                value={editedBooking.date}
                onChange={handleChange}
                className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md"
                required
              />
            ) : (
              <dd className="mt-1 text-gray-900">{booking.date}</dd>
            )}
          </div>

          {/* Time */}
          <div>
            <dt className="text-sm font-medium text-gray-500">Time</dt>
            {isEditing ? (
              <input
                type="time"
                name="time"
                value={editedBooking.time}
                onChange={handleChange}
                className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md"
                required
              />
            ) : (
              <dd className="mt-1 text-gray-900">{booking.time}</dd>
            )}
          </div>

          {/* People */}
          <div>
            <dt className="text-sm font-medium text-gray-500">People</dt>
            {isEditing ? (
              <input
                type="number"
                name="people"
                min="1"
                value={editedBooking.people}
                onChange={handleChange}
                className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md"
                required
              />
            ) : (
              <dd className="mt-1 text-gray-900">{booking.people}</dd>
            )}
          </div>

          {/* Location */}
          <div>
            <dt className="text-sm font-medium text-gray-500">Location</dt>
            {isEditing ? (
              <input
                type="text"
                name="location"
                value={editedBooking.location}
                onChange={handleChange}
                className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md"
                required
              />
            ) : (
              <dd className="mt-1 text-gray-900">{booking.location}</dd>
            )}
          </div>

          {/* Coffee Flavors */}
          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">
              Coffee Flavors
            </dt>
            {isEditing ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 mt-1 bg-white p-3 rounded-md border border-gray-300">
                {coffeeOptions.map((option) => (
                  <div key={option.value} className="flex items-center mb-2">
                    <input
                      type="checkbox"
                      id={`coffee-${option.value}`}
                      value={option.value}
                      checked={editedBooking.coffeeFlavors.includes(
                        option.value
                      )}
                      onChange={(e) => handleCheckBoxChange(e, "coffeeFlavors")}
                      className="h-4 w-4 border-gray-300 rounded mr-2"
                    />
                    <label
                      htmlFor={`coffee-${option.value}`}
                      className="text-gray-700"
                    >
                      {option.label}
                    </label>
                  </div>
                ))}
              </div>
            ) : (
              <dd className="mt-1">
                <ul className="list-disc pl-5 text-gray-900">
                  {booking.coffeeFlavors.map((flavor, index) => (
                    <li key={index}>
                      {coffeeOptions.find((o) => o.value === flavor)?.label ||
                        flavor}
                    </li>
                  ))}
                </ul>
              </dd>
            )}
          </div>

          {/* Milk Options */}
          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">Milk Options</dt>
            {isEditing ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 mt-1 bg-white p-3 rounded-md border border-gray-300">
                {milkOptions.map((option) => (
                  <div key={option.value} className="flex items-center mb-2">
                    <input
                      type="checkbox"
                      id={`milk-${option.value}`}
                      value={option.value}
                      checked={editedBooking.milkOptions.includes(option.value)}
                      onChange={(e) => handleCheckBoxChange(e, "milkOptions")}
                      className="h-4 w-4 border-gray-300 rounded mr-2"
                    />
                    <label
                      htmlFor={`milk-${option.value}`}
                      className="text-gray-700"
                    >
                      {option.label}
                    </label>
                  </div>
                ))}
              </div>
            ) : (
              <dd className="mt-1">
                <ul className="list-disc pl-5 text-gray-900">
                  {booking.milkOptions.map((option, index) => (
                    <li key={index}>
                      {milkOptions.find((o) => o.value === option)?.label ||
                        option}
                    </li>
                  ))}
                </ul>
              </dd>
            )}
          </div>

          {/* Package */}
          {(booking.package || isEditing) && (
            <div>
              <dt className="text-sm font-medium text-gray-500">Package</dt>
              {isEditing ? (
                <select
                  name="package"
                  value={editedBooking.package || ""}
                  onChange={handleChange}
                  className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md"
                >
                  <option value="">No Package</option>
                  <option value="Group">Group</option>
                  <option value="Crowd">Crowd</option>
                  <option value="Party">Party</option>
                </select>
              ) : (
                <dd className="mt-1 text-gray-900">{booking.package}</dd>
              )}
            </div>
          )}

          {/* Outdoor Options */}
          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">Outdoor Setup</dt>
            {isEditing ? (
              <div className="mt-1 space-y-2">
                <div className="flex items-center">
                  <input
                    type="checkbox"
                    id="is-outdoor"
                    name="isOutdoor"
                    checked={editedBooking.isOutdoor}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setEditedBooking((prev) => {
                        if (!prev) return null;
                        // If turning off outdoor, also turn off shade
                        return {
                          ...prev,
                          isOutdoor: checked,
                          hasShade: checked ? prev.hasShade : false,
                        };
                      });
                    }}
                    className="h-4 w-4 border-gray-300 rounded mr-2"
                  />
                  <label htmlFor="is-outdoor" className="text-gray-700">
                    Outdoor Event
                  </label>
                </div>

                {editedBooking.isOutdoor && (
                  <div className="flex items-center ml-6">
                    <input
                      type="checkbox"
                      id="has-shade"
                      name="hasShade"
                      checked={editedBooking.hasShade}
                      onChange={(e) => {
                        setEditedBooking((prev) => {
                          if (!prev) return null;
                          return {
                            ...prev,
                            hasShade: e.target.checked,
                          };
                        });
                      }}
                      disabled={!editedBooking.isOutdoor}
                      className="h-4 w-4 border-gray-300 rounded mr-2"
                    />
                    <label htmlFor="has-shade" className="text-gray-700">
                      Shaded Area Available
                    </label>
                  </div>
                )}
              </div>
            ) : (
              <dd className="mt-1 text-gray-900">
                {booking.isOutdoor ? (
                  <>
                    <div className="flex items-center">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-5 w-5 text-green-600 mr-2"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      Outdoor Event
                    </div>
                    <div className="ml-7 mt-1">
                      {booking.hasShade ? (
                        <span className="text-green-600">Shaded area available</span>
                      ) : (
                        <span className="text-amber-600">No shade available</span>
                      )}
                    </div>
                  </>
                ) : (
                  <div className="flex items-center">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5 text-blue-600 mr-2"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      />
                    </svg>
                    Indoor Event
                  </div>
                )}
              </dd>
            )}
          </div>

          {/* Notes */}
          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">Notes</dt>
            {isEditing ? (
              <textarea
                name="notes"
                value={editedBooking.notes}
                onChange={handleChange}
                rows={3}
                className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md"
              ></textarea>
            ) : (
              <dd className="mt-1 text-gray-900 whitespace-pre-line">
                {booking.notes || "No notes provided"}
              </dd>
            )}
          </div>
        </div>
      </div>

      {/* Delete button - only shown when not in edit mode */}
      {!isEditing && (
        <div className="mt-6 flex justify-end space-x-2">
          <button
            onClick={handleArchiveToggle}
            disabled={isArchiving}
            className={`px-4 py-2 text-white rounded disabled:opacity-50 flex items-center ${
              booking.archived
                ? "bg-green-600 hover:bg-green-700"
                : "bg-amber-600 hover:bg-amber-700"
            }`}
          >
            {isArchiving ? (
              <>
                <div className="animate-spin h-4 w-4 border-b-2 border-white mr-2"></div>
                <span>
                  {booking.archived ? "Unarchiving..." : "Archiving..."}
                </span>
              </>
            ) : booking.archived ? (
              "Unarchive Booking"
            ) : (
              "Archive Booking"
            )}
          </button>
          <button
            onClick={handleDelete}
            disabled={isDeleting}
            className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50 flex items-center"
          >
            {isDeleting ? (
              <>
                <div className="animate-spin h-4 w-4 border-b-2 border-white mr-2"></div>
                <span>Deleting...</span>
              </>
            ) : (
              "Delete Booking"
            )}
          </button>
        </div>
      )}
    </div>
  );
}
