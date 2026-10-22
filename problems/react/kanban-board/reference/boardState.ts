import type { Card, Column } from "./types";

export type Action =
  | { type: "addCard"; columnId: string; card: Card }
  | { type: "editCard"; cardId: string; title: string }
  | { type: "deleteCard"; cardId: string }
  | { type: "moveCard"; cardId: string; toColumnId: string; toIndex: number }
  | { type: "addColumn"; column: Column };

export function findCard(columns: Column[], cardId: string): { column: Column; index: number } | null {
  for (const column of columns) {
    const index = column.cards.findIndex((c) => c.id === cardId);
    if (index !== -1) return { column, index };
  }
  return null;
}

export function boardReducer(columns: Column[], action: Action): Column[] {
  switch (action.type) {
    case "addColumn":
      return [...columns, action.column];
    case "addCard":
      return columns.map((col) =>
        col.id === action.columnId ? { ...col, cards: [...col.cards, action.card] } : col,
      );
    case "editCard":
      return columns.map((col) => ({
        ...col,
        cards: col.cards.map((c) => (c.id === action.cardId ? { ...c, title: action.title } : c)),
      }));
    case "deleteCard":
      return columns.map((col) => ({ ...col, cards: col.cards.filter((c) => c.id !== action.cardId) }));
    case "moveCard": {
      const found = findCard(columns, action.cardId);
      if (!found) return columns;
      const card = found.column.cards[found.index];
      // Remove first, then insert, so toIndex refers to the destination list without the card.
      const removed = columns.map((col) =>
        col.id === found.column.id ? { ...col, cards: col.cards.filter((c) => c.id !== card.id) } : col,
      );
      return removed.map((col) => {
        if (col.id !== action.toColumnId) return col;
        const index = Math.max(0, Math.min(action.toIndex, col.cards.length));
        return { ...col, cards: [...col.cards.slice(0, index), card, ...col.cards.slice(index)] };
      });
    }
  }
}
