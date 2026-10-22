import { useMemo, useReducer, useRef, useState } from "react";
import type { FormEvent } from "react";
import "./styles.css";
import TodoItem from "./TodoItem";
import { todosReducer } from "./state";
import type { StatusFilter, Todo } from "./state";

const DEFAULT_CATEGORIES: string[] = ["Work", "Home", "Errands"];

const SEED_TODOS: Todo[] = [
  { id: "1", text: "Review pull requests", category: "Work", done: false },
  { id: "2", text: "Water the plants", category: "Home", done: true },
  { id: "3", text: "Pick up dry cleaning", category: "Errands", done: false },
];

const STATUS_OPTIONS: { value: StatusFilter; label: string }[] = [
  { value: "all", label: "All" },
  { value: "active", label: "Active" },
  { value: "done", label: "Done" },
];

export default function App() {
  const [todos, dispatch] = useReducer(todosReducer, SEED_TODOS);
  const [categories, setCategories] = useState<string[]>(DEFAULT_CATEGORIES);
  const [categoryFilter, setCategoryFilter] = useState<string>("all");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");
  const [text, setText] = useState("");
  const [category, setCategory] = useState<string>(DEFAULT_CATEGORIES[0]);
  const [newCategory, setNewCategory] = useState("");
  const [message, setMessage] = useState("");

  const textRef = useRef<HTMLInputElement>(null);
  const listHeadingRef = useRef<HTMLHeadingElement>(null);

  // Derived state: never stored.
  const visible = useMemo(
    () =>
      todos.filter(
        (t) =>
          (categoryFilter === "all" || t.category === categoryFilter) &&
          (statusFilter === "all" || (statusFilter === "done") === t.done),
      ),
    [todos, categoryFilter, statusFilter],
  );
  const doneCount = todos.filter((t) => t.done).length;
  const countFor = (c: string) => todos.filter((t) => t.category === c).length;

  const addTodo = (e: FormEvent) => {
    e.preventDefault();
    const trimmed = text.trim();
    if (!trimmed) {
      setMessage("Enter some text before adding.");
      textRef.current?.focus();
      return;
    }
    dispatch({ type: "add", text: trimmed, category });
    setText("");
    setMessage(`Added "${trimmed}" to ${category}.`);
    textRef.current?.focus();
  };

  const addCategory = (e: FormEvent) => {
    e.preventDefault();
    const name = newCategory.trim();
    if (!name) return;
    if (categories.some((c) => c.toLowerCase() === name.toLowerCase())) {
      setMessage(`Category "${name}" already exists.`);
      return;
    }
    setCategories([...categories, name]);
    setCategory(name);
    setNewCategory("");
    setMessage(`Category "${name}" created.`);
  };

  const deleteTodo = (id: string) => {
    const todo = todos.find((t) => t.id === id);
    dispatch({ type: "delete", id });
    setMessage(`Deleted "${todo?.text ?? "todo"}".`);
    listHeadingRef.current?.focus();
  };

  const clearDone = () => {
    dispatch({ type: "clearDone" });
    setMessage(`Cleared ${doneCount} completed ${doneCount === 1 ? "todo" : "todos"}.`);
    listHeadingRef.current?.focus();
  };

  const toggleTodo = (id: string) => {
    const todo = todos.find((t) => t.id === id);
    dispatch({ type: "toggle", id });
    if (todo) setMessage(`"${todo.text}" marked ${todo.done ? "active" : "done"}.`);
  };

  // If the selected category filter no longer exists, fall back to "all".
  const activeCategoryFilter =
    categoryFilter === "all" || categories.includes(categoryFilter) ? categoryFilter : "all";

  return (
    <main className="app">
      <h1>Todos</h1>

      <form className="row" onSubmit={addTodo} aria-label="Add a todo">
        <label>
          <span className="sr-only">Todo text</span>
          <input
            ref={textRef}
            type="text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="What needs doing?"
          />
        </label>
        <label>
          <span className="sr-only">Category</span>
          <select value={category} onChange={(e) => setCategory(e.target.value)}>
            {categories.map((c) => (
              <option key={c} value={c}>
                {c}
              </option>
            ))}
          </select>
        </label>
        <button type="submit">Add</button>
      </form>

      <form className="row" onSubmit={addCategory} aria-label="Add a category">
        <label>
          <span className="sr-only">New category name</span>
          <input
            type="text"
            value={newCategory}
            onChange={(e) => setNewCategory(e.target.value)}
            placeholder="New category"
          />
        </label>
        <button type="submit">Add category</button>
      </form>

      <section className="filters" aria-label="Filters">
        <label>
          Category{" "}
          <select value={activeCategoryFilter} onChange={(e) => setCategoryFilter(e.target.value)}>
            <option value="all">All ({todos.length})</option>
            {categories.map((c) => (
              <option key={c} value={c}>
                {c} ({countFor(c)})
              </option>
            ))}
          </select>
        </label>
        <fieldset>
          <legend>Status</legend>
          {STATUS_OPTIONS.map((o) => (
            <label key={o.value} className="radio">
              <input
                type="radio"
                name="status"
                value={o.value}
                checked={statusFilter === o.value}
                onChange={() => setStatusFilter(o.value)}
              />
              {o.label}
            </label>
          ))}
        </fieldset>
      </section>

      <p className="counts">
        {todos.length} total, {todos.length - doneCount} active, {doneCount} done. Showing {visible.length}.
      </p>

      <h2 ref={listHeadingRef} tabIndex={-1}>
        Your todos
      </h2>
      {visible.length === 0 ? (
        <p className="empty">Nothing to show. Add a todo or change the filters.</p>
      ) : (
        <ul className="todos">
          {visible.map((t) => (
            <TodoItem
              key={t.id}
              todo={t}
              onToggle={toggleTodo}
              onSave={(id, newText) => dispatch({ type: "edit", id, text: newText })}
              onDelete={deleteTodo}
            />
          ))}
        </ul>
      )}

      <button type="button" onClick={clearDone} disabled={doneCount === 0}>
        Clear completed
      </button>

      <div role="status" aria-live="polite" className="sr-only">
        {message}
        {` Showing ${visible.length} of ${todos.length} todos.`}
      </div>
    </main>
  );
}
