# Nested Checklist

Build a checklist whose items can have sub-items, with a parent checkbox that reflects its children (checked, unchecked, or partially checked). The starter stores a flat list with a `depth` per item; candidates may keep that model or convert to a tree, and explaining the tradeoff is part of the exercise. The key insight to look for: a parent's checked state should be derived from its children, not stored.

## Steps

1. **Render and add.** Render the items as a nested list (indentation is not enough; use real nested lists). Add a top-level item via a form. Toggle a leaf item.
2. **Edit and delete.** Items have an editable text field. Delete removes the item and its entire subtree. Enter in a text field adds a sibling after the item (and its subtree).
3. **Add sibling and add child.** Per-item buttons add a sibling or a child below the item and focus the new item's text field.
4. **Indent and outdent.** "Indent" makes an item a child of the previous item (only valid when there is a previous item and the result is at most one level deeper). "Outdent" moves an item one level up. The item's subtree moves with it. Disable controls when the move is invalid.
5. **Parent state.** A parent is checked when all its children are checked, unchecked when none are, and shows the indeterminate (partial) state otherwise. Checking a parent checks all descendants; unchecking clears them. Show "x of y done" progress overall.
6. **Stretch.** Collapse/expand subtrees, Alt+Arrow keyboard shortcuts for indent/outdent, drag to reorder, "hide completed".

## Accessibility checklist

- [ ] Hierarchy is exposed through nested `<ul>`/`<li>`, not just padding.
- [ ] Checkboxes are native `<input type="checkbox">`; the partial state is the DOM `indeterminate` property (or `aria-checked="mixed"` on a custom control), not just a different icon.
- [ ] Every checkbox, text field, and button is labelled with the item it acts on (e.g. "Indent Book flights").
- [ ] Indent/outdent and add-child/sibling are reachable by keyboard; Tab is NOT hijacked for indentation.
- [ ] Focus management: after add, focus goes to the new item's text field; after indent/outdent, the item keeps focus even though it re-parents and remounts; after delete, focus goes to a neighbour.
- [ ] A polite live region announces the result of indent/outdent/delete and progress changes ("Book flights is now a child of Plan the trip").
- [ ] Disabled controls are actually `disabled` when a move is invalid.
- [ ] Completed state is not conveyed by strike-through/colour alone.
- [ ] Visible focus indicators everywhere.

## Discussion topics

- Flat list with depth vs. a recursive tree: which operations are easy or hard in each (indent, outdent, delete subtree, render)?
- What does outdenting an item that has following siblings mean? (In the flat model they become its children; discuss alternatives.)
- Derived state: parent status computed from leaves; why storing `checked` on parents leads to drift. Complexity of the naive per-item recompute vs. a single pass.
- Immutable updates of a subtree range with slice/spread or `map` over index ranges.
- Keys: id-based keys; why re-parenting remounts DOM nodes and loses focus unless you restore it.
- The `indeterminate` property cannot be set via JSX; use a ref and effect.
- Performance and testing: recompute cost, `React.memo`, assertions on roles and `toBePartiallyChecked`.
