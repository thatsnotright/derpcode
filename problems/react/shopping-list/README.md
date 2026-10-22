# Grouped Shopping List

Build a shopping list where items have a quantity and an aisle, are grouped by aisle, can be marked as purchased, sorted, and totalled. The starter has an `Item` type, an `AISLES` list and seed data. This is the gentlest of the React set; it is a good warm-up and a good test of how cleanly the candidate separates stored state from derived state.

## Steps

1. **Render and add.** Render the items. Add a form with name, quantity (number, minimum 1), aisle (select) and optional unit price. Validate: non-blank name, integer quantity >= 1; show an inline error message otherwise. Merge with an existing item of the same name in the same aisle by adding the quantities (discuss whether to do this).
2. **Update.** Toggle purchased. Change quantity with increment and decrement controls (never below 1). Delete an item.
3. **Group by aisle.** Render a section per aisle (in the order of `AISLES`) with a heading and an item count. Hide aisles with no items.
4. **Sort.** A "Sort by" select (Name, Quantity, Price) applied within each aisle, plus a "Purchased last" toggle that sinks bought items to the bottom of their group.
5. **Totals.** Show items remaining, items purchased, estimated total (quantity x price over all items) and "still to buy" total. Format currency with `Intl.NumberFormat`. Add a "Clear purchased" button.
6. **Stretch.** Inline-edit item names, persist to `localStorage`, add a text filter, collapse aisles.

## Accessibility checklist

- [ ] Every input and select has a programmatic label; the quantity input has a minimum and its error is tied to it with `aria-describedby`.
- [ ] Aisle groups use headings (`<h2>`) and each group's items are in a `<ul>`.
- [ ] Purchased state uses a real checkbox with a name that includes the item, not a styled div.
- [ ] Increment/decrement buttons have names such as "Increase quantity of Milk" and the current quantity is exposed (e.g. `<output>` or text) so it is announced.
- [ ] Validation errors are announced (`role="alert"` or live region) and focus moves to the offending field.
- [ ] After add, focus returns to the name field; after delete, focus moves to a sensible target.
- [ ] A polite live region announces totals changes or actions (add, delete, clear).
- [ ] Purchased items are not conveyed by strike-through/colour alone.
- [ ] Visible focus on every control; all functionality works with the keyboard only.

## Discussion topics

- Stored versus derived state: grouping, sorting and all totals are computed during render. What would go wrong if totals were kept in state?
- Money: floating-point in dollars vs. integer cents. Where would rounding bite?
- Grouping implementation (reduce into a record vs. filter per aisle) and how to keep a stable aisle order.
- Controlled inputs for number fields: storing a string while editing, parsing on submit.
- `useMemo` for the grouped/sorted view: when is it worth it, and what are its dependencies?
- Keys when items move between groups after sorting.
- Component boundaries: `AddItemForm`, `AisleGroup`, `ItemRow`, `Totals`.
- Testing: add two items in different aisles and assert group headings and totals.
