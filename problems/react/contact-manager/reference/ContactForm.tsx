import { useEffect, useId, useRef, useState } from "react";
import type { ChangeEvent, FormEvent } from "react";
import { EMPTY_DRAFT, validate } from "./types";
import type { ContactDraft, FormErrors } from "./types";

type Props = {
  initial: ContactDraft;
  legend: string;
  submitLabel: string;
  clearOnSubmit?: boolean;
  autoFocus?: boolean;
  onSubmit: (draft: ContactDraft) => void;
  onCancel?: () => void;
};

type FieldDef = {
  key: keyof ContactDraft;
  label: string;
  type: "text" | "email" | "tel";
  required: boolean;
  autoComplete: string;
};

const FIELDS: FieldDef[] = [
  { key: "name", label: "Name (required)", type: "text", required: true, autoComplete: "name" },
  { key: "email", label: "Email (required)", type: "email", required: true, autoComplete: "email" },
  { key: "phone", label: "Phone (optional)", type: "tel", required: false, autoComplete: "tel" },
  { key: "company", label: "Company (optional)", type: "text", required: false, autoComplete: "organization" },
];

export default function ContactForm({
  initial,
  legend,
  submitLabel,
  clearOnSubmit = false,
  autoFocus = false,
  onSubmit,
  onCancel,
}: Props) {
  const uid = useId();
  const [draft, setDraft] = useState<ContactDraft>(initial);
  const [errors, setErrors] = useState<FormErrors>({});
  const formRef = useRef<HTMLFormElement>(null);

  const focusField = (key: keyof ContactDraft) => {
    const el = formRef.current?.elements.namedItem(key);
    if (el instanceof HTMLElement) el.focus();
  };

  useEffect(() => {
    if (autoFocus) {
      const el = formRef.current?.elements.namedItem("name");
      if (el instanceof HTMLElement) el.focus();
    }
  }, [autoFocus]);

  const handleChange = (key: keyof ContactDraft) => (e: ChangeEvent<HTMLInputElement>) => {
    setDraft({ ...draft, [key]: e.target.value });
    // Clear the error for a field as soon as the user starts fixing it.
    if (errors[key]) setErrors({ ...errors, [key]: undefined });
  };

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const trimmed: ContactDraft = {
      name: draft.name.trim(),
      email: draft.email.trim(),
      phone: draft.phone.trim(),
      company: draft.company.trim(),
    };
    const found = validate(trimmed);
    setErrors(found);
    const firstInvalid = FIELDS.find((f) => found[f.key]);
    if (firstInvalid) {
      focusField(firstInvalid.key);
      return;
    }
    onSubmit(trimmed);
    if (clearOnSubmit) {
      setDraft(EMPTY_DRAFT);
      focusField("name");
    }
  };

  return (
    <form ref={formRef} className="contact-form" onSubmit={handleSubmit} noValidate>
      <fieldset>
        <legend>{legend}</legend>
        {FIELDS.map((f) => {
          const inputId = `${uid}-${f.key}`;
          const errorId = `${inputId}-error`;
          const error = errors[f.key];
          return (
            <div className="field" key={f.key}>
              <label htmlFor={inputId}>{f.label}</label>
              <input
                id={inputId}
                name={f.key}
                type={f.type}
                autoComplete={f.autoComplete}
                value={draft[f.key]}
                onChange={handleChange(f.key)}
                aria-required={f.required}
                aria-invalid={error ? true : undefined}
                aria-describedby={error ? errorId : undefined}
              />
              {error && (
                <p id={errorId} className="field-error">
                  Error: {error}
                </p>
              )}
            </div>
          );
        })}
        <div className="actions">
          <button type="submit">{submitLabel}</button>
          {onCancel && (
            <button type="button" onClick={onCancel}>
              Cancel
            </button>
          )}
        </div>
      </fieldset>
    </form>
  );
}
