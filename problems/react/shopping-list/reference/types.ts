export type Item = {
  id: string;
  name: string;
  quantity: number;
  unitPrice: number; // dollars, 0 if unknown
  aisle: string;
  purchased: boolean;
};

export type SortKey = "name" | "quantity" | "price";

export const AISLES: string[] = ["Produce", "Dairy", "Bakery", "Pantry"];

export const SEED_ITEMS: Item[] = [
  { id: "1", name: "Bananas", quantity: 6, unitPrice: 0.25, aisle: "Produce", purchased: false },
  { id: "2", name: "Milk", quantity: 1, unitPrice: 3.49, aisle: "Dairy", purchased: false },
  { id: "3", name: "Sourdough loaf", quantity: 1, unitPrice: 5, aisle: "Bakery", purchased: true },
  { id: "4", name: "Spaghetti", quantity: 2, unitPrice: 1.89, aisle: "Pantry", purchased: false },
];
