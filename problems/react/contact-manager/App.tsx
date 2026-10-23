import { useState } from "react";
import "./styles.css";

type Contact = {
  id: string;
  name: string;
  email: string;
  phone: string;
  company: string;
  favorite: boolean;
};

const SEED_CONTACTS: Contact[] = [
  { id: "1", name: "Ada Lovelace", email: "ada@example.com", phone: "555-0100", company: "Analytical Engines", favorite: true },
  { id: "2", name: "Grace Hopper", email: "grace@example.com", phone: "555-0101", company: "Navy", favorite: false },
  { id: "3", name: "Alan Turing", email: "alan@example.com", phone: "", company: "Bletchley Park", favorite: false },
];

export default function App() {
  const [contacts] = useState<Contact[]>(SEED_CONTACTS);

  return (
    <main className="app">
      <h1>Contacts</h1>
      <p>{contacts.length} contacts.</p>
    </main>
  );
}
