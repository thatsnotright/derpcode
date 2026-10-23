export type Bookmark = {
  id: string;
  title: string;
  url: string;
  tags: string[];
};

export type BookmarkDraft = Omit<Bookmark, "id">;

export function normalizeTag(raw: string): string {
  return raw.trim().toLowerCase().replace(/\s+/g, "-");
}

// Returns a normalized http(s) URL, or null if the input is not a usable URL.
export function parseUrl(raw: string): string | null {
  const trimmed = raw.trim();
  if (!trimmed) return null;
  const withScheme = /^https?:\/\//i.test(trimmed) ? trimmed : `https://${trimmed}`;
  try {
    const url = new URL(withScheme);
    return url.hostname.includes(".") ? url.href : null;
  } catch {
    return null;
  }
}

// [tag, count] pairs sorted by tag name.
export function countTags(bookmarks: Bookmark[]): [string, number][] {
  const counts = new Map<string, number>();
  for (const b of bookmarks) {
    for (const tag of b.tags) {
      counts.set(tag, (counts.get(tag) ?? 0) + 1);
    }
  }
  return Array.from(counts.entries()).sort((a, b) => a[0].localeCompare(b[0]));
}

// AND semantics: a bookmark must have every selected tag.
export function filterByTags(bookmarks: Bookmark[], selected: string[]): Bookmark[] {
  return bookmarks.filter((b) => selected.every((tag) => b.tags.includes(tag)));
}

// Renames `from` to `to` in one tag list, removing any duplicate this creates (a merge).
export function renameInList(tags: string[], from: string, to: string): string[] {
  const mapped = tags.map((t) => (t === from ? to : t));
  return mapped.filter((t, i) => mapped.indexOf(t) === i);
}

export function renameTag(bookmarks: Bookmark[], from: string, to: string): Bookmark[] {
  return bookmarks.map((b) => (b.tags.includes(from) ? { ...b, tags: renameInList(b.tags, from, to) } : b));
}
