// Tiny regex-based syntax highlighter for Go, TS/TSX/JS and CSS.
// Returns escaped HTML with <span class="tok-*"> wrappers.

const GO_KEYWORDS =
  "break case chan const continue default defer else fallthrough for func go goto if import interface map package range return select struct switch type var";
const GO_TYPES =
  "bool byte complex64 complex128 error float32 float64 int int8 int16 int32 int64 rune string uint uint8 uint16 uint32 uint64 uintptr any comparable true false nil iota";
const GO_BUILTINS = "append cap clear close copy delete len make max min new panic print println recover";

const TS_KEYWORDS =
  "as async await break case catch class const continue debugger default delete do else enum export extends finally for from function if implements import in instanceof interface let new of private protected public readonly return satisfies static super switch this throw try type typeof var void while with yield";
const TS_TYPES = "string number boolean unknown never any void null undefined true false object symbol bigint Record Partial Readonly";

const words = (s) => new Set(s.split(" "));

const LANGS = {
  go: { keywords: words(GO_KEYWORDS), types: words(GO_TYPES), builtins: words(GO_BUILTINS), strings: /^("(?:\\.|[^"\\\n])*"?|`[^`]*`?|'(?:\\.|[^'\\\n])*'?)/ },
  ts: { keywords: words(TS_KEYWORDS), types: words(TS_TYPES), builtins: new Set(), strings: /^("(?:\\.|[^"\\\n])*"?|'(?:\\.|[^'\\\n])*'?|`(?:\\.|[^`\\])*`?)/ },
};

export function languageFor(name) {
  const ext = name.slice(name.lastIndexOf(".") + 1).toLowerCase();
  if (ext === "go") return "go";
  if (["ts", "tsx", "js", "jsx"].includes(ext)) return "ts";
  if (ext === "css") return "css";
  return "text";
}

export function escapeHtml(s) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

const span = (cls, text) => `<span class="tok-${cls}">${escapeHtml(text)}</span>`;

function highlightCode(code, lang) {
  const spec = LANGS[lang];
  let out = "";
  let i = 0;
  while (i < code.length) {
    const rest = code.slice(i);
    let m;
    if ((m = rest.match(/^\/\/[^\n]*/)) || (m = rest.match(/^\/\*[\s\S]*?(\*\/|$)/))) {
      out += span("comment", m[0]);
    } else if ((m = rest.match(spec.strings))) {
      out += span("string", m[0]);
    } else if ((m = rest.match(/^(0[xX][\da-fA-F_]+|\d[\d_]*(\.\d+)?([eE][+-]?\d+)?)/))) {
      out += span("number", m[0]);
    } else if ((m = rest.match(/^[A-Za-z_$][\w$]*/))) {
      const w = m[0];
      const next = code.slice(i + w.length).match(/^\s*[(<]/);
      if (spec.keywords.has(w)) out += span("keyword", w);
      else if (spec.types.has(w)) out += span("type", w);
      else if (spec.builtins.has(w)) out += span("builtin", w);
      else if (/^[A-Z]/.test(w)) out += span(next ? "func" : "type", w);
      else if (next && next[0].trim() === "(") out += span("func", w);
      else out += escapeHtml(w);
    } else if ((m = rest.match(/^<\/?[A-Za-z][\w.]*/)) && lang === "ts" && /[\s(=>{,:?]$|^$/.test(code.slice(Math.max(0, i - 1), i))) {
      out += escapeHtml(m[0].startsWith("</") ? "</" : "<") + span("tag", m[0].replace(/^<\/?/, ""));
    } else {
      m = [code[i]];
      out += escapeHtml(code[i]);
    }
    i += m[0].length;
  }
  return out;
}

function highlightCss(code) {
  return escapeHtml(code)
    .replace(/(\/\*[\s\S]*?\*\/)/g, '<span class="tok-comment">$1</span>')
    .replace(/^([^{}\n<][^{}\n]*?)(\s*\{)/gm, '<span class="tok-tag">$1</span>$2')
    .replace(/([\w-]+)(\s*:)(?![^<]*<\/span>)/g, '<span class="tok-keyword">$1</span>$2')
    .replace(/(#[\da-fA-F]{3,8}\b|\b\d+(\.\d+)?(px|rem|em|%|vh|vw|s|ms)?\b)/g, '<span class="tok-number">$1</span>');
}

export function highlight(code, lang) {
  if (lang === "go" || lang === "ts") return highlightCode(code, lang);
  if (lang === "css") return highlightCss(code);
  return escapeHtml(code);
}
