import { BrowserRouter, Routes, Route } from "react-router-dom";
import Layout from "./components/Layout";
import BookingList from "./components/BookingList";
import BookingDetail from "./components/BookingDetail";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<BookingList />} />
          <Route path="booking/:id" element={<BookingDetail />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
