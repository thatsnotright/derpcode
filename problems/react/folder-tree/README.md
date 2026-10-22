# Folder Tree (ARIA Tree)

Build a file-explorer style tree of nested folders and files. The starter defines a normalized `Nodes` map (id -> node with a `children` id list) and seed data. This is the hardest of the React set: the tree pattern and keyboard model are where strong candidates separate themselves. Steps 1-3 are the core; 4-5 are stretch.

## Steps

1. **Render the tree.** Recursively render the seed data as nested lists with a folder/file distinction (text or icon, not colour alone). Folders show a collapse/expand toggle; start with all folders expanded.
2. **Collapse and expand.** Toggling a folder hides or shows its descendants. Keep expanded state separate from the node data.
3. **Add, rename, delete.** Add a file or a folder as a child of the selected folder (or of the selected file's parent). New nodes enter rename mode immediately. Rename is inline (Enter commits, Escape cancels, empty names rejected). Delete removes a node and its whole subtree. The root cannot be deleted or renamed.
4. **Move.** Move the selected node into another folder. The destination list must exclude the node itself and its descendants (no cycles) and its current parent. A "Move to" select plus a button is fine; drag and drop is optional.
5. **ARIA tree keyboard model.** Use `role="tree"`, `treeitem`, and `group`, with a roving tabindex (only one item has `tabindex=0`). Keys: Up/Down move through visible items, Right expands a closed folder or moves to its first child, Left collapses an open folder or moves to the parent, Home/End jump to first/last visible item, F2 renames, Delete deletes, Enter or Space selects.
6. **Stretch.** Sort folders before files, prevent duplicate sibling names, type-ahead (typing a letter jumps to the next item starting with it).

## Accessibility checklist

- [ ] The tree container has `role="tree"` and an accessible name (`aria-label` or `aria-labelledby`); items are `role="treeitem"`, nested lists `role="group"`.
- [ ] Folders expose `aria-expanded` (true/false); files do not. Items expose `aria-level`, and `aria-selected` marks the selected item.
- [ ] Roving tabindex: exactly one treeitem is in the tab order; Tab leaves the tree rather than stepping through every node.
- [ ] Full keyboard model works: Up/Down/Left/Right/Home/End/F2/Delete/Enter.
- [ ] Every action reachable by mouse (add, rename, delete, move) is reachable by keyboard.
- [ ] Toolbar buttons have clear accessible names that mention the target (e.g. "Delete src") or the toolbar tells you what is selected.
- [ ] Focus management: after add, focus lands in the rename input; after rename commit/cancel focus returns to the item; after delete focus moves to a neighbour or parent, never `<body>`; after a move the moved item keeps focus.
- [ ] An `aria-live="polite"` region announces add, rename, delete, move and expand/collapse results.
- [ ] Collapse/expand state is not conveyed by an icon alone (icon is `aria-hidden`, state is on the treeitem).
- [ ] Visible focus indicator on the focused treeitem and the selected item are distinguishable.

## Discussion topics

- State shape: normalized `Record<id, node>` with `children: id[]` vs. a recursive nested object. Cost of rename/move/delete in each, and why normalized wins for move.
- Cycle prevention on move and how to compute descendants.
- Derived state: the flattened list of visible nodes used for keyboard navigation (recomputed from nodes + expanded).
- Where to keep `selectedId`, `expanded`, and `renamingId` (UI state) versus the nodes (data). Reducer for nodes, plain state for UI.
- Focus management in React: refs and effects versus `document.getElementById`; why the roving tabindex needs the DOM node to exist before it can be focused.
- Keys and remounting: why moving a node remounts its row and how that affects focus.
- Performance with 10k nodes: windowing, `React.memo` on rows, stable callbacks.
- Testing the keyboard model with `userEvent.keyboard`.
