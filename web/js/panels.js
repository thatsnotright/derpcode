// Right-hand panels (Go results, React preview) and the reveal modal.
import { escapeHtml, highlight, languageFor } from "./highlight.js";
import { renderMarkdown } from "./markdown.js";
import { api, diffLines, h } from "./util.js";

// ---------- Go ----------

export class GoPanel {
  constructor(parent) {
    this.results = h("div.panel-scroll", {}, h("p.muted.pad", { text: "Press Run (⌘/Ctrl+Enter) to vet and test with -race." }));
    this.docOut = h("pre.output.doc-out", { hidden: true });
    const input = h("input", { placeholder: "e.g. strings.Builder, sort.Slice, sync.Mutex", "aria-label": "Go docs lookup" });
    const form = h(
      "form.doc-form",
      {
        onSubmit: async (e) => {
          e.preventDefault();
          if (!input.value.trim()) return;
          this.docOut.hidden = false;
          this.docOut.textContent = "Looking up…";
          try {
            this.docOut.textContent = (await api.goDoc(input.value.trim())).output;
          } catch (err) {
            this.docOut.textContent = String(err);
          }
        },
      },
      input,
      h("button", { type: "submit", text: "go doc" }),
    );
    this.el = h(
      "section.side-panel",
      { "aria-label": "Test results" },
      h("div.panel-title", {}, "Results"),
      this.results,
      h("div.docs", {}, h("div.panel-title", {}, "Go docs (offline)"), form, this.docOut),
    );
    parent.appendChild(this.el);
  }

  running() {
    this.results.replaceChildren(h("p.muted.pad", { text: "Running go vet + go test -race…" }));
  }

  error(err) {
    this.results.replaceChildren(h("pre.output.err", { text: String(err) }));
  }

  show(r) {
    const passed = r.tests.filter((t) => t.status === "pass" && !t.name.includes("/")).length;
    const failed = r.tests.filter((t) => t.status === "fail" && !t.name.includes("/")).length;
    const nodes = [
      h(
        "div.summary",
        { className: r.ok ? "ok" : "bad", role: "status" },
        h("strong", { text: r.ok ? "✓ All checks pass" : "✗ Not passing" }),
        h("span", { text: `${passed} passed · ${failed} failed` }),
        h("span.muted", { text: `${(r.durationMs / 1000).toFixed(1)}s` }),
        r.race && h("span.badge.badge-red", { text: "DATA RACE" }),
        r.timedOut && h("span.badge.badge-red", { text: "timed out" }),
        !r.vet.ok && h("span.badge.badge-amber", { text: "vet" }),
        !r.gofmt.ok && h("span.badge.badge-amber", { text: "gofmt" }),
      ),
    ];
    if (!r.build.ok) nodes.push(h("h4", { text: "Build errors" }), h("pre.output.err", { text: r.build.output }));
    for (const t of r.tests) nodes.push(testRow(t));
    if (!r.vet.ok && r.build.ok) nodes.push(h("h4", { text: "go vet" }), h("pre.output.warn", { text: r.vet.output }));
    if (!r.gofmt.ok) nodes.push(h("h4", { text: "gofmt" }), h("pre.output.warn", { text: `Not gofmt-ed: ${r.gofmt.files.join(", ")} (use Format)` }));
    nodes.push(h("details", {}, h("summary", { text: "Raw output" }), h("pre.output", { text: r.raw || "(empty)" })));
    this.results.replaceChildren(...nodes);
  }
}

function testRow(t) {
  const depth = t.name.split("/").length - 1;
  const out = h("pre.output", { text: t.output || "(no output)", hidden: t.status !== "fail" || depth > 0 });
  const icon = t.status === "pass" ? "✓" : t.status === "fail" ? "✗" : "–";
  const head = h(
    "button.test-head",
    {
      style: `padding-left:${8 + depth * 16}px`,
      "aria-expanded": String(!out.hidden),
      onClick: () => {
        out.hidden = !out.hidden;
        head.setAttribute("aria-expanded", String(!out.hidden));
      },
    },
    h("span", { className: `status ${t.status}`, text: icon }),
    h("span.test-name", { text: depth ? t.name.slice(t.name.lastIndexOf("/") + 1) : t.name }),
    h("span.muted", { text: `${t.elapsed.toFixed(2)}s` }),
  );
  return h("div.test-row", { className: t.status }, head, out);
}

// ---------- React preview ----------

// Runs inside the sandboxed iframe: forwards console output and errors to the parent.
const BRIDGE = `(function(){
  function fmt(a){ if (typeof a === "string") return a; if (a instanceof Error) return a.stack || a.message;
    try { return JSON.stringify(a, null, 1); } catch (e) { return String(a); } }
  function send(level, args){ parent.postMessage({ __preview: true, level: level, text: Array.prototype.map.call(args, fmt).join(" ") }, "*"); }
  ["log","info","warn","error"].forEach(function(l){ var o = console[l]; console[l] = function(){ send(l, arguments); o.apply(console, arguments); }; });
  window.addEventListener("error", function(e){ send("error", [e.error || e.message]); });
  window.addEventListener("unhandledrejection", function(e){ send("error", ["Unhandled rejection:", e.reason]); });
})();`;

