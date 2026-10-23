export type ContactDraft = {
  name: string;
  email: string;
  phone: string;
  company: string;
};

export type Contact = ContactDraft & {
  id: string;
  favorite: boolean;
};

export type SortKey = "name" | "company";

export type FormErrors = Partial<Record<keyof ContactDraft, string>>;

export const EMPTY_DRAFT: ContactDraft = { name: "", email: "", phone: "", company: "" };

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const PHONE_RE = /^[\d\s()+.-]{7,}$/;

// Expects an already-trimmed draft.
export function validate(draft: ContactDraft): FormErrors {
  const errors: FormErrors = {};
  if (!draft.name) {
    errors.name = "Name is required.";
  }
  if (!draft.email) {
    errors.email = "Email is required.";
  } else if (!EMAIL_RE.test(draft.email)) {
    errors.email = "Enter an email like name@example.com.";
  }
  if (draft.phone && !PHONE_RE.test(draft.phone)) {
    errors.phone = "Phone may only contain digits, spaces and + ( ) . - and needs at least 7 characters.";
  }
  return errors;
}
