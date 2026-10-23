import { useState } from "react";
import "./styles.css";

type Song = {
  id: string;
  title: string;
  artist: string;
  seconds: number; // duration in seconds
};

const SEED_SONGS: Song[] = [
  { id: "1", title: "Blue in Green", artist: "Miles Davis", seconds: 337 },
  { id: "2", title: "Take Five", artist: "Dave Brubeck", seconds: 324 },
  { id: "3", title: "So What", artist: "Miles Davis", seconds: 562 },
];

export default function App() {
  const [songs] = useState<Song[]>(SEED_SONGS);

  return (
    <main className="app">
      <h1>Playlist</h1>
      <p>{songs.length} songs.</p>
    </main>
  );
}
