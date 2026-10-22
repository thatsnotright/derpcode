import { ROOT_ID } from "./types";
import type { NodeKind, Nodes } from "./types";

export type Action =
  | { type: "add"; id: string; parentId: string; kind: NodeKind; name: string }
  | { type: "rename"; id: string; name: string }
  | { type: "delete"; id: string }
  | { type: "move"; id: string; targetId: string };

export type VisibleItem = { id: string; parentId: string | null; level: number };

export function findParent(nodes: Nodes, id: string): string | null {
  for (const node of Object.values(nodes)) {
    if (node.children.includes(id)) return node.id;
  }
  return null;
}

export function collectSubtree(nodes: Nodes, id: string): string[] {
  const node = nodes[id];
  if (!node) return [];
  return [id, ...node.children.flatMap((c) => collectSubtree(nodes, c))];
}

/** Folders a node may be moved into: not itself, not a descendant, not its current parent. */
export function validMoveTargets(nodes: Nodes, id: string): string[] {
  if (id === ROOT_ID) return [];
  const blocked = new Set(collectSubtree(nodes, id));
  const parent = findParent(nodes, id);
  return Object.values(nodes)
    .filter((n) => n.kind === "folder" && !blocked.has(n.id) && n.id !== parent)
    .map((n) => n.id);
}

/** Flat, ordered list of nodes that are currently visible (ancestors all expanded). */
export function flattenVisible(
  nodes: Nodes,
  expanded: Record<string, boolean>,
  id: string = ROOT_ID,
  parentId: string | null = null,
  level = 1,
): VisibleItem[] {
  const node = nodes[id];
  const self: VisibleItem = { id, parentId, level };
  if (node.kind === "file" || !expanded[id]) return [self];
  return [self, ...node.children.flatMap((c) => flattenVisible(nodes, expanded, c, id, level + 1))];
}

export function treeReducer(nodes: Nodes, action: Action): Nodes {
  switch (action.type) {
    case "add": {
      const parent = nodes[action.parentId];
      if (!parent || parent.kind !== "folder") return nodes;
      return {
        ...nodes,
        [action.id]: { id: action.id, name: action.name, kind: action.kind, children: [] },
        [parent.id]: { ...parent, children: [...parent.children, action.id] },
      };
    }
    case "rename": {
      const node = nodes[action.id];
      if (!node || action.id === ROOT_ID) return nodes;
      return { ...nodes, [action.id]: { ...node, name: action.name } };
    }
    case "delete": {
      if (action.id === ROOT_ID) return nodes;
      const parentId = findParent(nodes, action.id);
      if (!parentId) return nodes;
      const removed = new Set(collectSubtree(nodes, action.id));
      const next: Nodes = {};
      for (const node of Object.values(nodes)) {
        if (removed.has(node.id)) continue;
        next[node.id] =
          node.id === parentId
            ? { ...node, children: node.children.filter((c) => c !== action.id) }
            : node;
      }
      return next;
    }
    case "move": {
      if (!validMoveTargets(nodes, action.id).includes(action.targetId)) return nodes;
      const parentId = findParent(nodes, action.id);
      if (!parentId) return nodes;
      const oldParent = nodes[parentId];
      const target = nodes[action.targetId];
      return {
        ...nodes,
        [parentId]: { ...oldParent, children: oldParent.children.filter((c) => c !== action.id) },
        [target.id]: { ...target, children: [...target.children, action.id] },
      };
    }
  }
}
