import { useId, useRef, useState } from "react";
import type { ChangeEvent, FormEvent } from "react";
import { parseDuration, sameSong } from "./types";
import type { Song, SongDraft } from "./types";

type Props = {
  songs: Song[];
  onAdd: (draft: SongDraft) => void;
};

type FieldError = { field: "title" | "duration"; message: string } | null;

export default function SongForm({ songs, onAdd }: Props) {
  const uid = useId();
  const [title, setTitle] = useState("");
  const [artist, setArtist] = useState("");
  const [duration, setDuration] = useState("");
  const [error, setError] = useState<FieldError>(null);
  const [warning, setWarning] = useState(false);
  const titleRef = useRef<HTMLInputElement>(null);
  const durationRef = useRef<HTMLInputElement>(null);

  const titleErrorId = `${uid}-title-error`;
  const durationErrorId = `${uid}-duration-error`;
  const durationHintId = `${uid}-duration-hint`;

  // Validates the fields. Returns the draft, or null after showing an error.
  const readDraft = (): SongDraft | null => {
    const trimmedTitle = title.trim();
    if (!trimmedTitle) {
      setError({ field: "title", message: "Title is required." });
      titleRef.current?.focus();
      return null;
    }
    const seconds = parseDuration(duration);
    if (seconds === null) {
      setError({ field: "duration", message: "Enter a duration like 3:45." });
      durationRef.current?.focus();
      return null;
    }
    setError(null);
    return { title: trimmedTitle, artist: artist.trim(), seconds };
  };

  const commit = (draft: SongDraft) => {
    onAdd(draft);
    setTitle("");
    setArtist("");
    setDuration("");
    setWarning(false);
    titleRef.current?.focus();
  };

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const draft = readDraft();
    if (!draft) return;
    if (!warning && songs.some((s) => sameSong(s, draft))) {
      setWarning(true);
      return;
    }
    commit(draft);
  };

  const addAnyway = () => {
    const draft = readDraft();
    if (draft) commit(draft);
  };

  const dismissWarning = () => {
    setWarning(false);
    titleRef.current?.focus();
  };

  const edit = (setter: (value: string) => void) => (e: ChangeEvent<HTMLInputElement>) => {
    setter(e.target.value);
    setWarning(false);
  };

  return (
    <form className="song-form" onSubmit={handleSubmit} noValidate aria-label="Add a song">
      <div className="field">
        <label htmlFor={`${uid}-title`}>Title</label>
        <input
          id={`${uid}-title`}
          ref={titleRef}
          type="text"
          value={title}
          onChange={edit(setTitle)}
          aria-required={true}
          aria-invalid={error?.field === "title" ? true : undefined}
          aria-describedby={error?.field === "title" ? titleErrorId : undefined}
        />
        {error?.field === "title" && (
          <p id={titleErrorId} className="field-error">
            Error: {error.message}
          </p>
        )}
      </div>
      <div className="field">
        <label htmlFor={`${uid}-artist`}>Artist</label>
        <input id={`${uid}-artist`} type="text" value={artist} onChange={edit(setArtist)} />
      </div>
      <div className="field">
        <label htmlFor={`${uid}-duration`}>Duration</label>
        <input
          id={`${uid}-duration`}
          ref={durationRef}
          type="text"
          inputMode="numeric"
          placeholder="3:45"
          value={duration}
          onChange={edit(setDuration)}
          aria-required={true}
          aria-invalid={error?.field === "duration" ? true : undefined}
          aria-describedby={error?.field === "duration" ? `${durationHintId} ${durationErrorId}` : durationHintId}
        />
        <span id={durationHintId} className="hint">
          Minutes and seconds, like 3:45.
        </span>
        {error?.field === "duration" && (
          <p id={durationErrorId} className="field-error">
            Error: {error.message}
          </p>
        )}
      </div>
      <button type="submit">Add song</button>

      {warning && (
        <div role="alert" className="warning">
          <span>
            Warning: &quot;{title.trim()}&quot; by {artist.trim() || "an unknown artist"} is already in the playlist.
          </span>
          <button type="button" onClick={addAnyway}>
            Add anyway
          </button>
          <button type="button" onClick={dismissWarning}>
            Dismiss
          </button>
        </div>
      )}
    </form>
  );
}
