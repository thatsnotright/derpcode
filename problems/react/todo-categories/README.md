# Todo List with Categories

Build a todo list where every todo belongs to a category. The starter renders a heading and defines a `Todo` type plus a little seed data. Everything else is up to the candidate. Treat it as a pairing session: ask the candidate to talk through their state shape before typing.

## Steps

1. **Render and add.** Render the seed todos as a list. Add a form (text input + category select + "Add" button) that appends a todo. Ignore blank input (trim it). Clear the input after adding.
2. **Complete, edit, delete.** Each todo has a checkbox to toggle completion, an "Edit" button that swaps the text for an input (Enter saves, Escape cancels), and a "Delete" button. Completed todos should look different (strike-through) without relying on colour alone.
3. **Filter and counts.** Add a category filter (All + each category) and a status filter (All / Active / Done). Show counts: total, active, done, and the count for the current filtered view. Derive these; do not store them.
4. **Empty states and bulk action.** Show a friendly message when the filtered list is empty. Add a "Clear completed" button that is disabled when there is nothing to clear.
5. **Stretch: custom categories.** Let the user create a new category. Per-category counts appear next to each filter option (e.g. "Work (2)"). Deleting the last todo in a category must not break the filter.
6. **Stretch: persistence.** Persist todos to `localStorage` with a safe read (handle malformed JSON).

## Accessibility checklist

- [ ] Every input and select has a programmatically associated label (visible `<label>` or `aria-label`).
- [ ] Todos are rendered in a semantic `<ul>`/`<li>`.
- [ ] Every button has an accessible name that identifies its todo (e.g. "Delete Water the plants"), not just "Delete".
- [ ] The whole app is keyboard operable: add via Enter in the form, toggle via Space, edit and delete via buttons.
- [ ] Inline edit moves focus into the input when it opens and returns focus to the "Edit" button when it closes.
- [ ] After deleting a todo, focus goes somewhere sensible (not lost to `<body>`).
- [ ] An `aria-live="polite"` region announces adds, deletes and filter result counts.
- [ ] Completed state is conveyed by the checkbox state, not only colour or strike-through.
- [ ] Focus indicators are visible on every interactive element.
- [ ] Filters are real form controls (select or radio group/fieldset with legend), not clickable divs.

## Discussion topics

- State shape: flat array of todos with a `category` string vs. normalized categories with ids. Where does `editingId` live (list vs. item)?
- `useState` vs. `useReducer` for the todo list; why a reducer helps once there are 5+ actions.
- Derived state: filtered list and counts computed during render (optionally `useMemo`), never duplicated in state.
- Keys: stable ids, not array index, and what breaks with index keys when items are deleted mid-edit.
- Component boundaries: `TodoForm`, `Filters`, `TodoList`, `TodoItem`; which ones own state and which are presentational.
- Performance: when would `React.memo` on `TodoItem` help, and what would need to be stable (callbacks)?
- Testing: what to assert with React Testing Library (roles and labels, not class names).
