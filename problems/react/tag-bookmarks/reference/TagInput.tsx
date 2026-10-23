import { useId, useRef, useState } from "react";
import type { ChangeEvent, KeyboardEvent } from "react";
import { normalizeTag } from "./tags";

type Props = {
  label: string;
  tags: string[];
  onChange: (tags: string[]) => void;
};

export default function TagInput({ label, tags, onChange }: Props) {
  const uid = useId();
  const inputId = `${uid}-input`;
  const hintId = `${uid}-hint`;
  const [draft, setDraft] = useState("");
  const [status, setStatus] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  const addTags = (raw: string[]) => {
    const fresh: string[] = [];
    for (const r of raw) {
      const tag = normalizeTag(r);
      if (tag && !tags.includes(tag) && !fresh.includes(tag)) fresh.push(tag);
    }
    if (fresh.length === 0) {
      if (raw.some((r) => r.trim())) setStatus("That tag is already added.");
      return;
    }
    onChange([...tags, ...fresh]);
    setStatus(`Added tag ${fresh.join(", ")}.`);
  };

  // Commas are handled here (not in onKeyDown) so pasted "a, b, c" works too.
  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    if (value.includes(",")) {
      const parts = value.split(",");
      const rest = parts.pop() ?? "";
      addTags(parts);
      setDraft(rest);
    } else {
      setDraft(value);
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && draft.trim()) {
      // Only swallow Enter when it adds a tag; an empty Enter submits the form.
      e.preventDefault();
      addTags([draft]);
      setDraft("");
    } else if (e.key === "Backspace" && draft === "" && tags.length > 0) {
      const last = tags[tags.length - 1];
      onChange(tags.slice(0, -1));
      setStatus(`Removed tag ${last}.`);
    }
  };

  const handleBlur = () => {
    if (draft.trim()) {
      addTags([draft]);
      setDraft("");
    }
  };

  const removeTag = (tag: string) => {
    onChange(tags.filter((t) => t !== tag));
    setStatus(`Removed tag ${tag}.`);
    inputRef.current?.focus();
  };

  return (
    <div className="tag-input">
      <label htmlFor={inputId}>{label}</label>
      <div className="tag-box">
        {tags.length > 0 && (
          <ul className="chips" aria-label="Tags added so far">
            {tags.map((t) => (
              <li key={t} className="chip">
                <span>{t}</span>
                <button type="button" aria-label={`Remove tag ${t}`} onClick={() => removeTag(t)}>
                  x
                </button>
              </li>
            ))}
          </ul>
        )}
        <input
          id={inputId}
          ref={inputRef}
          type="text"
          value={draft}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          onBlur={handleBlur}
          aria-describedby={hintId}
        />
      </div>
      <p id={hintId} className="hint">
        Press Enter or comma to add a tag. Press Backspace in the empty field to remove the last tag.
      </p>
      <div role="status" aria-live="polite" className="sr-only">
        {status}
      </div>
    </div>
  );
}
