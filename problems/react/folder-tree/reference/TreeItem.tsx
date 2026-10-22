import { useEffect, useRef, useState } from "react";
import type { FormEvent, KeyboardEvent } from "react";
import type { Nodes } from "./types";

type Props = {
  nodes: Nodes;
  id: string;
  level: number;
  posInSet: number;
  setSize: number;
  expanded: Record<string, boolean>;
  selectedId: string;
  renamingId: string | null;
  onSelect: (id: string) => void;
  onToggle: (id: string) => void;
  onCommitRename: (id: string, name: string) => void;
  onCancelRename: (id: string) => void;
};

export default function TreeItem(props: Props) {
  const { nodes, id, level, posInSet, setSize, expanded, selectedId, renamingId } = props;
  const node = nodes[id];
  const isFolder = node.kind === "folder";
  const isOpen = isFolder && Boolean(expanded[id]);
  const isSelected = selectedId === id;
  const isRenaming = renamingId === id;

  return (
    <li
      id={`treeitem-${id}`}
      role="treeitem"
      aria-labelledby={`treelabel-${id}`}
      aria-level={level}
      aria-posinset={posInSet}
      aria-setsize={setSize}
      aria-expanded={isFolder ? isOpen : undefined}
      aria-selected={isSelected}
      tabIndex={isSelected ? 0 : -1}
      className="tree-li"
    >
      <div
        className={isSelected ? "tree-row selected" : "tree-row"}
        onClick={() => props.onSelect(id)}
      >
        {isFolder ? (
          <span
            className="twisty"
            aria-hidden="true"
            onClick={(e) => {
              e.stopPropagation();
              props.onSelect(id);
              props.onToggle(id);
            }}
          >
            {isOpen ? "v" : ">"}
          </span>
        ) : (
          <span className="twisty" aria-hidden="true" />
        )}
        <span className="kind" aria-hidden="true">
          {isFolder ? "[dir]" : "[file]"}
        </span>
        {isRenaming ? (
          <RenameInput
            id={id}
            initial={node.name}
            onCommit={props.onCommitRename}
            onCancel={props.onCancelRename}
          />
        ) : (
          <span id={`treelabel-${id}`} className="tree-label">
            {node.name}
            <span className="sr-only">{isFolder ? " folder" : " file"}</span>
          </span>
        )}
      </div>
      {isOpen && node.children.length > 0 && (
        <ul role="group">
          {node.children.map((childId, i) => (
            <TreeItem
              key={childId}
              {...props}
              id={childId}
              level={level + 1}
              posInSet={i + 1}
              setSize={node.children.length}
            />
          ))}
        </ul>
      )}
    </li>
  );
}

type RenameProps = {
  id: string;
  initial: string;
  onCommit: (id: string, name: string) => void;
  onCancel: (id: string) => void;
};

function RenameInput({ id, initial, onCommit, onCancel }: RenameProps) {
  const [value, setValue] = useState(initial);
  const ref = useRef<HTMLInputElement>(null);

  useEffect(() => {
    ref.current?.focus();
    ref.current?.select();
  }, []);

  const submit = (e: FormEvent) => {
    e.preventDefault();
    const name = value.trim();
    if (name) onCommit(id, name);
  };

  const onKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Escape") {
      e.stopPropagation();
      onCancel(id);
    }
  };

  return (
    <form onSubmit={submit} onClick={(e) => e.stopPropagation()}>
      {/* id is on the input so aria-labelledby on the treeitem still resolves to something. */}
      <input
        id={`treelabel-${id}`}
        ref={ref}
        type="text"
        value={value}
        aria-label={`New name for ${initial}`}
        onChange={(e) => setValue(e.target.value)}
        onKeyDown={onKeyDown}
        onBlur={() => onCancel(id)}
      />
    </form>
  );
}
