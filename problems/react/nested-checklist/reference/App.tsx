import { useEffect, useMemo, useReducer, useRef, useState } from "react";
import type { FormEvent } from "react";
import "./styles.css";
import ChecklistTree from "./ChecklistNode";
import { buildTree, canIndent, canOutdent, checklistReducer, leafProgress, subtreeEnd } from "./model";
import type { ChecklistItem } from "./model";

const SEED_ITEMS: ChecklistItem[] = [
  { id: "1", text: "Plan the trip", checked: false, depth: 0 },
  { id: "2", text: "Book flights", checked: true, depth: 1 },
  { id: "3", text: "Reserve hotel", checked: false, depth: 1 },
  { id: "4", text: "Pack", checked: false, depth: 0 },
];

export default function App() {
  const [items, dispatch] = useReducer(checklistReducer, SEED_ITEMS);
  const [draft, setDraft] = useState("");
  const [message, setMessage] = useState("");
  const pendingFocus = useRef<string | null>(null);
  const topInputRef = useRef<HTMLInputElement>(null);

  // Re-parenting remounts rows, so restore focus after render.
  useEffect(() => {
    if (pendingFocus.current) {
      document
        .querySelector<HTMLElement>(`[data-row-id="${pendingFocus.current}"] input[type="text"]`)
        ?.focus();
      pendingFocus.current = null;
    }
  });

  const tree = useMemo(() => buildTree(items), [items]);
  const progress = useMemo(() => leafProgress(items), [items]);

  const textOf = (id: string) => items.find((i) => i.id === id)?.text || "untitled item";

  const addTop = (e: FormEvent) => {
    e.preventDefault();
    const text = draft.trim();
    if (!text) return;
    dispatch({ type: "add", item: { id: crypto.randomUUID(), text, checked: false, depth: 0 }, afterId: null, asChild: false });
    setDraft("");
    setMessage(`Added ${text}.`);
    topInputRef.current?.focus();
  };

  const add = (afterId: string, asChild: boolean) => {
    const id = crypto.randomUUID();
    dispatch({ type: "add", item: { id, text: "", checked: false, depth: 0 }, afterId, asChild });
    pendingFocus.current = id;
    setMessage(`Added ${asChild ? "child" : "sibling"} ${asChild ? "to" : "after"} ${textOf(afterId)}.`);
  };

  const remove = (id: string) => {
    const index = items.findIndex((i) => i.id === id);
    const end = subtreeEnd(items, index);
    const neighbour = items[end] ?? items[index - 1];
    dispatch({ type: "delete", id });
    setMessage(`Deleted ${textOf(id)}${end - index > 1 ? " and its sub-items" : ""}.`);
    if (neighbour) pendingFocus.current = neighbour.id;
    else topInputRef.current?.focus();
  };

  const indent = (id: string) => {
    const index = items.findIndex((i) => i.id === id);
    if (index < 1) return;
    dispatch({ type: "indent", id });
    pendingFocus.current = id;
    setMessage(`${textOf(id)} is now a child of ${textOf(items[index - 1].id)}.`);
  };

  const outdent = (id: string) => {
    dispatch({ type: "outdent", id });
    pendingFocus.current = id;
    setMessage(`${textOf(id)} moved up one level.`);
  };

  return (
    <main className="app">
      <h1>Checklist</h1>
      <p>
        <progress value={progress.done} max={Math.max(progress.total, 1)} aria-label="Overall progress" />{" "}
        {progress.done} of {progress.total} tasks done
      </p>

      <form className="row" onSubmit={addTop} aria-label="Add a top-level item">
        <label>
          <span className="sr-only">New item</span>
          <input ref={topInputRef} type="text" value={draft} onChange={(e) => setDraft(e.target.value)} placeholder="New item" />
        </label>
        <button type="submit">Add</button>
      </form>

      {items.length === 0 ? (
        <p className="empty">Nothing here yet.</p>
      ) : (
        <ChecklistTree
          nodes={tree}
          canIndent={(id) => canIndent(items, items.findIndex((i) => i.id === id))}
          canOutdent={(id) => canOutdent(items, items.findIndex((i) => i.id === id))}
          onToggle={(id, checked) => dispatch({ type: "toggle", id, checked })}
          onText={(id, text) => dispatch({ type: "setText", id, text })}
          onAddSibling={(id) => add(id, false)}
          onAddChild={(id) => add(id, true)}
          onIndent={indent}
          onOutdent={outdent}
          onDelete={remove}
        />
      )}

      <p className="hint">In a text field: Enter adds a sibling, Alt+Right indents, Alt+Left outdents, Backspace on an empty item deletes it.</p>
      <div role="status" aria-live="polite" className="sr-only">
        {message}
      </div>
    </main>
  );
}
