import { useEffect, useRef, useState } from "react";
import type { FormEvent, KeyboardEvent } from "react";
import type { Card } from "./types";

type Props = {
  card: Card;
  canLeft: boolean;
  canRight: boolean;
  canUp: boolean;
  canDown: boolean;
  onEdit: (id: string, title: string) => void;
  onDelete: (id: string) => void;
  onMove: (id: string, dir: "left" | "right" | "up" | "down") => void;
};

export default function CardItem({ card, canLeft, canRight, canUp, canDown, onEdit, onDelete, onMove }: Props) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(card.title);
  const inputRef = useRef<HTMLInputElement>(null);
  const editRef = useRef<HTMLButtonElement>(null);
  const wasEditing = useRef(false);

  useEffect(() => {
    if (editing) {
      inputRef.current?.focus();
      inputRef.current?.select();
    } else if (wasEditing.current) {
      editRef.current?.focus();
    }
    wasEditing.current = editing;
  }, [editing]);

  const save = (e: FormEvent) => {
    e.preventDefault();
    const title = draft.trim();
    if (title) onEdit(card.id, title);
    setEditing(false);
  };

  const onKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Escape") setEditing(false);
  };

  const name = card.title;

  return (
    // tabIndex -1 lets us restore focus to the card itself after a move.
    <li className="card" data-card-id={card.id} tabIndex={-1}>
      {editing ? (
        <form onSubmit={save} className="card-edit">
          <input
            ref={inputRef}
            type="text"
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onKeyDown={onKeyDown}
            aria-label={`Title for card ${name}`}
          />
          <button type="submit">Save</button>
          <button type="button" onClick={() => setEditing(false)}>
            Cancel
          </button>
        </form>
      ) : (
        <>
          <p className="card-title">{name}</p>
          <div className="card-actions">
            <button type="button" onClick={() => onMove(card.id, "left")} disabled={!canLeft} aria-label={`Move ${name} left`}>
              &larr;
            </button>
            <button type="button" onClick={() => onMove(card.id, "right")} disabled={!canRight} aria-label={`Move ${name} right`}>
              &rarr;
            </button>
            <button type="button" onClick={() => onMove(card.id, "up")} disabled={!canUp} aria-label={`Move ${name} up`}>
              &uarr;
            </button>
            <button type="button" onClick={() => onMove(card.id, "down")} disabled={!canDown} aria-label={`Move ${name} down`}>
              &darr;
            </button>
            <button
              ref={editRef}
              type="button"
              onClick={() => {
                setDraft(card.title);
                setEditing(true);
              }}
              aria-label={`Edit ${name}`}
            >
              Edit
            </button>
            <button type="button" onClick={() => onDelete(card.id)} aria-label={`Delete ${name}`}>
              Delete
            </button>
          </div>
        </>
      )}
    </li>
  );
}
