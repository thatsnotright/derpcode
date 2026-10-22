import type { KeyboardEvent } from "react";
import type { ChecklistItem, TreeNode } from "./model";
import TriCheckbox from "./TriCheckbox";

type Handlers = {
  canIndent: (id: string) => boolean;
  canOutdent: (id: string) => boolean;
  onToggle: (id: string, checked: boolean) => void;
  onText: (id: string, text: string) => void;
  onAddSibling: (id: string) => void;
  onAddChild: (id: string) => void;
  onIndent: (id: string) => void;
  onOutdent: (id: string) => void;
  onDelete: (id: string) => void;
};

type Props = { nodes: TreeNode[] } & Handlers;

export default function ChecklistTree({ nodes, ...handlers }: Props) {
  return (
    <ul className="checklist">
      {nodes.map((node) => (
        <ChecklistNode key={node.item.id} node={node} {...handlers} />
      ))}
    </ul>
  );
}

function ChecklistNode({ node, ...h }: { node: TreeNode } & Handlers) {
  const item: ChecklistItem = node.item;
  const name = item.text || "untitled item";

  const onKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      h.onAddSibling(item.id);
    } else if (e.altKey && e.key === "ArrowRight") {
      e.preventDefault();
      if (h.canIndent(item.id)) h.onIndent(item.id);
    } else if (e.altKey && e.key === "ArrowLeft") {
      e.preventDefault();
      if (h.canOutdent(item.id)) h.onOutdent(item.id);
    } else if (e.key === "Backspace" && item.text === "" && node.children.length === 0) {
      e.preventDefault();
      h.onDelete(item.id);
    }
  };

  return (
    <li>
      <div className={node.state === "checked" ? "row done" : "row"} data-row-id={item.id}>
        <TriCheckbox
          state={node.state}
          label={`${name}${node.state === "mixed" ? " (partially complete)" : ""}`}
          onChange={(checked) => h.onToggle(item.id, checked)}
        />
        <input
          type="text"
          className="text"
          value={item.text}
          onChange={(e) => h.onText(item.id, e.target.value)}
          onKeyDown={onKeyDown}
          aria-label={`Text of ${name}`}
        />
        <button type="button" onClick={() => h.onAddSibling(item.id)} aria-label={`Add sibling after ${name}`}>
          + Sibling
        </button>
        <button type="button" onClick={() => h.onAddChild(item.id)} aria-label={`Add child to ${name}`}>
          + Child
        </button>
        <button type="button" onClick={() => h.onOutdent(item.id)} disabled={!h.canOutdent(item.id)} aria-label={`Outdent ${name}`}>
          &larr;
        </button>
        <button type="button" onClick={() => h.onIndent(item.id)} disabled={!h.canIndent(item.id)} aria-label={`Indent ${name}`}>
          &rarr;
        </button>
        <button type="button" onClick={() => h.onDelete(item.id)} aria-label={`Delete ${name}`}>
          Delete
        </button>
      </div>
      {node.children.length > 0 && <ChecklistTree nodes={node.children} {...h} />}
    </li>
  );
}
