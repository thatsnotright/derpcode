import { useEffect, useRef, useState } from "react";
import "./styles.css";
import SongForm from "./SongForm";
import SongRow from "./SongRow";
import { formatDuration, moveItem } from "./types";
import type { Song, SongDraft } from "./types";

const SEED_SONGS: Song[] = [
  { id: "1", title: "Blue in Green", artist: "Miles Davis", seconds: 337 },
  { id: "2", title: "Take Five", artist: "Dave Brubeck", seconds: 324 },
  { id: "3", title: "So What", artist: "Miles Davis", seconds: 562 },
];

type HistoryEntry = {
  label: string;
  songs: Song[]; // the playlist as it was before the action
};

const HINT_ID = "playlist-hint";

export default function App() {
  const [songs, setSongs] = useState<Song[]>(SEED_SONGS);
  const [history, setHistory] = useState<HistoryEntry[]>([]);
  const [message, setMessage] = useState("");
  const listHeadingRef = useRef<HTMLHeadingElement>(null);
  // DOM id to focus after the next render (the target may not exist yet).
  const pendingFocus = useRef<string | null>(null);

  useEffect(() => {
    if (pendingFocus.current) {
      document.getElementById(pendingFocus.current)?.focus();
      pendingFocus.current = null;
    }
  });

  // Derived state: never stored.
  const totalSeconds = songs.reduce((sum, s) => sum + s.seconds, 0);

  // Every change goes through here so it can be undone.
  const commit = (next: Song[], label: string) => {
    setHistory([...history, { label, songs }]);
    setSongs(next);
  };

  const addSong = (draft: SongDraft) => {
    commit([...songs, { ...draft, id: crypto.randomUUID() }], `add "${draft.title}"`);
    setMessage(`Added "${draft.title}".`);
  };

  const removeSong = (id: string) => {
    const song = songs.find((s) => s.id === id);
    if (!song) return;
    commit(
      songs.filter((s) => s.id !== id),
      `remove "${song.title}"`,
    );
    setMessage(`Removed "${song.title}".`);
    listHeadingRef.current?.focus();
  };

  const moveSong = (id: string, direction: -1 | 1, via: "key" | "button") => {
    const from = songs.findIndex((s) => s.id === id);
    const to = from + direction;
    if (from === -1 || to < 0 || to >= songs.length) return;
    commit(moveItem(songs, from, to), `move "${songs[from].title}"`);
    setMessage(`${songs[from].title} moved to position ${to + 1} of ${songs.length}.`);
    // Keep focus on the same button, unless it is about to become disabled.
    const atEdge = direction === -1 ? to === 0 : to === songs.length - 1;
    pendingFocus.current =
      via === "button" && !atEdge ? `song-${id}-${direction === -1 ? "up" : "down"}` : `song-${id}`;
  };

  const undo = () => {
    const last = history[history.length - 1];
    if (!last) return;
    setSongs(last.songs);
    setHistory(history.slice(0, -1));
    setMessage(`Undid: ${last.label}.`);
    listHeadingRef.current?.focus();
  };

  const lastEntry = history[history.length - 1];

  return (
    <main className="app">
      <h1>Playlist</h1>

      <SongForm songs={songs} onAdd={addSong} />

      <div className="summary">
        <h2 ref={listHeadingRef} tabIndex={-1}>
          Songs
        </h2>
        <button type="button" onClick={undo} disabled={!lastEntry}>
          {lastEntry ? `Undo: ${lastEntry.label}` : "Undo"}
        </button>
      </div>
      <p className="counts">
        {songs.length} {songs.length === 1 ? "song" : "songs"}, total duration {formatDuration(totalSeconds)}.
      </p>
      <p id={HINT_ID} className="hint">
        Focus a song and press Alt+Up or Alt+Down to move it, or Delete to remove it.
      </p>

      {songs.length === 0 ? (
        <p className="empty">The playlist is empty. Add a song above.</p>
      ) : (
        <ol className="songs">
          {songs.map((s, i) => (
            <SongRow
              key={s.id}
              song={s}
              index={i}
              count={songs.length}
              hintId={HINT_ID}
              onMove={moveSong}
              onRemove={removeSong}
            />
          ))}
        </ol>
      )}

      <div role="status" aria-live="polite" className="sr-only">
        {message}
      </div>
    </main>
  );
}