function srcdoc(js, css) {
  return `<!doctype html><html><head><meta charset="utf-8"><style>${css.replace(/<\/style/gi, "<\\/style")}</style></head>
<body><div id="root"></div><script>${BRIDGE}<\/script><script>${js.replace(/<\/script/gi, "<\\/script")}<\/script></body></html>`;
}

export class ReactPanel {
  constructor(parent) {
    this.seq = 0;
    this.frameWrap = h("div.preview-frame");
    this.overlay = h("div.overlay", { hidden: true });
    this.logs = h("div.console-lines", { role: "log", "aria-live": "polite" });
    this.statusEl = h("span.muted");
    this.el = h(
      "section.side-panel",
      { "aria-label": "Preview" },
      h("div.panel-title", {}, "Preview ", this.statusEl, h("button.ml-auto", { text: "Reload", onClick: () => this.update(this.files, true) })),
      h("div.preview-wrap", {}, this.frameWrap, this.overlay),
      h("div.console", {}, h("div.panel-title", {}, "Console", h("button.ml-auto", { text: "Clear", onClick: () => this.logs.replaceChildren() })), this.logs),
    );
    parent.appendChild(this.el);
    this.onMessage = (e) => {
      if (e.source !== this.frame?.contentWindow || !e.data?.__preview) return;
      this.logs.append(h("div.log", { className: e.data.level, text: e.data.text }));
      this.logs.scrollTop = this.logs.scrollHeight;
    };
    window.addEventListener("message", this.onMessage);
  }

  async update(files, force = false) {
    this.files = files;
    const id = ++this.seq;
    this.statusEl.textContent = "bundling…";
    let r;
    try {
      r = await api.reactBundle(files);
    } catch (err) {
      r = { ok: false, errors: [String(err)] };
    }
    if (id !== this.seq) return;
    this.statusEl.textContent = "";
    if (!r.ok) {
      this.overlay.hidden = false;
      this.overlay.replaceChildren(h("h4", { text: "Build error" }), ...r.errors.map((e) => h("pre.output.err", { text: e })));
      return;
    }
    this.overlay.hidden = true;
    const doc = srcdoc(r.js, r.css ?? "");
    if (!force && doc === this.lastDoc) return;
    this.lastDoc = doc;
    this.logs.replaceChildren();
    this.frame = h("iframe", { title: "App preview", sandbox: "allow-scripts allow-modals allow-forms" });
    this.frame.srcdoc = doc;
    this.frameWrap.replaceChildren(this.frame);
  }

  destroy() {
    window.removeEventListener("message", this.onMessage);
  }
}

// ---------- Reveal modal ----------

export function openReveal(reveal, userFiles) {
  const prev = document.activeElement;
  let view = reveal.files.length ? "solution" : "review";
  let active = reveal.files[0]?.name;
  const body = h("div.modal-body");
  const tabs = h("div.tabs", { role: "tablist" });
  const seg = h("div.seg");

  function renderBody() {
    tabs.hidden = view === "review";
    tabs.replaceChildren(
      ...reveal.files.map((f) =>
        h("button.tab", { role: "tab", "aria-selected": String(f.name === active), className: f.name === active ? "active" : "", text: f.name, onClick: () => ((active = f.name), renderBody()) }),
      ),
    );
    seg.querySelectorAll("button").forEach((b) => b.classList.toggle("on", b.dataset.view === view));
    const ref = reveal.files.find((f) => f.name === active);
    if (view === "review") {
      body.replaceChildren(h("div.markdown.pad", { html: reveal.review ? renderMarkdown(reveal.review) : "<p class='muted'>No review notes for this problem.</p>" }));
    } else if (view === "solution" && ref) {
      body.replaceChildren(h("pre.code-view", { html: highlight(ref.content, languageFor(ref.name)) }));
    } else if (ref) {
      const mine = userFiles.find((f) => f.name === ref.name)?.content ?? "";
      const html = diffLines(mine, ref.content)
        .map((d) => `<span class="diff-${d.type}">${d.type === "add" ? "+" : d.type === "del" ? "-" : " "} ${escapeHtml(d.text)}</span>`)
        .join("\n");
      body.replaceChildren(h("p.muted.pad", { text: "− your code   + reference" }), h("pre.code-view.diff", { html }));
    }
  }

  for (const [v, label] of [["solution", "Solution"], ["diff", "Diff vs mine"], ["review", "Review notes"]]) {
    if (v !== "review" && !reveal.files.length) continue;
    const b = h("button", { text: label, onClick: () => ((view = v), renderBody()) });
    b.dataset.view = v;
    seg.append(b);
  }

  const close = () => {
    backdrop.remove();
    document.removeEventListener("keydown", onKey);
    prev?.focus?.();
  };
  const onKey = (e) => e.key === "Escape" && close();
  const closeBtn = h("button.ml-auto", { text: "Close", onClick: close });
  const backdrop = h(
    "div.modal-backdrop",
    { onClick: (e) => e.target === backdrop && close() },
    h("div.modal", { role: "dialog", "aria-modal": "true", "aria-label": "Reference solution" }, h("div.modal-head", {}, h("strong", { text: "Reference" }), seg, closeBtn), tabs, body),
  );
  document.addEventListener("keydown", onKey);
  document.body.append(backdrop);
  renderBody();
  closeBtn.focus();
}
