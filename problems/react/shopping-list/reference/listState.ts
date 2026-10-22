import type { Item, SortKey } from "./types";
import { AISLES } from "./types";

export type Action =
  | { type: "add"; item: Item }
  | { type: "togglePurchased"; id: string }
  | { type: "changeQuantity"; id: string; delta: number }
  | { type: "delete"; id: string }
  | { type: "clearPurchased" };

export function listReducer(items: Item[], action: Action): Item[] {
  switch (action.type) {
    case "add": {
      const existing = items.find(
        (i) =>
          i.aisle === action.item.aisle && i.name.toLowerCase() === action.item.name.toLowerCase() && !i.purchased,
      );
      if (existing) {
        return items.map((i) =>
          i.id === existing.id ? { ...i, quantity: i.quantity + action.item.quantity } : i,
        );
      }
      return [...items, action.item];
    }
    case "togglePurchased":
      return items.map((i) => (i.id === action.id ? { ...i, purchased: !i.purchased } : i));
    case "changeQuantity":
      return items.map((i) =>
        i.id === action.id ? { ...i, quantity: Math.max(1, i.quantity + action.delta) } : i,
      );
    case "delete":
      return items.filter((i) => i.id !== action.id);
    case "clearPurchased":
      return items.filter((i) => !i.purchased);
  }
}

export type AisleGroup = { aisle: string; items: Item[] };

export function groupAndSort(items: Item[], sortKey: SortKey, purchasedLast: boolean): AisleGroup[] {
  const compare = (a: Item, b: Item): number => {
    if (purchasedLast && a.purchased !== b.purchased) return a.purchased ? 1 : -1;
    switch (sortKey) {
      case "name":
        return a.name.localeCompare(b.name);
      case "quantity":
        return b.quantity - a.quantity || a.name.localeCompare(b.name);
      case "price":
        return b.unitPrice - a.unitPrice || a.name.localeCompare(b.name);
    }
  };
  return AISLES.map((aisle) => ({
    aisle,
    items: items.filter((i) => i.aisle === aisle).sort(compare),
  })).filter((g) => g.items.length > 0);
}

export type Totals = {
  remaining: number;
  purchased: number;
  estimated: number;
  toBuy: number;
};

export function computeTotals(items: Item[]): Totals {
  let remaining = 0;
  let purchased = 0;
  let estimatedCents = 0;
  let toBuyCents = 0;
  for (const item of items) {
    // Integer cents avoid float drift when summing.
    const cents = Math.round(item.unitPrice * 100) * item.quantity;
    estimatedCents += cents;
    if (item.purchased) {
      purchased += 1;
    } else {
      remaining += 1;
      toBuyCents += cents;
    }
  }
  return { remaining, purchased, estimated: estimatedCents / 100, toBuy: toBuyCents / 100 };
}
