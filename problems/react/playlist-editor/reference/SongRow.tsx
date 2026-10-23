import type { KeyboardEvent } from "react";
import { formatDuration } from "./types";
import type { Song } from "./types";

type Props = {
  song: Song;
  index: number;
  count: number;
  hintId: string;
  onMove: (id: string, direction: -1 | 1, via: "key" | "button") => void;
  onRemove: (id: string) => void;
};

export default function SongRow({ song, index, count, hintId, onMove, onRemove }: Props) {
  const onKeyDown = (e: KeyboardEvent<HTMLLIElement>) => {
    if (e.altKey && e.key === "ArrowUp") {
      e.preventDefault();
      onMove(song.id, -1, "key");
    } else if (e.altKey && e.key === "ArrowDown") {
      e.preventDefault();
      onMove(song.id, 1, "key");
    } else if (e.key === "Delete" && e.target === e.currentTarget) {
      onRemove(song.id);
    }
  };

  return (
    <li
      id={`song-${song.id}`}
      className="song"
      tabIndex={0}
      aria-describedby={hintId}
      aria-keyshortcuts="Alt+ArrowUp Alt+ArrowDown Delete"
      onKeyDown={onKeyDown}
    >
      <span className="position" aria-hidden="true">
        {index + 1}.
      </span>
      <span className="text">
        <strong>{song.title}</strong>
        {song.artist && <span className="artist"> by {song.artist}</span>}
      </span>
      <span className="duration">{formatDuration(song.seconds)}</span>
      <button
        id={`song-${song.id}-up`}
        type="button"
        onClick={() => onMove(song.id, -1, "button")}
        disabled={index === 0}
        aria-label={`Move ${song.title} up`}
      >
        Up
      </button>
      <button
        id={`song-${song.id}-down`}
        type="button"
        onClick={() => onMove(song.id, 1, "button")}
        disabled={index === count - 1}
        aria-label={`Move ${song.title} down`}
      >
        Down
      </button>
      <button type="button" onClick={() => onRemove(song.id)} aria-label={`Remove ${song.title}`}>
        Remove
      </button>
    </li>
  );
}
