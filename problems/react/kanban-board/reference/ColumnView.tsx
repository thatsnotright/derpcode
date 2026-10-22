import { useState } from "react";
import type { FormEvent } from "react";
import CardItem from "./CardItem";
import type { Column } from "./types";

type Props = {
  column: Column;
  index: number;
  columnCount: number;
  onAddCard: (columnId: string, title: string) => void;
  onEdit: (id: string, title: string) => void;
  onDelete: (id: string) => void;
  onMove: (id: string, dir: "left" | "right" | "up" | "down") => void;
};

export default function ColumnView({ column, index, columnCount, onAddCard, onEdit, onDelete, onMove }: Props) {
  const [title, setTitle] = useState("");
  const headingId = `col-${column.id}-heading`;

  const submit = (e: FormEvent) => {
    e.preventDefault();
    const trimmed = title.trim();
    if (!trimmed) return;
    onAddCard(column.id, trimmed);
    setTitle("");
  };

  return (
    <section className="column" aria-labelledby={headingId}>
      <h2 id={headingId}>
        {column.title} <span className="count">({column.cards.length})</span>
      </h2>
      {column.cards.length === 0 ? (
        <p className="empty">No cards yet.</p>
      ) : (
        <ul className="cards">
          {column.cards.map((card, i) => (
            <CardItem
              key={card.id}
              card={card}
              canLeft={index > 0}
              canRight={index < columnCount - 1}
              canUp={i > 0}
              canDown={i < column.cards.length - 1}
              onEdit={onEdit}
              onDelete={onDelete}
              onMove={onMove}
            />
          ))}
        </ul>
      )}
      <form onSubmit={submit} className="add-card">
        <input
          type="text"
          data-add-input={column.id}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          aria-label={`Add card to ${column.title}`}
          placeholder="New card"
        />
        <button type="submit">Add</button>
      </form>
    </section>
  );
}
