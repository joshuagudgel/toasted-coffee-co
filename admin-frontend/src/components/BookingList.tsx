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
  const [includeArchived, setIncludeArchived] = useState(false);
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    console.log("BookingList rendering");
    async function fetchBookings() {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(
          `${API_URL}/api/v1/bookings?include_archived=${includeArchived}`,
          {
            credentials: "include",
          }
        );

        if (response.status === 401) {
          setError("Your session has expired");
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
  }, [isAuthenticated, navigate, includeArchived]);

  const toggleArchivedView = () => {
    setIncludeArchived((prev) => !prev);
  };

  if (loading) {
    return (
      <div className="p-8 flex flex-col items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-terracotta mb-4"></div>
        <p className="text-gray-600">Loading bookings...</p>
      </div>
    );
  }
  if (error) return <p className="text-red-500">Error: {error}</p>;

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Bookings</h1>
        <div className="flex space-x-2">
          <button
            onClick={toggleArchivedView}
            className="px-3 py-1 bg-gray-200 rounded hover:bg-gray-300"
          >
            {includeArchived ? "Hide Archived" : "Show Archived"}
          </button>
        </div>
      </div>

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
                  className={`border-t hover:bg-gray-50 cursor-pointer ${
                    booking.archived ? "bg-gray-100 text-gray-500" : ""
                  }`}
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
