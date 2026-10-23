import { useEffect, useRef, useState } from "react";
import ContactForm from "./ContactForm";
import type { Contact, ContactDraft } from "./types";

type Props = {
  contact: Contact;
  onToggleFavorite: (id: string) => void;
  onSave: (id: string, draft: ContactDraft) => void;
  onDelete: (id: string) => void;
};

type Mode = "view" | "edit" | "confirm";

export default function ContactItem({ contact, onToggleFavorite, onSave, onDelete }: Props) {
  const [mode, setMode] = useState<Mode>("view");
  const editRef = useRef<HTMLButtonElement>(null);
  const deleteRef = useRef<HTMLButtonElement>(null);
  const cancelDeleteRef = useRef<HTMLButtonElement>(null);
  const previousMode = useRef<Mode>("view");

  // Focus management: into the confirm prompt (on the safe choice), and back to
  // the originating button when edit or confirm ends without deleting.
  useEffect(() => {
    if (mode === "confirm") {
      cancelDeleteRef.current?.focus();
    } else if (mode === "view" && previousMode.current === "edit") {
      editRef.current?.focus();
    } else if (mode === "view" && previousMode.current === "confirm") {
      deleteRef.current?.focus();
    }
    previousMode.current = mode;
  }, [mode]);

  if (mode === "edit") {
    return (
      <li className="contact">
        <ContactForm
          initial={{ name: contact.name, email: contact.email, phone: contact.phone, company: contact.company }}
          legend={`Edit ${contact.name}`}
          submitLabel="Save"
          autoFocus
          onSubmit={(draft) => {
            onSave(contact.id, draft);
            setMode("view");
          }}
          onCancel={() => setMode("view")}
        />
      </li>
    );
  }

  return (
    <li className="contact">
      <div className="contact-row">
        <div className="info">
          <strong>{contact.name}</strong>
          {contact.favorite && <span className="badge">Favorite</span>}
          <div className="meta">
            {contact.company && <span>{contact.company}</span>}
            <a href={`mailto:${contact.email}`}>{contact.email}</a>
            {contact.phone && <span>{contact.phone}</span>}
          </div>
        </div>
        <div className="actions">
          <button
            type="button"
            aria-pressed={contact.favorite}
            aria-label={`Favorite ${contact.name}`}
            onClick={() => onToggleFavorite(contact.id)}
          >
            {contact.favorite ? "Favorited" : "Favorite"}
          </button>
          <button ref={editRef} type="button" aria-label={`Edit ${contact.name}`} onClick={() => setMode("edit")}>
            Edit
          </button>
          <button ref={deleteRef} type="button" aria-label={`Delete ${contact.name}`} onClick={() => setMode("confirm")}>
            Delete
          </button>
        </div>
      </div>
      {mode === "confirm" && (
        <div className="confirm" role="group" aria-label={`Confirm deleting ${contact.name}`}>
          <span>Delete {contact.name}? This cannot be undone.</span>
          <button ref={cancelDeleteRef} type="button" onClick={() => setMode("view")}>
            Cancel
          </button>
          <button type="button" className="danger" onClick={() => onDelete(contact.id)}>
            Confirm delete
          </button>
        </div>
      )}
    </li>
  );
}
