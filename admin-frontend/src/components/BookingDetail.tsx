import { useState, useEffect } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { Booking } from "../types/booking";

export default function BookingDetail() {
  const { id } = useParams<{ id: string }>();
  const [booking, setBooking] = useState<Booking | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchBookingDetails() {
      try {
        const token = localStorage.getItem("authToken");
        if (!token) {
          setError("Not authenticated");
          setLoading(false);
          return;
        }
        const response = await fetch(`${API_URL}/api/v1/bookings/${id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (response.status === 401) {
          setError("Unauthorized. Please sign in again.");
          setLoading(false);
          return;
        }
        if (!response.ok) {
          throw new Error(`Failed to fetch booking: ${response.status}`);
        }
        const data = await response.json();
        setBooking(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Unknown error");
      } finally {
        setLoading(false);
      }
    }

    fetchBookingDetails();
  }, [id, API_URL]);

  // Handle deletion
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
      const token = localStorage.getItem("authToken");
      if (!token) {
        alert("Not authenticated");
        setIsDeleting(false);
        return;
      }

      const response = await fetch(`${API_URL}/api/v1/bookings/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status === 401) {
        alert("Unauthorized. Please sign in again.");
        setIsDeleting(false);
        return;
      }

      if (!response.ok) {
        throw new Error(`Failed to delete booking: ${response.status}`);
      }

      alert("Booking deleted successfully");
      navigate("/"); // Redirect to booking list
    } catch (err) {
      alert(err instanceof Error ? err.message : "Error deleting booking");
    } finally {
      setIsDeleting(false);
    }
  };

  if (loading) return <p>Loading booking details...</p>;
  if (error) return <p className="text-red-500">Error: {error}</p>;
  if (!booking) return <p>Booking not found</p>;

  return (
    <div>
      <Link
        to="/"
        className="text-blue-600 hover:text-blue-800 flex items-center mb-6"
      >
        ← Back to Bookings
      </Link>

      <h1 className="text-2xl font-bold mb-6">Booking #{booking.id}</h1>

      <div className="bg-white shadow rounded-lg p-6">
        <dl className="grid grid-cols-1 md:grid-cols-2 gap-x-4 gap-y-6">
          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">Customer Name</dt>
            <dd className="mt-1 text-lg text-gray-900">{booking.name}</dd>
          </div>

          <div>
            <dt className="text-sm font-medium text-gray-500">
              Contact Information
            </dt>
            <dd className="mt-1 text-gray-900">
              {booking.email && <div>Email: {booking.email}</div>}
              {booking.phone && <div>Phone: {booking.phone}</div>}
            </dd>
          </div>

          <div>
            <dt className="text-sm font-medium text-gray-500">Date</dt>
            <dd className="mt-1 text-gray-900">{booking.date}</dd>
          </div>

          <div>
            <dt className="text-sm font-medium text-gray-500">Time</dt>
            <dd className="mt-1 text-gray-900">{booking.time}</dd>
          </div>

          <div>
            <dt className="text-sm font-medium text-gray-500">People</dt>
            <dd className="mt-1 text-gray-900">{booking.people}</dd>
          </div>

          <div>
            <dt className="text-sm font-medium text-gray-500">Location</dt>
            <dd className="mt-1 text-gray-900">{booking.location}</dd>
          </div>

          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">
              Coffee Flavors
            </dt>
            <dd className="mt-1">
              <ul className="list-disc pl-5 text-gray-900">
                {booking.coffeeFlavors.map((flavor, index) => (
                  <li key={index}>{flavor}</li>
                ))}
              </ul>
            </dd>
          </div>

          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">Milk Options</dt>
            <dd className="mt-1">
              <ul className="list-disc pl-5 text-gray-900">
                {booking.milkOptions.map((option, index) => (
                  <li key={index}>{option}</li>
                ))}
              </ul>
            </dd>
          </div>

          {booking.package && (
            <div>
              <dt className="text-sm font-medium text-gray-500">Package</dt>
              <dd className="mt-1 text-gray-900">{booking.package}</dd>
            </div>
          )}

          <div className="col-span-2">
            <dt className="text-sm font-medium text-gray-500">Notes</dt>
            <dd className="mt-1 text-gray-900 whitespace-pre-line">
              {booking.notes || "No notes provided"}
            </dd>
          </div>
        </dl>
      </div>

      {/* Delete button at the bottom */}
      <div className="mt-6 flex justify-end">
        <button
          onClick={handleDelete}
          disabled={isDeleting}
          className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50"
        >
          {isDeleting ? "Deleting..." : "Delete Booking"}
        </button>
      </div>
    </div>
  );
}
