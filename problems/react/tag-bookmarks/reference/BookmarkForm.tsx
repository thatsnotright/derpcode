import { useId, useRef, useState } from "react";
import type { FormEvent } from "react";
import TagInput from "./TagInput";
import { parseUrl } from "./tags";
import type { BookmarkDraft } from "./tags";

type Props = {
  onAdd: (draft: BookmarkDraft) => void;
};

type Errors = { title?: string; url?: string };

export default function BookmarkForm({ onAdd }: Props) {
  const uid = useId();
  const [title, setTitle] = useState("");
  const [url, setUrl] = useState("");
  const [tags, setTags] = useState<string[]>([]);
  const [errors, setErrors] = useState<Errors>({});
  const titleRef = useRef<HTMLInputElement>(null);
  const urlRef = useRef<HTMLInputElement>(null);

  const titleErrorId = `${uid}-title-error`;
  const urlErrorId = `${uid}-url-error`;

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const trimmedTitle = title.trim();
    const parsedUrl = parseUrl(url);
    const found: Errors = {};
    if (!trimmedTitle) found.title = "Title is required.";
    if (!parsedUrl) found.url = "Enter a valid web address like https://example.com.";
    setErrors(found);
    if (found.title) {
      titleRef.current?.focus();
      return;
    }
    if (found.url || !parsedUrl) {
      urlRef.current?.focus();
      return;
    }
    onAdd({ title: trimmedTitle, url: parsedUrl, tags });
    setTitle("");
    setUrl("");
    setTags([]);
    titleRef.current?.focus();
  };

  return (
    <form className="bookmark-form" onSubmit={handleSubmit} noValidate aria-label="Add a bookmark">
      <div className="field">
        <label htmlFor={`${uid}-title`}>Title</label>
        <input
          id={`${uid}-title`}
          ref={titleRef}
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          aria-required={true}
          aria-invalid={errors.title ? true : undefined}
          aria-describedby={errors.title ? titleErrorId : undefined}
        />
        {errors.title && (
          <p id={titleErrorId} className="field-error">
            Error: {errors.title}
          </p>
        )}
      </div>
      <div className="field">
        <label htmlFor={`${uid}-url`}>URL</label>
        <input
          id={`${uid}-url`}
          ref={urlRef}
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          aria-required={true}
          aria-invalid={errors.url ? true : undefined}
          aria-describedby={errors.url ? urlErrorId : undefined}
        />
        {errors.url && (
          <p id={urlErrorId} className="field-error">
            Error: {errors.url}
          </p>
        )}
      </div>
      <TagInput label="Tags" tags={tags} onChange={setTags} />
      <button type="submit">Add bookmark</button>
    </form>
  );
}
