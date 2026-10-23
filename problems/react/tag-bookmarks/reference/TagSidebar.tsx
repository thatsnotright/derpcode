import { useEffect, useRef, useState } from "react";
import type { FormEvent, KeyboardEvent } from "react";
import { normalizeTag } from "./tags";

type Props = {
  counts: [string, number][];
  selected: string[];
  onToggle: (tag: string) => void;
  onRename: (from: string, to: string) => void;
  onClear: () => void;
};

export default function TagSidebar({ counts, selected, onToggle, onRename, onClear }: Props) {
  const [renaming, setRenaming] = useState<string | null>(null);
  // DOM id to focus after the next render (the target may not exist yet).
  const focusAfter = useRef<string | null>(null);

  useEffect(() => {
    if (focusAfter.current) {
      document.getElementById(focusAfter.current)?.focus();
      focusAfter.current = null;
    }
  });

  const commit = (from: string, rawName: string) => {
    const to = normalizeTag(rawName);
    if (to && to !== from) {
      onRename(from, to);
      focusAfter.current = `tag-rename-${to}`;
    } else {
      focusAfter.current = `tag-rename-${from}`;
    }
    setRenaming(null);
  };

  const cancel = (from: string) => {
    focusAfter.current = `tag-rename-${from}`;
    setRenaming(null);
  };

  return (
    <aside className="sidebar" aria-labelledby="tags-heading">
      <h2 id="tags-heading">Tags</h2>
      {counts.length === 0 ? (
        <p className="empty">No tags yet.</p>
      ) : (
        <ul className="tag-list">
          {counts.map(([tag, count]) => (
            <li key={tag} className="tag-row">
              {renaming === tag ? (
                <RenameForm tag={tag} onCommit={commit} onCancel={cancel} />
              ) : (
                <>
                  <button
                    type="button"
                    id={`tag-toggle-${tag}`}
                    className={selected.includes(tag) ? "tag-toggle selected" : "tag-toggle"}
                    aria-pressed={selected.includes(tag)}
                    aria-label={`${tag}, ${count} ${count === 1 ? "bookmark" : "bookmarks"}`}
                    onClick={() => onToggle(tag)}
                  >
                    {selected.includes(tag) ? "[x] " : ""}
                    {tag} ({count})
                  </button>
                  <button
                    type="button"
                    id={`tag-rename-${tag}`}
                    aria-label={`Rename tag ${tag}`}
                    onClick={() => setRenaming(tag)}
                  >
                    Rename
                  </button>
                </>
              )}
            </li>
          ))}
        </ul>
      )}
      <button type="button" onClick={onClear} disabled={selected.length === 0}>
        Clear filters
      </button>
    </aside>
  );
}

type RenameProps = {
  tag: string;
  onCommit: (from: string, name: string) => void;
  onCancel: (from: string) => void;
};

function RenameForm({ tag, onCommit, onCancel }: RenameProps) {
  const [value, setValue] = useState(tag);
  const ref = useRef<HTMLInputElement>(null);

  useEffect(() => {
    ref.current?.focus();
    ref.current?.select();
  }, []);

  const submit = (e: FormEvent) => {
    e.preventDefault();
    onCommit(tag, value);
  };

  const onKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Escape") onCancel(tag);
  };

  return (
    <form className="rename-form" onSubmit={submit}>
      <input
        ref={ref}
        type="text"
        value={value}
        aria-label={`New name for tag ${tag}`}
        onChange={(e) => setValue(e.target.value)}
        onKeyDown={onKeyDown}
      />
      <button type="submit">Save</button>
      <button type="button" onClick={() => onCancel(tag)}>
        Cancel
      </button>
    </form>
  );
}
