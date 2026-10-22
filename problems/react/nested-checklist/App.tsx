import { useState } from "react";
import "./styles.css";

// A flat, ordered list. Hierarchy is implied by `depth`: an item's children are
// the items that follow it with a greater depth, up to the next item at <= its depth.
type ChecklistItem = {
  id: string;
  text: string;
  checked: boolean;
  depth: number; // 0 = top level
};

const SEED_ITEMS: ChecklistItem[] = [
  { id: "1", text: "Plan the trip", checked: false, depth: 0 },
  { id: "2", text: "Book flights", checked: true, depth: 1 },
  { id: "3", text: "Reserve hotel", checked: false, depth: 1 },
  { id: "4", text: "Pack", checked: false, depth: 0 },
];

export default function App() {
  const [items] = useState<ChecklistItem[]>(SEED_ITEMS);

  return (
    <main className="app">
      <h1>Checklist</h1>
      <p>{items.length} items.</p>
    </main>
  );
}
