import { useState } from "react";
import "./styles.css";

type Todo = {
  id: string;
  text: string;
  category: string;
  done: boolean;
};

const CATEGORIES: string[] = ["Work", "Home", "Errands"];

const SEED_TODOS: Todo[] = [
  { id: "1", text: "Review pull requests", category: "Work", done: false },
  { id: "2", text: "Water the plants", category: "Home", done: true },
  { id: "3", text: "Pick up dry cleaning", category: "Errands", done: false },
];

export default function App() {
  const [todos] = useState<Todo[]>(SEED_TODOS);

  return (
    <main className="app">
      <h1>Todos</h1>
      <p>
        {todos.length} todos across {CATEGORIES.length} categories.
      </p>
    </main>
  );
}
