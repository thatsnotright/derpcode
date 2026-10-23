import { useState } from "react";
import "./styles.css";

type Bookmark = {
  id: string;
  title: string;
  url: string;
  tags: string[]; // lowercase, no spaces
};

const SEED_BOOKMARKS: Bookmark[] = [
  { id: "1", title: "MDN Web Docs", url: "https://developer.mozilla.org/", tags: ["reference", "web"] },
  { id: "2", title: "React documentation", url: "https://react.dev/", tags: ["react", "reference"] },
  { id: "3", title: "WAI-ARIA Authoring Practices", url: "https://www.w3.org/WAI/ARIA/apg/", tags: ["a11y", "reference"] },
];

export default function App() {
  const [bookmarks] = useState<Bookmark[]>(SEED_BOOKMARKS);

  return (
    <main className="app">
      <h1>Bookmarks</h1>
      <p>{bookmarks.length} bookmarks.</p>
    </main>
  );
}
