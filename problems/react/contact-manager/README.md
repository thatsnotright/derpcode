# Contact Manager

Build a small address book. The starter renders a heading and defines a `Contact` type plus a little seed data. Everything else is up to the candidate. The interesting parts are form validation that is actually accessible, and how state is split between the list and each row. Ask the candidate to describe their component tree before typing.

## Steps

1. **Render and add.** Render the seed contacts as a list (name, company, email, phone). Add a form with name, email, phone and company fields and an "Add contact" button. Clear the form after a successful add.
2. **Validate.** Name and email are required, email must look like an email, phone is optional but must look like a phone number if present. Show an error message under each invalid field on submit, and move focus to the first invalid field. Do not rely on native browser validation bubbles (use `noValidate`).
3. **Edit and delete.** Each contact has an "Edit" button that swaps the row for the same form pre-filled (Save / Cancel, same validation), and a "Delete" button that asks for confirmation (inline "Delete? Confirm / Cancel", not `window.confirm`) before removing the contact.
4. **Search and sort.** Add a search box that filters by name, email, phone or company (case-insensitive, trimmed), and a sort control for name or company. Show "Showing X of Y contacts" and an empty state when nothing matches.
5. **Favorites.** Add a toggle per contact to mark it as a favorite, and a "Favorites only" filter that combines with search.
6. **Stretch: polish.** Detect duplicate emails when adding or editing. Persist contacts to `localStorage` with a safe read. Sort contacts with an empty company last.

## Accessibility checklist

- [ ] Every input has a visible `<label>` associated with it (`htmlFor`/`id` or wrapping), not just a placeholder.
- [ ] Required fields are indicated in text (not only colour or a bare asterisk) and use `aria-required` or `required`.
- [ ] Invalid inputs have `aria-invalid="true"` and `aria-describedby` pointing at the id of their error message.
- [ ] On a failed submit, focus moves to the first invalid field.
- [ ] Error messages do not rely on colour alone.
- [ ] Contacts are rendered in a semantic `<ul>`/`<li>`.
- [ ] Every row button names its contact (e.g. "Edit Ada Lovelace", "Delete Ada Lovelace"); the favorite toggle uses `aria-pressed`.
- [ ] Delete confirmation is keyboard operable, moves focus to the safe choice (Cancel), and returns focus to the "Delete" button if cancelled.
- [ ] Opening the edit form moves focus into it; closing it returns focus to the "Edit" button.
- [ ] After deleting a contact, focus goes somewhere sensible (e.g. the list heading), not `<body>`.
- [ ] An `aria-live="polite"` region announces adds, edits, deletes and the number of search results.
- [ ] The search input has an accessible name and the sort control is a labelled `<select>`.

## Discussion topics

- Where validation lives: a pure `validate(draft)` function returning an errors object vs. validating on each keystroke. When to show errors (on submit, on blur, on change after the first error).
- Controlled inputs and the draft object. Why the edit form should copy the contact into local draft state rather than editing the list in place.
- Where `editing` and `confirmingDelete` state live (per row vs. a single `editingId` in the parent) and the trade-offs.
- Derived state: the filtered and sorted list is computed in render, never stored. Sort a copy rather than mutating state.
- Keys: stable ids, and what goes wrong with index keys when rows are deleted mid-edit.
- Why `window.confirm` is a poor accessibility and testing choice, and what the inline alternative needs to get right.
- Component boundaries: `ContactForm` shared between add and edit; which components are presentational.
- Testing with React Testing Library: `getByRole("textbox", { name: "Email" })`, `toHaveAccessibleDescription`.
