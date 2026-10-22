// Small Markdown renderer covering what the problem READMEs use:
// headings, paragraphs, (nested) lists, task lists, fenced code, tables,
// blockquotes, hr, and inline code / bold / italic / links.
import { escapeHtml, highlight, languageFor } from "./highlight.js";

function inline(text) {
  const codes = [];
  let s = text.replace(/`([^`]+)`/g, (_, c) => {
    codes.push(`<code>${escapeHtml(c)}</code>`);
    return `\u0000${codes.length - 1}\u0000`;
  });
  s = escapeHtml(s)
    .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
    .replace(/(^|[^\w*])\*([^*\s][^*]*)\*/g, "$1<em>$2</em>")
    .replace(/(^|[^\w])_([^_\s][^_]*)_(?=[^\w]|$)/g, "$1<em>$2</em>")
    .replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, (_, label, href) =>
      /^(https?:|#|\.?\/)/.test(href) ? `<a href="${href}" target="_blank" rel="noreferrer">${label}</a>` : label,
    );
  return s.replace(/\u0000(\d+)\u0000/g, (_, i) => codes[Number(i)]);
}

const LIST_RE = /^(\s*)([-*+]|\d+[.)])\s+(.*)$/;

export function renderMarkdown(src) {
  const lines = src.replace(/\r\n/g, "\n").split("\n");
  let html = "";
  let task = 0;
  let i = 0;

  const isBlockStart = (l) => /^(#{1,6}\s|```|>|\s*([-*+]|\d+[.)])\s|\|)/.test(l) || /^(-{3,}|\*{3,})\s*$/.test(l);

  function renderList(baseIndent) {
    const ordered = /\d/.test(lines[i].match(LIST_RE)[2]);
    let out = ordered ? "<ol>" : "<ul>";
    while (i < lines.length) {
      const m = lines[i].match(LIST_RE);
      if (!m || m[1].length < baseIndent) break;
      if (m[1].length > baseIndent) {
        out = out.replace(/<\/li>$/, "") + renderList(m[1].length) + "</li>";
        continue;
      }
      let body = m[3];
      i++;
      // lazy continuation lines
      while (i < lines.length && lines[i].trim() && !LIST_RE.test(lines[i]) && !isBlockStart(lines[i].trimStart())) {
        body += " " + lines[i].trim();
        i++;
      }
      const t = body.match(/^\[([ xX])\]\s+(.*)$/);
      if (t) {
        const id = task++;
        out += `<li class="task"><label><input type="checkbox" data-task="${id}"${t[1] !== " " ? " checked" : ""}> <span>${inline(t[2])}</span></label></li>`;
      } else {
        out += `<li>${inline(body)}</li>`;
      }
      while (i < lines.length && !lines[i].trim() && LIST_RE.test(lines[i + 1] ?? "")) i++;
    }
    return out + (ordered ? "</ol>" : "</ul>");
  }

  while (i < lines.length) {
    const line = lines[i];
    if (!line.trim()) {
      i++;
    } else if (line.startsWith("```")) {
      const lang = line.slice(3).trim();
      const body = [];
      for (i++; i < lines.length && !lines[i].startsWith("```"); i++) body.push(lines[i]);
      i++;
      const l = languageFor("x." + (lang === "tsx" || lang === "typescript" || lang === "jsx" ? "ts" : lang));
      html += `<pre class="md-code"><code>${highlight(body.join("\n"), l)}</code></pre>`;
    } else if (/^#{1,6}\s/.test(line)) {
      const level = line.match(/^#+/)[0].length;
      html += `<h${level}>${inline(line.slice(level).trim())}</h${level}>`;
      i++;
    } else if (/^(-{3,}|\*{3,})\s*$/.test(line)) {
      html += "<hr>";
      i++;
    } else if (line.startsWith(">")) {
      const body = [];
      for (; i < lines.length && lines[i].startsWith(">"); i++) body.push(lines[i].replace(/^>\s?/, ""));
      html += `<blockquote>${renderMarkdown(body.join("\n"))}</blockquote>`;
    } else if (line.trim().startsWith("|") && /^\s*\|?[\s:-]+\|/.test(lines[i + 1] ?? "")) {
      const cells = (l) => l.trim().replace(/^\||\|$/g, "").split("|").map((c) => inline(c.trim()));
      html += "<table><thead><tr>" + cells(line).map((c) => `<th>${c}</th>`).join("") + "</tr></thead><tbody>";
      for (i += 2; i < lines.length && lines[i].trim().startsWith("|"); i++) {
        html += "<tr>" + cells(lines[i]).map((c) => `<td>${c}</td>`).join("") + "</tr>";
      }
      html += "</tbody></table>";
    } else if (LIST_RE.test(line)) {
      html += renderList(line.match(LIST_RE)[1].length);
    } else {
      const body = [];
      for (; i < lines.length && lines[i].trim() && !isBlockStart(lines[i]); i++) body.push(lines[i].trim());
      if (!body.length) body.push(lines[i++].trim());
      html += `<p>${inline(body.join(" "))}</p>`;
    }
  }
  return html;
}
