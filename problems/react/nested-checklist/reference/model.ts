export type ChecklistItem = {
  id: string;
  text: string;
  checked: boolean; // only meaningful for leaves; parents are derived
  depth: number;
};

export type CheckState = "checked" | "unchecked" | "mixed";

export type TreeNode = {
  item: ChecklistItem;
  state: CheckState;
  children: TreeNode[];
};

/** Index just past the last descendant of items[index]. */
export function subtreeEnd(items: ChecklistItem[], index: number): number {
  const depth = items[index].depth;
  let end = index + 1;
  while (end < items.length && items[end].depth > depth) end += 1;
  return end;
}

export function canIndent(items: ChecklistItem[], index: number): boolean {
  return index > 0 && items[index].depth <= items[index - 1].depth;
}

export function canOutdent(items: ChecklistItem[], index: number): boolean {
  return items[index].depth > 0;
}

/** Build a nested tree from the flat list and derive each node's check state from its leaves. */
export function buildTree(items: ChecklistItem[]): TreeNode[] {
  const roots: TreeNode[] = [];
  const stack: TreeNode[] = [];
  for (const item of items) {
    const node: TreeNode = { item, state: item.checked ? "checked" : "unchecked", children: [] };
    while (stack.length > item.depth) stack.pop();
    if (stack.length === 0) roots.push(node);
    else stack[stack.length - 1].children.push(node);
    stack.push(node);
  }
  const derive = (node: TreeNode): void => {
    if (node.children.length === 0) return;
    node.children.forEach(derive);
    const states = node.children.map((c) => c.state);
    if (states.every((s) => s === "checked")) node.state = "checked";
    else if (states.every((s) => s === "unchecked")) node.state = "unchecked";
    else node.state = "mixed";
  };
  roots.forEach(derive);
  return roots;
}

export function leafProgress(items: ChecklistItem[]): { done: number; total: number } {
  let done = 0;
  let total = 0;
  items.forEach((item, i) => {
    const isLeaf = i === items.length - 1 || items[i + 1].depth <= item.depth;
    if (!isLeaf) return;
    total += 1;
    if (item.checked) done += 1;
  });
  return { done, total };
}

export type Action =
  | { type: "add"; item: ChecklistItem; afterId: string | null; asChild: boolean }
  | { type: "setText"; id: string; text: string }
  | { type: "toggle"; id: string; checked: boolean }
  | { type: "delete"; id: string }
  | { type: "indent"; id: string }
  | { type: "outdent"; id: string };

function shiftSubtree(items: ChecklistItem[], index: number, delta: number): ChecklistItem[] {
  const end = subtreeEnd(items, index);
  return items.map((it, i) => (i >= index && i < end ? { ...it, depth: it.depth + delta } : it));
}

export function checklistReducer(items: ChecklistItem[], action: Action): ChecklistItem[] {
  switch (action.type) {
    case "add": {
      if (action.afterId === null) return [...items, action.item];
      const index = items.findIndex((i) => i.id === action.afterId);
      if (index === -1) return items;
      const parent = items[index];
      const insertAt = subtreeEnd(items, index);
      const depth = action.asChild ? parent.depth + 1 : parent.depth;
      const next = [...items];
      // Adding a child to a checked parent: the parent's state is derived, so
      // the new unchecked child naturally makes it partial/unchecked.
      next.splice(insertAt, 0, { ...action.item, depth });
      return next;
    }
    case "setText":
      return items.map((i) => (i.id === action.id ? { ...i, text: action.text } : i));
    case "toggle": {
      const index = items.findIndex((i) => i.id === action.id);
      if (index === -1) return items;
      const end = subtreeEnd(items, index);
      // Setting the whole range is correct for leaves (range of 1) and parents alike.
      return items.map((it, i) => (i >= index && i < end ? { ...it, checked: action.checked } : it));
    }
    case "delete": {
      const index = items.findIndex((i) => i.id === action.id);
      if (index === -1) return items;
      return [...items.slice(0, index), ...items.slice(subtreeEnd(items, index))];
    }
    case "indent": {
      const index = items.findIndex((i) => i.id === action.id);
      if (index === -1 || !canIndent(items, index)) return items;
      return shiftSubtree(items, index, 1);
    }
    case "outdent": {
      const index = items.findIndex((i) => i.id === action.id);
      if (index === -1 || !canOutdent(items, index)) return items;
      return shiftSubtree(items, index, -1);
    }
  }
}
