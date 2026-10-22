import { useState } from "react";
import "./styles.css";

type Item = {
  id: string;
  name: string;
  quantity: number;
  unitPrice: number; // dollars, 0 if unknown
  aisle: string;
  purchased: boolean;
};

const AISLES: string[] = ["Produce", "Dairy", "Bakery", "Pantry"];

const SEED_ITEMS: Item[] = [
  { id: "1", name: "Bananas", quantity: 6, unitPrice: 0.25, aisle: "Produce", purchased: false },
  { id: "2", name: "Milk", quantity: 1, unitPrice: 3.49, aisle: "Dairy", purchased: false },
  { id: "3", name: "Sourdough loaf", quantity: 1, unitPrice: 5, aisle: "Bakery", purchased: true },
  { id: "4", name: "Spaghetti", quantity: 2, unitPrice: 1.89, aisle: "Pantry", purchased: false },
];

export default function App() {
  const [items] = useState<Item[]>(SEED_ITEMS);

  return (
    <main className="app">
      <h1>Shopping list</h1>
      <p>
        {items.length} items in {AISLES.length} aisles.
      </p>
    </main>
  );
}
