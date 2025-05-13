import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import { Booking } from "../types/booking";

export default function BookingDetail() {
  const { id } = useParams<{ id: string }>();
  const [booking, setBooking] = useState<Booking | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  useEffect(() => {
    async function fetchBookingDetails() {
      try {
        const response = await fetch(`${API_URL}/api/v1/bookings/${id}`);

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
  }, [id]);

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
    </div>
  );
}
