import { useMemo, useRef, useState } from "react";
import "./styles.css";
import BookmarkForm from "./BookmarkForm";
import BookmarkItem from "./BookmarkItem";
import TagSidebar from "./TagSidebar";
import { countTags, filterByTags, renameInList, renameTag } from "./tags";
import type { Bookmark, BookmarkDraft } from "./tags";

const SEED_BOOKMARKS: Bookmark[] = [
  { id: "1", title: "MDN Web Docs", url: "https://developer.mozilla.org/", tags: ["reference", "web"] },
  { id: "2", title: "React documentation", url: "https://react.dev/", tags: ["react", "reference"] },
  { id: "3", title: "WAI-ARIA Authoring Practices", url: "https://www.w3.org/WAI/ARIA/apg/", tags: ["a11y", "reference"] },
];

export default function App() {
  const [bookmarks, setBookmarks] = useState<Bookmark[]>(SEED_BOOKMARKS);
  const [selected, setSelected] = useState<string[]>([]);
  const [message, setMessage] = useState("");
  const listHeadingRef = useRef<HTMLHeadingElement>(null);

  // Derived state: never stored.
  const counts = useMemo(() => countTags(bookmarks), [bookmarks]);
  // A selected tag whose last bookmark was deleted silently drops out of the filter.
  const activeTags = selected.filter((t) => counts.some(([name]) => name === t));
  const visible = filterByTags(bookmarks, activeTags);

  const addBookmark = (draft: BookmarkDraft) => {
    setBookmarks([...bookmarks, { ...draft, id: crypto.randomUUID() }]);
    setMessage(`Added "${draft.title}".`);
  };

  const deleteBookmark = (id: string) => {
    const bookmark = bookmarks.find((b) => b.id === id);
    setBookmarks(bookmarks.filter((b) => b.id !== id));
    setMessage(`Deleted "${bookmark?.title ?? "bookmark"}".`);
    listHeadingRef.current?.focus();
  };

  const toggleTag = (tag: string) => {
    setSelected(activeTags.includes(tag) ? activeTags.filter((t) => t !== tag) : [...activeTags, tag]);
  };

  const rename = (from: string, to: string) => {
    const merged = counts.some(([name]) => name === to);
    setBookmarks(renameTag(bookmarks, from, to));
    setSelected(renameInList(activeTags, from, to));
    setMessage(`Renamed tag "${from}" to "${to}"${merged ? ", merged with the existing tag" : ""}.`);
  };

  return (
    <main className="app">
      <h1>Bookmarks</h1>

      <BookmarkForm onAdd={addBookmark} />

      <div className="layout">
        <TagSidebar
          counts={counts}
          selected={activeTags}
          onToggle={toggleTag}
          onRename={rename}
          onClear={() => setSelected([])}
        />

        <section aria-labelledby="bookmarks-heading">
          <h2 id="bookmarks-heading" ref={listHeadingRef} tabIndex={-1}>
            Your bookmarks
          </h2>
          <p className="counts">
            Showing {visible.length} of {bookmarks.length} bookmarks
            {activeTags.length > 0 && <>, tagged {activeTags.join(" AND ")}</>}.
          </p>
          {visible.length === 0 ? (
            <p className="empty">No bookmarks match. Add one or clear the filters.</p>
          ) : (
            <ul className="bookmarks">
              {visible.map((b) => (
                <BookmarkItem key={b.id} bookmark={b} onDelete={deleteBookmark} />
              ))}
            </ul>
          )}
        </section>
      </div>

      <div role="status" aria-live="polite" className="sr-only">
        {message}
        {` Showing ${visible.length} of ${bookmarks.length} bookmarks.`}
      </div>
    </main>
  );
}
