import { useState } from "react";
import "./styles.css";

type Card = {
  id: string;
  title: string;
};

type Column = {
  id: string;
  title: string;
  cards: Card[];
};

const SEED_COLUMNS: Column[] = [
  {
    id: "todo",
    title: "To do",
    cards: [
      { id: "c1", title: "Write project brief" },
      { id: "c2", title: "Set up CI" },
    ],
  },
  { id: "doing", title: "In progress", cards: [{ id: "c3", title: "Design login page" }] },
  { id: "done", title: "Done", cards: [] },
];

export default function App() {
  const [columns] = useState<Column[]>(SEED_COLUMNS);

  return (
    <main className="app">
      <h1>Board</h1>
      <p>{columns.length} columns.</p>
    </main>
  );
}
