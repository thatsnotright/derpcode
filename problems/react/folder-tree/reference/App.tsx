import { useEffect, useReducer, useRef, useState } from "react";
import type { KeyboardEvent } from "react";
import "./styles.css";
import TreeItem from "./TreeItem";
import { collectSubtree, findParent, flattenVisible, treeReducer, validMoveTargets } from "./treeState";
import { ROOT_ID, SEED_NODES } from "./types";
import type { NodeKind } from "./types";

export default function App() {
  const [nodes, dispatch] = useReducer(treeReducer, SEED_NODES);
  const [expanded, setExpanded] = useState<Record<string, boolean>>({ root: true, src: true });
  const [selectedId, setSelectedId] = useState<string>(ROOT_ID);
  const [renamingId, setRenamingId] = useState<string | null>(null);
  const [moveTarget, setMoveTarget] = useState("");
  const [message, setMessage] = useState("");
  const pendingFocus = useRef<string | null>(null);

  // Focus the DOM node of the item we just navigated to (after React has rendered it).
  useEffect(() => {
    if (pendingFocus.current && renamingId === null) {
      document.getElementById(`treeitem-${pendingFocus.current}`)?.focus();
      pendingFocus.current = null;
    }
  });

  const visible = flattenVisible(nodes, expanded);
  const selected = nodes[selectedId];
  const targets = validMoveTargets(nodes, selectedId);
  const effectiveMoveTarget = targets.includes(moveTarget) ? moveTarget : (targets[0] ?? "");

  const focusItem = (id: string) => {
    setSelectedId(id);
    pendingFocus.current = id;
  };

  const toggle = (id: string, open?: boolean) => {
    const next = open ?? !expanded[id];
    setExpanded((e) => ({ ...e, [id]: next }));
    setMessage(`${nodes[id].name} ${next ? "expanded" : "collapsed"}.`);
  };

  const add = (kind: NodeKind) => {
    const parentId = selected.kind === "folder" ? selected.id : (findParent(nodes, selected.id) ?? ROOT_ID);
    const id = crypto.randomUUID();
    const name = kind === "folder" ? "new-folder" : "untitled.txt";
    dispatch({ type: "add", id, parentId, kind, name });
    setExpanded((e) => ({ ...e, [parentId]: true }));
    setSelectedId(id);
    setRenamingId(id);
    setMessage(`Added ${kind} in ${nodes[parentId].name}. Type a name and press Enter.`);
  };

  const startRename = () => {
    if (selectedId === ROOT_ID) {
      setMessage("The root folder cannot be renamed.");
      return;
    }
    setRenamingId(selectedId);
  };

  const commitRename = (id: string, name: string) => {
    dispatch({ type: "rename", id, name });
    setRenamingId(null);
    focusItem(id);
    setMessage(`Renamed to ${name}.`);
  };

  const cancelRename = (id: string) => {
    if (renamingId !== id) return;
    setRenamingId(null);
    focusItem(id);
  };

  const remove = () => {
    if (selectedId === ROOT_ID) {
      setMessage("The root folder cannot be deleted.");
      return;
    }
    const index = visible.findIndex((v) => v.id === selectedId);
    const removedIds = new Set(collectSubtree(nodes, selectedId));
    const neighbour =
      visible.slice(index + 1).find((v) => !removedIds.has(v.id)) ?? visible[index - 1] ?? visible[0];
    const name = selected.name;
    dispatch({ type: "delete", id: selectedId });
    focusItem(neighbour.id);
    setMessage(`Deleted ${name}.`);
  };

  const move = () => {
    if (!effectiveMoveTarget) return;
    dispatch({ type: "move", id: selectedId, targetId: effectiveMoveTarget });
    setExpanded((e) => ({ ...e, [effectiveMoveTarget]: true }));
    focusItem(selectedId);
    setMessage(`Moved ${selected.name} into ${nodes[effectiveMoveTarget].name}.`);
  };

  const onKeyDown = (e: KeyboardEvent<HTMLUListElement>) => {
    if ((e.target as HTMLElement).tagName === "INPUT") return;
    const index = visible.findIndex((v) => v.id === selectedId);
    const current = visible[index];
    if (!current) return;
    const isFolder = selected.kind === "folder";
    const isOpen = isFolder && Boolean(expanded[selectedId]);

    switch (e.key) {
      case "ArrowDown":
        if (visible[index + 1]) focusItem(visible[index + 1].id);
        break;
      case "ArrowUp":
        if (visible[index - 1]) focusItem(visible[index - 1].id);
        break;
      case "ArrowRight":
        if (isFolder && !isOpen) toggle(selectedId, true);
        else if (isOpen && nodes[selectedId].children.length > 0) focusItem(nodes[selectedId].children[0]);
        break;
      case "ArrowLeft":
        if (isOpen) toggle(selectedId, false);
        else if (current.parentId) focusItem(current.parentId);
        break;
      case "Home":
        focusItem(visible[0].id);
        break;
      case "End":
        focusItem(visible[visible.length - 1].id);
        break;
      case "F2":
        startRename();
        break;
      case "Delete":
        remove();
        break;
      case "Enter":
      case " ":
        if (isFolder) toggle(selectedId);
        break;
      default:
        return;
    }
    e.preventDefault();
  };

  const select = (id: string) => {
    setSelectedId(id);
    pendingFocus.current = id;
  };

  const isRoot = selectedId === ROOT_ID;

  return (
    <main className="app">
      <h1 id="tree-title">File explorer</h1>
      <p>
        Selected: <strong>{selected.name}</strong>
      </p>

      <div role="toolbar" aria-label="Actions for selected item" className="toolbar">
        <button type="button" onClick={() => add("file")}>
          Add file
        </button>
        <button type="button" onClick={() => add("folder")}>
          Add folder
        </button>
        <button type="button" onClick={startRename} disabled={isRoot}>
          Rename
        </button>
        <button type="button" onClick={remove} disabled={isRoot}>
          Delete
        </button>
        <label>
          Move to{" "}
          <select
            value={effectiveMoveTarget}
            onChange={(e) => setMoveTarget(e.target.value)}
            disabled={targets.length === 0}
          >
            {targets.map((t) => (
              <option key={t} value={t}>
                {nodes[t].name}
              </option>
            ))}
          </select>
        </label>
        <button type="button" onClick={move} disabled={targets.length === 0}>
          Move
        </button>
      </div>

      <ul role="tree" aria-labelledby="tree-title" className="tree" onKeyDown={onKeyDown}>
        <TreeItem
          nodes={nodes}
          id={ROOT_ID}
          level={1}
          posInSet={1}
          setSize={1}
          expanded={expanded}
          selectedId={selectedId}
          renamingId={renamingId}
          onSelect={select}
          onToggle={(id) => toggle(id)}
          onCommitRename={commitRename}
          onCancelRename={cancelRename}
        />
      </ul>

      <p className="hint">
        Arrow keys navigate, Right/Left expand and collapse, F2 renames, Delete removes, Home/End jump.
      </p>
      <div role="status" aria-live="polite" className="sr-only">
        {message}
      </div>
    </main>
  );
}
