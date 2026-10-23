# Playlist Editor

Build an editable playlist. The starter renders a heading and defines a `Song` type plus a little seed data. Everything else is up to the candidate. The interesting parts are reordering with a good keyboard story (focus must follow the moved item), a non-blocking warning, and an undo history. Ask the candidate how they will model undo before they start on it.

## Steps

1. **Render and total.** Render the songs as an ordered list showing position, title, artist and duration formatted as `m:ss`. Show the total playlist duration (and song count), derived from the data. Show an empty state when there are no songs.
2. **Add and remove.** Add a form with title, artist and duration (`m:ss`, e.g. `3:45`). Title is required and duration must parse; show an accessible error otherwise. Each row has a "Remove" button.
3. **Reorder.** Each row has "Move up" and "Move down" buttons (disabled at the ends). Also support keyboard shortcuts on a focused row: `Alt+ArrowUp` and `Alt+ArrowDown`. After a move, focus must stay with the moved song.
4. **Duplicate warning.** If the song being added has the same title and artist as an existing one (case-insensitive), do not add it straight away. Show a warning with "Add anyway" and "Dismiss" actions. The warning must be announced to screen reader users.
5. **Undo.** Add an "Undo" button that reverts the last add, remove or move, and is disabled when there is nothing to undo. Its label should say what it will undo.
6. **Stretch.** Multi-level undo with a redo button, `Ctrl/Cmd+Z`, flag existing duplicates in the list, drag-and-drop reordering as an enhancement (buttons must remain).

## Accessibility checklist

- [ ] The playlist is a semantic `<ol>` with `<li>` rows; the order is conveyed by the list, not only by visual position.
- [ ] Every input has an associated `<label>`; the duration format is explained in text.
- [ ] Invalid inputs use `aria-invalid` and `aria-describedby` pointing at the error message, and focus moves to the first invalid field.
- [ ] Every row button names its song (e.g. "Move Take Five up", "Remove Take Five"). Disabled state at the list ends uses the `disabled` attribute.
- [ ] Keyboard shortcuts are discoverable (visible hint and/or `aria-keyshortcuts`) and do not clash with screen reader or browser shortcuts.
- [ ] After a move, focus stays on the moved item (or on the same button, falling back to the row if that button becomes disabled).
- [ ] After a remove or undo, focus goes somewhere sensible, not `<body>`.
- [ ] The duplicate warning is in a `role="alert"` (or live) region and is keyboard operable.
- [ ] An `aria-live="polite"` region announces moves ("Take Five moved to position 1 of 3"), adds, removes and undo.
- [ ] The Undo button's accessible name says what will be undone.
- [ ] Focus indicators are visible on every interactive element.

## Discussion topics

- Undo design: snapshot stack of previous `songs` arrays vs. a stack of inverse actions vs. a reducer with action log. Memory trade-offs and why immutable updates make snapshots cheap.
- Where history lives (alongside `songs` in one reducer so they cannot get out of sync) and what should or should not be undoable (the form draft?).
- Reordering helper: a pure `moveItem(list, from, to)` that returns a new array; testing it in isolation.
- Moving focus after the DOM changes: refs, `useEffect`, and why you cannot focus in the event handler (the target has not re-rendered yet). Stable keys keep DOM nodes (and focus) across reorders.
- Keyboard shortcuts: why `Alt+Arrow`, scoping handlers to the row rather than `document`, and `aria-keyshortcuts`.
- Warning vs. error: why duplicates should not block, and how to announce a non-blocking warning.
- Derived state: total duration computed in render; formatting vs. storing seconds as numbers.
- Drag-and-drop as an enhancement, and why it can never be the only way to reorder.
