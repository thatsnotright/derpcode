import type { Item } from "./types";

type Props = {
  item: Item;
  currency: Intl.NumberFormat;
  onToggle: (id: string) => void;
  onQuantity: (id: string, delta: number) => void;
  onDelete: (id: string) => void;
};

export default function ItemRow({ item, currency, onToggle, onQuantity, onDelete }: Props) {
  const lineTotal = item.quantity * item.unitPrice;
  return (
    <li className={item.purchased ? "item purchased" : "item"}>
      <input
        type="checkbox"
        checked={item.purchased}
        onChange={() => onToggle(item.id)}
        aria-label={`Purchased: ${item.name}`}
      />
      <span className="name">
        {item.name}
        {item.purchased && <span className="sr-only"> (purchased)</span>}
      </span>
      <span className="qty">
        <button
          type="button"
          onClick={() => onQuantity(item.id, -1)}
          disabled={item.quantity <= 1}
          aria-label={`Decrease quantity of ${item.name}`}
        >
          -
        </button>
        <output aria-label={`Quantity of ${item.name}`}>{item.quantity}</output>
        <button type="button" onClick={() => onQuantity(item.id, 1)} aria-label={`Increase quantity of ${item.name}`}>
          +
        </button>
      </span>
      <span className="price">{item.unitPrice > 0 ? currency.format(lineTotal) : "--"}</span>
      <button type="button" onClick={() => onDelete(item.id)} aria-label={`Delete ${item.name}`}>
        Delete
      </button>
    </li>
  );
}
