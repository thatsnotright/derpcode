export type Card = {
  id: string;
  title: string;
};

export type Column = {
  id: string;
  title: string;
  cards: Card[];
};

export const SEED_COLUMNS: Column[] = [
  {
    id: "todo",
    title: "To do",
    cards: [
      { id: "c1", title: "Write project brief" },
      { id: "c2", title: "Set up CI" },
    ],
  },
  { id: "doing", title: "In progress", cards: [{ id: "c3", title: "Design login page" }] },
  { id: "done", title: "Done", cards: [] },
];
