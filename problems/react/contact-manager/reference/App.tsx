import { useMemo, useRef, useState } from "react";
import "./styles.css";
import ContactForm from "./ContactForm";
import ContactItem from "./ContactItem";
import { EMPTY_DRAFT } from "./types";
import type { Contact, ContactDraft, SortKey } from "./types";

const SEED_CONTACTS: Contact[] = [
  { id: "1", name: "Ada Lovelace", email: "ada@example.com", phone: "555-0100", company: "Analytical Engines", favorite: true },
  { id: "2", name: "Grace Hopper", email: "grace@example.com", phone: "555-0101", company: "Navy", favorite: false },
  { id: "3", name: "Alan Turing", email: "alan@example.com", phone: "", company: "Bletchley Park", favorite: false },
];

export default function App() {
  const [contacts, setContacts] = useState<Contact[]>(SEED_CONTACTS);
  const [query, setQuery] = useState("");
  const [sortKey, setSortKey] = useState<SortKey>("name");
  const [favoritesOnly, setFavoritesOnly] = useState(false);
  const [message, setMessage] = useState("");
  const listHeadingRef = useRef<HTMLHeadingElement>(null);

  // Derived state: never stored.
  const visible = useMemo(() => {
    const q = query.trim().toLowerCase();
    const filtered = contacts.filter(
      (c) =>
        (!favoritesOnly || c.favorite) &&
        (!q || [c.name, c.email, c.phone, c.company].some((v) => v.toLowerCase().includes(q))),
    );
    // Sort a copy; never mutate state.
    return [...filtered].sort((a, b) => a[sortKey].localeCompare(b[sortKey]) || a.name.localeCompare(b.name));
  }, [contacts, query, sortKey, favoritesOnly]);

  const addContact = (draft: ContactDraft) => {
    setContacts([...contacts, { ...draft, id: crypto.randomUUID(), favorite: false }]);
    setMessage(`Added ${draft.name}.`);
  };

  const saveContact = (id: string, draft: ContactDraft) => {
    setContacts(contacts.map((c) => (c.id === id ? { ...c, ...draft } : c)));
    setMessage(`Saved ${draft.name}.`);
  };

  const deleteContact = (id: string) => {
    const contact = contacts.find((c) => c.id === id);
    setContacts(contacts.filter((c) => c.id !== id));
    setMessage(`Deleted ${contact?.name ?? "contact"}.`);
    listHeadingRef.current?.focus();
  };

  const toggleFavorite = (id: string) => {
    const contact = contacts.find((c) => c.id === id);
    setContacts(contacts.map((c) => (c.id === id ? { ...c, favorite: !c.favorite } : c)));
    if (contact) {
      setMessage(`${contact.name} ${contact.favorite ? "removed from" : "added to"} favorites.`);
    }
  };

  return (
    <main className="app">
      <h1>Contacts</h1>

      <ContactForm
        initial={EMPTY_DRAFT}
        legend="Add a contact"
        submitLabel="Add contact"
        clearOnSubmit
        onSubmit={addContact}
      />

      <section className="toolbar" aria-label="Search and sort">
        <label>
          Search{" "}
          <input type="search" value={query} onChange={(e) => setQuery(e.target.value)} />
        </label>
        <label>
          Sort by{" "}
          <select value={sortKey} onChange={(e) => setSortKey(e.target.value as SortKey)}>
            <option value="name">Name</option>
            <option value="company">Company</option>
          </select>
        </label>
        <label className="check">
          <input
            type="checkbox"
            checked={favoritesOnly}
            onChange={(e) => setFavoritesOnly(e.target.checked)}
          />
          Favorites only
        </label>
      </section>

      <h2 ref={listHeadingRef} tabIndex={-1}>
        Your contacts
      </h2>
      <p className="counts">
        Showing {visible.length} of {contacts.length} contacts.
      </p>
      {visible.length === 0 ? (
        <p className="empty">No contacts match. Add one or change the search.</p>
      ) : (
        <ul className="contacts">
          {visible.map((c) => (
            <ContactItem
              key={c.id}
              contact={c}
              onToggleFavorite={toggleFavorite}
              onSave={saveContact}
              onDelete={deleteContact}
            />
          ))}
        </ul>
      )}

      <div role="status" aria-live="polite" className="sr-only">
        {message}
        {` Showing ${visible.length} of ${contacts.length} contacts.`}
      </div>
    </main>
  );
}
