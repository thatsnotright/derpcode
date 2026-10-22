import { useEffect, useRef, useState } from "react";
import type { FormEvent, KeyboardEvent } from "react";
import type { Todo } from "./state";

type Props = {
  todo: Todo;
  onToggle: (id: string) => void;
  onSave: (id: string, text: string) => void;
  onDelete: (id: string) => void;
};

export default function TodoItem({ todo, onToggle, onSave, onDelete }: Props) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(todo.text);
  const inputRef = useRef<HTMLInputElement>(null);
  const editButtonRef = useRef<HTMLButtonElement>(null);
  const wasEditing = useRef(false);

  // Move focus into the input when editing starts, back to "Edit" when it ends.
  useEffect(() => {
    if (editing) {
      inputRef.current?.focus();
      inputRef.current?.select();
    } else if (wasEditing.current) {
      editButtonRef.current?.focus();
    }
    wasEditing.current = editing;
  }, [editing]);

  const startEdit = () => {
    setDraft(todo.text);
    setEditing(true);
  };

  const save = (e: FormEvent) => {
    e.preventDefault();
    const text = draft.trim();
    if (text && text !== todo.text) onSave(todo.id, text);
    setEditing(false);
  };

  const onKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Escape") setEditing(false);
  };

  return (
    <li className={todo.done ? "todo done" : "todo"}>
      <input
        type="checkbox"
        checked={todo.done}
        onChange={() => onToggle(todo.id)}
        aria-label={`Mark "${todo.text}" as ${todo.done ? "not done" : "done"}`}
      />
      {editing ? (
        <form className="edit-form" onSubmit={save}>
          <input
            ref={inputRef}
            type="text"
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onKeyDown={onKeyDown}
            aria-label={`Edit text of "${todo.text}"`}
          />
          <button type="submit">Save</button>
          <button type="button" onClick={() => setEditing(false)}>
            Cancel
          </button>
        </form>
      ) : (
        <>
          <span className="text">{todo.text}</span>
          <span className="badge">{todo.category}</span>
          <button ref={editButtonRef} type="button" onClick={startEdit} aria-label={`Edit ${todo.text}`}>
            Edit
          </button>
          <button type="button" onClick={() => onDelete(todo.id)} aria-label={`Delete ${todo.text}`}>
            Delete
          </button>
        </>
      )}
    </li>
  );
}
