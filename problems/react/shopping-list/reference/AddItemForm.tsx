import { useRef, useState } from "react";
import type { FormEvent } from "react";
import type { Item } from "./types";
import { AISLES } from "./types";

type Props = {
  onAdd: (item: Item) => void;
};

export default function AddItemForm({ onAdd }: Props) {
  const [name, setName] = useState("");
  const [quantity, setQuantity] = useState("1");
  const [price, setPrice] = useState("");
  const [aisle, setAisle] = useState(AISLES[0]);
  const [error, setError] = useState("");
  const nameRef = useRef<HTMLInputElement>(null);
  const qtyRef = useRef<HTMLInputElement>(null);

  const submit = (e: FormEvent) => {
    e.preventDefault();
    const trimmed = name.trim();
    const qty = Number(quantity);
    const unitPrice = price === "" ? 0 : Number(price);
    if (!trimmed) {
      setError("Item name is required.");
      nameRef.current?.focus();
      return;
    }
    if (!Number.isInteger(qty) || qty < 1) {
      setError("Quantity must be a whole number of at least 1.");
      qtyRef.current?.focus();
      return;
    }
    if (!Number.isFinite(unitPrice) || unitPrice < 0) {
      setError("Price must be zero or more.");
      return;
    }
    setError("");
    onAdd({ id: crypto.randomUUID(), name: trimmed, quantity: qty, unitPrice, aisle, purchased: false });
    setName("");
    setQuantity("1");
    setPrice("");
    nameRef.current?.focus();
  };

  return (
    <form className="add-form" onSubmit={submit} noValidate aria-label="Add an item">
      <label>
        Item
        <input
          ref={nameRef}
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          aria-invalid={error.startsWith("Item") ? true : undefined}
          aria-describedby={error ? "form-error" : undefined}
        />
      </label>
      <label>
        Quantity
        <input
          ref={qtyRef}
          type="number"
          min={1}
          step={1}
          value={quantity}
          onChange={(e) => setQuantity(e.target.value)}
          aria-invalid={error.startsWith("Quantity") ? true : undefined}
          aria-describedby={error ? "form-error" : undefined}
        />
      </label>
      <label>
        Unit price
        <input type="number" min={0} step="0.01" value={price} onChange={(e) => setPrice(e.target.value)} />
      </label>
      <label>
        Aisle
        <select value={aisle} onChange={(e) => setAisle(e.target.value)}>
          {AISLES.map((a) => (
            <option key={a} value={a}>
              {a}
            </option>
          ))}
        </select>
      </label>
      <button type="submit">Add item</button>
      <p id="form-error" role="alert" className="error">
        {error}
      </p>
    </form>
  );
}
