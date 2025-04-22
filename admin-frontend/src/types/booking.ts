export interface Booking {
  id: number;
  name: string;
  date: string;
  time: string;
  people: number;
  coffeeFlavors: string[];
  milkOptions: string[];
  location: string;
  notes: string;
  package?: string;
  createdAt: string;
}