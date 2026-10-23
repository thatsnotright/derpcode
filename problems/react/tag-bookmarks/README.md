# Tagged Bookmarks

Build a bookmark manager where every bookmark has a title, a URL and a set of tags. The starter renders a heading and defines a `Bookmark` type plus a little seed data. Everything else is up to the candidate. This is the hardest of the set: a reusable tag input widget, multi-select filtering, and a cross-cutting rename that touches every bookmark. Ask the candidate to sketch the state shape and which values are derived before typing.

## Steps

1. **Render and add.** Render the bookmarks (title as a link, URL, tags). Add a form with title and URL fields. Validate that the title is non-empty and the URL is a valid http(s) URL, with an accessible error message on each.
2. **Tag input.** Add a tag input to the form: typing a tag and pressing Enter or comma adds it as a chip. Backspace in the empty input removes the last chip, and each chip has a remove button. Tags are normalized (trimmed, lowercase, spaces become hyphens), and duplicates and empty tags are ignored. Enter in an empty tag input still submits the form.
3. **Filter by tags.** Show a sidebar listing every tag with a count of bookmarks that use it. Clicking a tag toggles it as a filter. Multiple selected tags combine with AND (a bookmark must have all of them). Show "Showing X of Y bookmarks", an empty state, and a "Clear filters" button.
4. **Delete.** Each bookmark has a "Delete" button. When the last bookmark using a tag is deleted, that tag disappears from the sidebar and must also drop out of the active filter.
5. **Rename a tag.** Each tag in the sidebar has a "Rename" button that turns it into an inline input (Enter commits, Escape cancels). Renaming updates the tag on every bookmark and in the active filter. If the new name already exists the tags merge, and a bookmark must never end up with the same tag twice.
6. **Stretch.** Edit an existing bookmark's tags using the same tag input. Tag autocomplete from existing tags. OR/AND toggle for the filter. Persist to `localStorage` with a safe read.

## Accessibility checklist

- [ ] Every input (title, URL, tag input, rename input) has a visible `<label>` or `aria-label`.
- [ ] Invalid inputs have `aria-invalid="true"` and `aria-describedby` pointing at their error message, and focus moves to the first invalid field on a failed submit.
- [ ] The tag input explains its keys in text (referenced with `aria-describedby`): "Press Enter or comma to add, Backspace to remove the last tag".
- [ ] Added tag chips are in a list, and each remove button names its tag (e.g. "Remove tag react").
- [ ] Adding and removing tags is announced through a live region (chips alone are silent).
- [ ] After removing a chip with its button, focus returns to the tag input rather than being lost.
- [ ] Sidebar tag filters are real buttons using `aria-pressed`, whose names include the count (e.g. "react, 2 bookmarks"); selected state is not conveyed by colour alone.
- [ ] The "Rename" buttons name their tag; the rename input moves focus into it on open, and focus returns to a sensible control on commit or cancel.
- [ ] An `aria-live="polite"` region announces filter results, renames and deletes.
- [ ] External links make it clear they open in a new tab (or do not open a new tab).
- [ ] After deleting a bookmark, focus goes somewhere sensible, not `<body>`.
- [ ] Focus indicators are visible on every interactive element.

## Discussion topics

- State shape: tags stored as strings on each bookmark vs. a normalized `tags` table with ids. How each choice changes the cost and risk of "rename a tag everywhere" (and merging).
- Derived state: tag counts, the active filter and the visible list are computed in render. Why a stale selected tag (its last bookmark was deleted) is best fixed by deriving `activeTags` rather than syncing state in an effect.
- Reusable controlled component: the `TagInput` API (`tags` + `onChange`), what state belongs inside it (the draft text, announcements) and what does not.
- Keyboard handling: `Enter`/`,`/`Backspace` on `onKeyDown` vs. handling commas in `onChange` (paste support), and not breaking form submission on an empty Enter.
- Pure helpers (`normalizeTag`, `renameTag`, `filterByTags`) that can be unit tested without rendering anything.
- Focus management after the DOM changes (rename commit, chip removal), and why ids or refs are needed.
- Counts: counting across all bookmarks vs. within the current filter, and what each tells the user.
- Chips as a widget: why a plain list of buttons plus a text input is preferable to inventing a custom ARIA pattern.
- Scale: what changes with 10,000 bookmarks (memoizing counts, virtualization, indexing tags to bookmark ids).
