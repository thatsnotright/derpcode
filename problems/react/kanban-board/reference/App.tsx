import { useEffect, useReducer, useRef, useState } from "react";
import type { FormEvent } from "react";
import "./styles.css";
import ColumnView from "./ColumnView";
import { boardReducer, findCard } from "./boardState";
import { SEED_COLUMNS } from "./types";

type Direction = "left" | "right" | "up" | "down";

export default function App() {
  const [columns, dispatch] = useReducer(boardReducer, SEED_COLUMNS);
  const [newColumn, setNewColumn] = useState("");
  const [message, setMessage] = useState("");
  const pendingFocus = useRef<string | null>(null);

  // Runs after every render: restores focus once the new DOM exists.
  useEffect(() => {
    if (pendingFocus.current) {
      document.querySelector<HTMLElement>(pendingFocus.current)?.focus();
      pendingFocus.current = null;
    }
  });

  const addCard = (columnId: string, title: string) => {
    dispatch({ type: "addCard", columnId, card: { id: crypto.randomUUID(), title } });
    const col = columns.find((c) => c.id === columnId);
    setMessage(`Added ${title} to ${col?.title ?? "column"}.`);
    pendingFocus.current = `[data-add-input="${columnId}"]`;
  };

  const deleteCard = (cardId: string) => {
    const found = findCard(columns, cardId);
    if (!found) return;
    const { column, index } = found;
    const neighbour = column.cards[index + 1] ?? column.cards[index - 1];
    dispatch({ type: "deleteCard", cardId });
    setMessage(`Deleted ${column.cards[index].title}.`);
    pendingFocus.current = neighbour ? `[data-card-id="${neighbour.id}"]` : `[data-add-input="${column.id}"]`;
  };

  const moveCard = (cardId: string, dir: Direction) => {
    const found = findCard(columns, cardId);
    if (!found) return;
    const { column, index } = found;
    const card = column.cards[index];
    const colIndex = columns.findIndex((c) => c.id === column.id);
    let toColumn = column;
    let toIndex = index;
    if (dir === "up") toIndex = index - 1;
    if (dir === "down") toIndex = index + 1;
    if (dir === "left" || dir === "right") {
      const dest = columns[colIndex + (dir === "left" ? -1 : 1)];
      if (!dest) return;
      toColumn = dest;
      toIndex = dest.cards.length;
    }
    if (toIndex < 0 || toIndex > (toColumn === column ? column.cards.length - 1 : toColumn.cards.length)) return;
    dispatch({ type: "moveCard", cardId, toColumnId: toColumn.id, toIndex });
    const size = toColumn === column ? column.cards.length : toColumn.cards.length + 1;
    setMessage(`Moved ${card.title} to ${toColumn.title}, position ${toIndex + 1} of ${size}.`);
    pendingFocus.current = `[data-card-id="${cardId}"]`;
  };

  const addColumn = (e: FormEvent) => {
    e.preventDefault();
    const title = newColumn.trim();
    if (!title) return;
    dispatch({ type: "addColumn", column: { id: crypto.randomUUID(), title, cards: [] } });
    setNewColumn("");
    setMessage(`Added column ${title}.`);
  };

  return (
    <main className="app wide">
      <h1>Board</h1>
      <div className="board">
        {columns.map((column, i) => (
          <ColumnView
            key={column.id}
            column={column}
            index={i}
            columnCount={columns.length}
            onAddCard={addCard}
            onEdit={(id, title) => dispatch({ type: "editCard", cardId: id, title })}
            onDelete={deleteCard}
            onMove={moveCard}
          />
        ))}
        <form onSubmit={addColumn} className="column new-column">
          <label>
            <span className="sr-only">New column name</span>
            <input
              type="text"
              value={newColumn}
              onChange={(e) => setNewColumn(e.target.value)}
              placeholder="New column"
            />
          </label>
          <button type="submit">Add column</button>
        </form>
      </div>
      <div role="status" aria-live="polite" className="sr-only">
        {message}
      </div>
    </main>
  );
}
