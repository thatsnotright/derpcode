import type { Bookmark } from "./tags";

type Props = {
  bookmark: Bookmark;
  onDelete: (id: string) => void;
};

export default function BookmarkItem({ bookmark, onDelete }: Props) {
  return (
    <li className="bookmark">
      <div className="info">
        <a href={bookmark.url} target="_blank" rel="noopener noreferrer">
          {bookmark.title}
          <span className="sr-only"> (opens in a new tab)</span>
        </a>
        <div className="url">{bookmark.url}</div>
        {bookmark.tags.length > 0 && (
          <ul className="chips" aria-label={`Tags for ${bookmark.title}`}>
            {bookmark.tags.map((t) => (
              <li key={t} className="chip">
                {t}
              </li>
            ))}
          </ul>
        )}
      </div>
      <button type="button" onClick={() => onDelete(bookmark.id)} aria-label={`Delete bookmark ${bookmark.title}`}>
        Delete
      </button>
    </li>
  );
}
