import { useMemo, useReducer, useRef, useState } from "react";
import "./styles.css";
import AddItemForm from "./AddItemForm";
import ItemRow from "./ItemRow";
import { computeTotals, groupAndSort, listReducer } from "./listState";
import { SEED_ITEMS } from "./types";
import type { Item, SortKey } from "./types";

const currency = new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" });

export default function App() {
  const [items, dispatch] = useReducer(listReducer, SEED_ITEMS);
  const [sortKey, setSortKey] = useState<SortKey>("name");
  const [purchasedLast, setPurchasedLast] = useState(true);
  const [message, setMessage] = useState("");
  const listHeadingRef = useRef<HTMLHeadingElement>(null);

  const groups = useMemo(() => groupAndSort(items, sortKey, purchasedLast), [items, sortKey, purchasedLast]);
  const totals = useMemo(() => computeTotals(items), [items]);

  const add = (item: Item) => {
    dispatch({ type: "add", item });
    setMessage(`Added ${item.name} to ${item.aisle}.`);
  };

  const remove = (id: string) => {
    const item = items.find((i) => i.id === id);
    dispatch({ type: "delete", id });
    setMessage(`Deleted ${item?.name ?? "item"}.`);
    listHeadingRef.current?.focus();
  };

  const toggle = (id: string) => {
    const item = items.find((i) => i.id === id);
    dispatch({ type: "togglePurchased", id });
    if (item) setMessage(`${item.name} ${item.purchased ? "put back on the list" : "marked purchased"}.`);
  };

  const changeQuantity = (id: string, delta: number) => {
    const item = items.find((i) => i.id === id);
    dispatch({ type: "changeQuantity", id, delta });
    if (item) setMessage(`${item.name} quantity ${Math.max(1, item.quantity + delta)}.`);
  };

  const clearPurchased = () => {
    dispatch({ type: "clearPurchased" });
    setMessage(`Cleared ${totals.purchased} purchased ${totals.purchased === 1 ? "item" : "items"}.`);
    listHeadingRef.current?.focus();
  };

  return (
    <main className="app">
      <h1>Shopping list</h1>
      <AddItemForm onAdd={add} />

      <div className="controls">
        <label>
          Sort by{" "}
          <select value={sortKey} onChange={(e) => setSortKey(e.target.value as SortKey)}>
            <option value="name">Name</option>
            <option value="quantity">Quantity</option>
            <option value="price">Price</option>
          </select>
        </label>
        <label>
          <input type="checkbox" checked={purchasedLast} onChange={(e) => setPurchasedLast(e.target.checked)} />{" "}
          Purchased last
        </label>
        <button type="button" onClick={clearPurchased} disabled={totals.purchased === 0}>
          Clear purchased
        </button>
      </div>

      <h2 ref={listHeadingRef} tabIndex={-1} className="list-heading">
        Your list
      </h2>
      {groups.length === 0 && <p className="empty">Your list is empty.</p>}
      {groups.map((group) => (
        <section key={group.aisle} aria-labelledby={`aisle-${group.aisle}`} className="aisle">
          <h3 id={`aisle-${group.aisle}`}>
            {group.aisle} <span className="count">({group.items.length})</span>
          </h3>
          <ul className="items">
            {group.items.map((item) => (
              <ItemRow
                key={item.id}
                item={item}
                currency={currency}
                onToggle={toggle}
                onQuantity={changeQuantity}
                onDelete={remove}
              />
            ))}
          </ul>
        </section>
      ))}

      <section className="totals" aria-labelledby="totals-heading">
        <h2 id="totals-heading">Totals</h2>
        <dl>
          <dt>Still to buy</dt>
          <dd>
            {totals.remaining} items, {currency.format(totals.toBuy)}
          </dd>
          <dt>Purchased</dt>
          <dd>{totals.purchased} items</dd>
          <dt>Estimated total</dt>
          <dd>{currency.format(totals.estimated)}</dd>
        </dl>
      </section>

      <div role="status" aria-live="polite" className="sr-only">
        {message} {totals.remaining} items remaining.
      </div>
    </main>
  );
}
