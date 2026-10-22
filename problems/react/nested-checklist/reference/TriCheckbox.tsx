import { useEffect, useRef } from "react";
import type { CheckState } from "./model";

type Props = {
  state: CheckState;
  label: string;
  onChange: (checked: boolean) => void;
};

/** Native checkbox; `indeterminate` can only be set through the DOM property. */
export default function TriCheckbox({ state, label, onChange }: Props) {
  const ref = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (ref.current) ref.current.indeterminate = state === "mixed";
  }, [state]);

  return (
    <input
      ref={ref}
      type="checkbox"
      checked={state === "checked"}
      // Clicking a mixed checkbox checks everything (the browser reports checked = true).
      onChange={(e) => onChange(e.target.checked)}
      aria-label={label}
    />
  );
}
