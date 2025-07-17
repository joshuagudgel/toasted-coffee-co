export interface Booking {
  id: number;
  name: string;
  email?: string;
  phone?: string;
  date: string;
  time: string;
  people: number;
  coffeeFlavors: string[];
  milkOptions: string[];
  location: string;
  notes: string;
  package?: string;
  createdAt: string;
  archived: boolean;
  isOutdoor: boolean;
  hasShade: boolean;
}