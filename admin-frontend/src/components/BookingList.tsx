// src/components/BookingList.tsx
import { useState, useEffect } from "react";
import { Booking } from "../types/booking";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";

interface BookingListProps {
  hiddenColumns?: string[];
}

export default function BookingList({ hiddenColumns = [] }: BookingListProps) {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  // TODO authorize API requests with JWT token
  useEffect(() => {
    async function fetchBookings() {
      try {
        // Get JWT token from localStorage
        const token = localStorage.getItem("authToken");
        if (!token) {
          navigate("/signin");
          return;
        }

        const response = await fetch(`${API_URL}/api/v1/bookings`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          localStorage.removeItem("authToken");
          navigate("/signin");
          return;
        }

        if (!response.ok) {
          throw new Error(`Failed to fetch bookings: ${response.status}`);
        }

        const data = await response.json();
        setBookings(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Unknown error");
      } finally {
        setLoading(false);
      }
    }

    if (isAuthenticated) {
      fetchBookings();
    }
  }, [isAuthenticated, navigate]);

  if (loading) return <p>Loading bookings...</p>;
  if (error) return <p className="text-red-500">Error: {error}</p>;

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Bookings</h1>

      <div className="overflow-x-auto">
        <table className="min-w-full bg-white border rounded-lg">
          <thead>
            <tr className="bg-gray-100 text-gray-700">
              {!hiddenColumns.includes("id") && (
                <th className="py-3 px-4 text-left">ID</th>
              )}
              {!hiddenColumns.includes("name") && (
                <th className="py-3 px-4 text-left">Name</th>
              )}
              {!hiddenColumns.includes("contact") && (
                <th className="py-3 px-4 text-left">Contact</th>
              )}
              {!hiddenColumns.includes("date") && (
                <th className="py-3 px-4 text-left">Date</th>
              )}
              {!hiddenColumns.includes("time") && (
                <th className="py-3 px-4 text-left">Time</th>
              )}
              {!hiddenColumns.includes("people") && (
                <th className="py-3 px-4 text-left">People</th>
              )}
              {!hiddenColumns.includes("package") && (
                <th className="py-3 px-4 text-left">Package</th>
              )}
              {!hiddenColumns.includes("location") && (
                <th className="py-3 px-4 text-left">Location</th>
              )}
              {!hiddenColumns.includes("notes") && (
                <th className="py-3 px-4 text-left">Notes</th>
              )}
              {!hiddenColumns.includes("createdAt") && (
                <th className="py-3 px-4 text-left">Created</th>
              )}
            </tr>
          </thead>
          <tbody>
            {bookings.length === 0 ? (
              <tr>
                <td colSpan={7} className="py-4 px-4 text-center text-gray-500">
                  No bookings found
                </td>
              </tr>
            ) : (
              bookings.map((booking) => (
                <tr
                  key={booking.id}
                  className="border-t hover:bg-gray-50 cursor-pointer"
                  onClick={() => navigate(`/booking/${booking.id}`)}
                >
                  {!hiddenColumns.includes("id") && (
                    <td className="py-3 px-4">{booking.id}</td>
                  )}
                  {!hiddenColumns.includes("name") && (
                    <td className="py-3 px-4">{booking.name}</td>
                  )}
                  {!hiddenColumns.includes("contact") && (
                    <td className="py-3 px-4">
                      {booking.email && <div>{booking.email}</div>}
                      {booking.phone && <div>{booking.phone}</div>}
                    </td>
                  )}
                  {!hiddenColumns.includes("date") && (
                    <td className="py-3 px-4">{booking.date}</td>
                  )}
                  {!hiddenColumns.includes("time") && (
                    <td className="py-3 px-4">{booking.time}</td>
                  )}
                  {!hiddenColumns.includes("people") && (
                    <td className="py-3 px-4">{booking.people}</td>
                  )}
                  {!hiddenColumns.includes("package") && (
                    <td className="py-3 px-4">{booking.package || "N/A"}</td>
                  )}
                  {!hiddenColumns.includes("location") && (
                    <td className="py-3 px-4">{booking.location}</td>
                  )}
                  {!hiddenColumns.includes("notes") && (
                    <td className="py-3 px-4">{booking.notes}</td>
                  )}
                  {!hiddenColumns.includes("createdAt") && (
                    <td className="py-3 px-4">{booking.createdAt}</td>
                  )}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
