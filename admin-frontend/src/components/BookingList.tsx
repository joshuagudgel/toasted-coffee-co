// src/components/BookingList.tsx
import { useState, useEffect } from "react";
import { Booking } from "../types/booking";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";

export default function BookingList() {
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
            'Authorization': `Bearer ${token}`
          }
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

    if( isAuthenticated) {
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
              <th className="py-3 px-4 text-left">ID</th>
              <th className="py-3 px-4 text-left">Name</th>
              <th className="py-3 px-4 text-left">Date</th>
              <th className="py-3 px-4 text-left">Time</th>
              <th className="py-3 px-4 text-left">People</th>
              <th className="py-3 px-4 text-left">Package</th>
              <th className="py-3 px-4 text-left">Location</th>
              <th className="py-3 px-4 text-left">Notes</th>
              <th className="py-3 px-4 text-left">Created</th>
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
                <tr key={booking.id} className="border-t hover:bg-gray-50">
                  <td className="py-3 px-4">{booking.id}</td>
                  <td className="py-3 px-4">{booking.name}</td>
                  <td className="py-3 px-4">{booking.date}</td>
                  <td className="py-3 px-4">{booking.time}</td>
                  <td className="py-3 px-4">{booking.people}</td>
                  <td className="py-3 px-4">{booking.package || "N/A"}</td>
                  <td className="py-3 px-4">{booking.location}</td>
                  <td className="py-3 px-4">{booking.notes}</td>
                  <td className="py-3 px-4">{booking.createdAt}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
