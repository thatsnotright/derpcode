// Lightweight code editor: a transparent <textarea> over a highlighted <pre>,
// with line numbers, auto-indent, tab handling and comment toggling.
// Edits go through execCommand("insertText") so native undo/redo keeps working.
import { highlight, languageFor } from "./highlight.js";

const PAIRS = { "(": ")", "[": "]", "{": "}" };

export class CodeEditor {
  constructor(parent, { onChange, onRun } = {}) {
    this.onChange = onChange ?? (() => {});
    this.onRun = onRun ?? (() => {});
    this.lang = "text";
    this.root = document.createElement("div");
    this.root.className = "code-editor";
    this.root.innerHTML = `<div class="gutter" aria-hidden="true"></div>
      <div class="code-area"><pre class="code-hl" aria-hidden="true"><code></code></pre>
      <textarea spellcheck="false" autocapitalize="off" autocomplete="off" autocorrect="off" aria-label="Code editor"></textarea></div>`;
    this.gutter = this.root.querySelector(".gutter");
    this.hl = this.root.querySelector("code");
    this.pre = this.root.querySelector("pre");
    this.ta = this.root.querySelector("textarea");
    parent.appendChild(this.root);

    this.ta.addEventListener("input", () => {
      this.render();
      this.onChange(this.ta.value);
    });
    this.ta.addEventListener("scroll", () => this.syncScroll());
    this.ta.addEventListener("keydown", (e) => this.onKeyDown(e));
  }

  setFile(name, content, readOnly) {
    this.lang = languageFor(name);
    this.indentUnit = this.lang === "go" ? "\t" : "  ";
    this.ta.value = content;
    this.ta.readOnly = readOnly;
    this.root.classList.toggle("readonly", readOnly);
    this.ta.setAttribute("aria-label", `Code editor: ${name}${readOnly ? " (read-only)" : ""}`);
    this.ta.scrollTop = 0;
    this.render();
  }

  get value() {
    return this.ta.value;
  }

  /** Replace all content (e.g. after gofmt) while keeping it undoable. */
  replaceAll(text) {
    const pos = this.ta.selectionStart;
    this.ta.focus();
    this.ta.select();
    if (!document.execCommand("insertText", false, text)) this.ta.value = text;
    this.ta.setSelectionRange(Math.min(pos, text.length), Math.min(pos, text.length));
    this.render();
    this.onChange(this.ta.value);
  }

  focus() {
    this.ta.focus();
  }

  render() {
    // Trailing newline keeps the <pre> as tall as the textarea.
    this.hl.innerHTML = highlight(this.ta.value, this.lang) + "\n";
    const lines = this.ta.value.split("\n").length;
    if (this.lineCount !== lines) {
      this.lineCount = lines;
      this.gutter.textContent = Array.from({ length: lines }, (_, i) => i + 1).join("\n");
    }
    this.syncScroll();
  }

  syncScroll() {
    this.pre.scrollTop = this.ta.scrollTop;
    this.pre.scrollLeft = this.ta.scrollLeft;
    this.gutter.scrollTop = this.ta.scrollTop;
  }

  insert(text) {
    if (!document.execCommand("insertText", false, text)) {
      this.ta.setRangeText(text, this.ta.selectionStart, this.ta.selectionEnd, "end");
      this.ta.dispatchEvent(new Event("input"));
    }
  }

  lineRange() {
    const { value, selectionStart, selectionEnd } = this.ta;
    const start = value.lastIndexOf("\n", selectionStart - 1) + 1;
    let end = value.indexOf("\n", selectionEnd - (selectionEnd > selectionStart && value[selectionEnd - 1] === "\n" ? 1 : 0));
    if (end === -1) end = value.length;
    return { start, end, text: value.slice(start, end) };
  }

  replaceLines(transform) {
    const { start, end, text } = this.lineRange();
    const next = text.split("\n").map(transform).join("\n");
    this.ta.setSelectionRange(start, end);
    this.insert(next);
    this.ta.setSelectionRange(start, start + next.length);
  }

  toggleComment() {
    const prefix = this.lang === "css" ? null : "// ";
    if (!prefix) return;
    const { text } = this.lineRange();
    const lines = text.split("\n").filter((l) => l.trim());
    const allCommented = lines.length > 0 && lines.every((l) => /^\s*\/\//.test(l));
    this.replaceLines((l) => {
      if (!l.trim()) return l;
      return allCommented ? l.replace(/^(\s*)\/\/ ?/, "$1") : l.replace(/^(\s*)/, `$1${prefix}`);
    });
  }

  onKeyDown(e) {
    const mod = e.metaKey || e.ctrlKey;
    if (mod && e.key === "Enter") {
      e.preventDefault();
      this.onRun();
      return;
    }
    if (this.ta.readOnly) return;
    const { value, selectionStart: s, selectionEnd: end } = this.ta;

    if (mod && e.key === "/") {
      e.preventDefault();
      this.toggleComment();
    } else if (e.key === "Tab") {
      e.preventDefault();
      const unit = this.indentUnit;
      if (e.shiftKey) {
        this.replaceLines((l) => (l.startsWith(unit) ? l.slice(unit.length) : l.replace(/^[ \t]{1,2}/, "")));
      } else if (s !== end && value.slice(s, end).includes("\n")) {
        this.replaceLines((l) => (l ? unit + l : l));
      } else {
        this.insert(unit);
      }
    } else if (e.key === "Enter" && !e.shiftKey && !e.altKey) {
      e.preventDefault();
      const lineStart = value.lastIndexOf("\n", s - 1) + 1;
      const indent = value.slice(lineStart, s).match(/^[ \t]*/)[0];
      const before = value[s - 1];
      const after = value[s];
      if (PAIRS[before]) {
        const inner = "\n" + indent + this.indentUnit;
        if (after === PAIRS[before]) {
          this.insert(inner + "\n" + indent);
          this.ta.setSelectionRange(s + inner.length, s + inner.length);
        } else {
          this.insert(inner);
        }
      } else {
        this.insert("\n" + indent);
      }
    } else if ((e.key === "}" || e.key === ")" || e.key === "]") && s === end) {
      // Dedent a closing bracket typed on a whitespace-only line.
      const lineStart = value.lastIndexOf("\n", s - 1) + 1;
      const typed = value.slice(lineStart, s);
      if (typed && /^[ \t]+$/.test(typed)) {
        e.preventDefault();
        const unit = this.indentUnit;
        const dedented = typed.endsWith(unit) ? typed.slice(0, -unit.length) : typed.slice(0, -1);
        this.ta.setSelectionRange(lineStart, s);
        this.insert(dedented + e.key);
      }
    }
  }
}
