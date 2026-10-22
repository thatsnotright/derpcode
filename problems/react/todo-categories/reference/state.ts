export type Todo = {
  id: string;
  text: string;
  category: string;
  done: boolean;
};

export type StatusFilter = "all" | "active" | "done";

export type Action =
  | { type: "add"; text: string; category: string }
  | { type: "toggle"; id: string }
  | { type: "edit"; id: string; text: string }
  | { type: "delete"; id: string }
  | { type: "clearDone" };

export function todosReducer(state: Todo[], action: Action): Todo[] {
  switch (action.type) {
    case "add":
      return [
        ...state,
        {
          id: crypto.randomUUID(),
          text: action.text,
          category: action.category,
          done: false,
        },
      ];
    case "toggle":
      return state.map((t) => (t.id === action.id ? { ...t, done: !t.done } : t));
    case "edit":
      return state.map((t) => (t.id === action.id ? { ...t, text: action.text } : t));
    case "delete":
      return state.filter((t) => t.id !== action.id);
    case "clearDone":
      return state.filter((t) => !t.done);
  }
}
