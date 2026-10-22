# Kanban Board

Build a small kanban board: columns containing cards that can be added, edited, moved between columns and reordered. Drag and drop is explicitly optional; the accessible, keyboard-operable path matters more and is what we want to see first. The starter defines `Card` and `Column` types and seed data.

## Steps

1. **Render the board.** Render each column with a heading, a card count and its cards. Lay columns out horizontally (scrolling is fine).
2. **Add and edit cards.** Each column has an "Add card" form (labelled input + button). Cards can be edited inline (title) and deleted. Reject blank titles.
3. **Move between columns.** Each card has "Move left" / "Move right" controls that send it to the adjacent column (appended at the bottom). Disable or hide a control that has no destination.
4. **Reorder within a column.** Add "Move up" / "Move down" controls. Moving a card must keep keyboard focus on that card, even if it changes columns and remounts.
5. **Add a column.** A form at the end of the board creates a new empty column. Moves left/right then work with it.
6. **Stretch.** Rename or delete columns (what happens to the cards?), WIP limits per column, HTML5 drag and drop that complements (does not replace) the buttons, persistence.

## Accessibility checklist

- [ ] Each column is a landmark-ish region (`<section aria-labelledby>`) with a heading; cards are in a `<ul>` with `<li>` items.
- [ ] Every control has a name that includes the card, e.g. "Move Set up CI right", not just "Move right".
- [ ] All card operations (edit, delete, move in four directions) are keyboard operable without drag and drop.
- [ ] Add-card inputs are labelled with the column name (e.g. "Add card to To do").
- [ ] After a move, focus stays on the moved card (or one of its controls); after delete, focus goes to a neighbouring card or the column's add input; after adding a card, focus returns to that column's input.
- [ ] An `aria-live="polite"` region announces results: "Moved Set up CI to In progress, position 2 of 3".
- [ ] Inline edit moves focus to the input, supports Enter/Escape, and returns focus to the Edit button.
- [ ] Buttons that cannot act are disabled (and do not silently drop focus), or are hidden consistently.
- [ ] Column counts and the card's current column are available to assistive tech, not just by position on screen.
- [ ] Visible focus styles on all controls; layout still works at 200% zoom.

## Discussion topics

- State shape: columns each holding `cards[]` vs. normalized `cards` map + `cardIds` per column. Tradeoffs for moving a card across columns.
- A single `moveCard(cardId, toColumnId, toIndex)` reducer action that powers buttons, keyboard, and drag and drop alike.
- Immutability: updating nested arrays without mutation; avoiding two-column inconsistent updates.
- Focus restoration after remounts (moving a card changes its parent, so React remounts it); refs and an effect vs. `autoFocus`.
- Component boundaries: `Board`, `Column`, `CardItem`, `AddForm`; what props crossing them look like (callbacks with ids, not indices).
- How to add drag and drop later without rewriting state, and how to keep it accessible.
- Performance: memoizing cards, large boards, keys.
- Testing: moving a card via role queries and asserting both columns.
