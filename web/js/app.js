import { CodeEditor } from "./editor.js";
import { renderMarkdown } from "./markdown.js";
import { GoPanel, openReveal, ReactPanel } from "./panels.js";
import { api, debounce, h, load, remove, save } from "./util.js";

const $ = (sel) => document.querySelector(sel);

const state = {
  problem: null, // { id, language, minutes, readme, files: [{name, content, readonly}] }
  files: [], // working copy
  active: 0,
  cursors: {}, // file name -> editor position, saved with the draft
  savedAt: null,
  panel: null,
  reactAvailable: false,
};

// ---------- Problem list ----------

async function loadProblemList() {
  const { groups } = await api.problems();
  const nav = $("#problem-list");
  nav.replaceChildren(
    ...groups.map((g) =>
      h(
        "div.group",
        {},
        h("div.group-label", { text: g.label }),
        ...g.problems.map((p) =>
          h(
            "button.problem-item",
            { "data-id": p.id, onClick: () => selectProblem(p.id) },
            h("span.problem-title", { text: p.title }),
            h("span.started", { text: "●", title: "In progress (saved)", hidden: !load(draftKey(p.id), null) }),
            h("span.badge", { className: `diff-${p.difficulty}`, text: p.difficulty }),
          ),
        ),
      ),
    ),
  );
  return groups;
}

// ---------- Selecting / persisting a problem ----------

const draftKey = (id) => `draft:${id}`;

async function selectProblem(id) {
  flushDraft(); // save the problem we're leaving before state changes
  const problem = await api.problem(id);
  state.problem = problem;
  const draft = load(draftKey(id), null);
  state.files = problem.files.map((f) => ({ ...f, content: draft?.files?.[f.name] ?? f.content }));
  state.active = Math.max(0, Math.min(draft?.active ?? 0, state.files.length - 1));
  state.cursors = draft?.cursors ?? {};
  state.savedAt = draft?.savedAt ?? null;
  save("last", id);
  location.hash = id;

  document.querySelectorAll(".problem-item").forEach((b) => b.classList.toggle("active", b.dataset.id === id));
  $("#problem-title").textContent = problem.title;
  document.title = `${problem.title} · Interview Practice`;
  renderReadme();
  renderTabs();
  showActiveFile();
  setupPanel();
  renderSaveStatus();
  $("#notes-text").value = load(`notes:${id}`, "");
  $("#btn-format").hidden = problem.language !== "go";
  $("#react-hint").hidden = problem.language !== "react" || state.reactAvailable;
  $("#btn-run").textContent = problem.language === "go" ? "▶ Run tests" : "↻ Refresh preview";
}

// Drafts: only files that differ from the starter are meaningful, but we store
// all of them so a problem's starter can change without clobbering work.
function saveDraft() {
  if (!state.problem) return;
  const id = state.problem.id;
  const edited = state.files.some((f, i) => f.content !== state.problem.files[i].content);
  rememberCursor();
  if (!edited && !load(draftKey(id), null)) return; // untouched: don't mark as started
  state.savedAt = Date.now();
  save(draftKey(id), {
    active: state.active,
    cursors: state.cursors,
    savedAt: state.savedAt,
    files: Object.fromEntries(state.files.map((f) => [f.name, f.content])),
  });
  document.querySelector(`.problem-item[data-id="${id}"] .started`)?.removeAttribute("hidden");
  renderSaveStatus();
}

const persist = debounce(saveDraft, 300);

/** Save now, including cursor moves that didn't trigger an edit. */
function flushDraft() {
  persist.cancel();
  saveDraft();
}

function rememberCursor() {
  const f = state.files[state.active];
  if (f) state.cursors[f.name] = editor.position();
}

function renderSaveStatus() {
  const t = state.savedAt;
  $("#save-status").textContent = t
    ? `Saved ${new Date(t).toLocaleTimeString([], { hour: "numeric", minute: "2-digit" })}`
    : "";
}

function renderReadme() {
  const el = $("#readme");
  const key = `checks:${state.problem.id}`;
  const checks = load(key, {});
  el.innerHTML = renderMarkdown(state.problem.readme);
  el.querySelectorAll("input[data-task]").forEach((box) => {
    const i = box.dataset.task;
    if (i in checks) box.checked = checks[i];
    box.addEventListener("change", () => {
      checks[i] = box.checked;
      save(key, checks);
    });
  });
  el.scrollTop = 0;
}

function renderTabs() {
  $("#tabs").replaceChildren(
    ...state.files.map((f, i) =>
      h(
        "button.tab",
        {
          role: "tab",
          className: i === state.active ? "active" : "",
          "aria-selected": String(i === state.active),
          title: f.readonly ? "Read-only (provided by the interviewer)" : f.name,
          onClick: () => {
            rememberCursor();
            state.active = i;
            renderTabs();
            showActiveFile();
            persist();
          },
        },
        f.readonly ? "🔒 " : "",
        f.name,
      ),
    ),
  );
}

function showActiveFile() {
  const f = state.files[state.active];
  editor.setFile(f.name, f.content, f.readonly);
  if (state.cursors[f.name]) editor.restorePosition(state.cursors[f.name]);
}

const runKey = (id) => `run:${id}`;

function setupPanel() {
  state.panel?.destroy?.();
  const host = $("#right");
  host.replaceChildren();
  state.panel = state.problem.language === "go" ? new GoPanel(host) : new ReactPanel(host);
  if (state.problem.language === "react") state.panel.update(state.files);
  else {
    const last = load(runKey(state.problem.id), null);
    if (last) state.panel.show(last, true);
  }
}

const refreshPreview = debounce(() => state.panel instanceof ReactPanel && state.panel.update(state.files), 600);

// ---------- Actions ----------

let running = false;
async function run() {
  if (!state.problem) return;
  if (state.problem.language === "react") return state.panel.update(state.files, true);
  if (running) return;
  running = true;
  $("#btn-run").disabled = true;
  state.panel.running();
  const id = state.problem.id;
  try {
    const result = await api.goRun(state.files);
    save(runKey(id), { ...result, raw: result.raw.slice(0, 20000) });
    if (state.problem.id === id) state.panel.show(result);
  } catch (err) {
    state.panel.error(err);
  } finally {
    running = false;
    $("#btn-run").disabled = false;
  }
}

async function format() {
  const f = state.files[state.active];
  if (f.readonly || !f.name.endsWith(".go")) return;
  const r = await api.goFormat(f.content);
  if (r.ok) editor.replaceAll(r.code);
  else if (state.panel instanceof GoPanel) state.panel.error(`gofmt: ${r.error}`);
}

function reset() {
  if (!state.problem || !confirm("Reset all files to the starter code? Your edits will be lost.")) return;
  const id = state.problem.id;
  persist.cancel();
  remove(draftKey(id));
  remove(runKey(id));
  state.problem = null; // so selectProblem's flush doesn't re-save the old draft
  document.querySelector(`.problem-item[data-id="${id}"] .started`)?.setAttribute("hidden", "");
  selectProblem(id);
}

async function reveal() {
  if (!state.problem) return;
  if (!confirm("Reveal the reference solution? Try to finish first — in the interview you won't have one.")) return;
  openReveal(await api.reveal(state.problem.id), state.files);
}

// ---------- Resizable columns ----------

function setupSplitters() {
  const main = $("#main");
  const widths = load("layout", { left: 340, right: 440 });
  const apply = () => {
    main.style.setProperty("--left", `${widths.left}px`);
    main.style.setProperty("--right", `${widths.right}px`);
  };
  apply();
  for (const handle of document.querySelectorAll(".splitter")) {
    const side = handle.dataset.side;
    handle.addEventListener("pointerdown", (e) => {
      e.preventDefault();
      handle.setPointerCapture(e.pointerId);
      const startX = e.clientX;
      const start = widths[side];
      document.body.classList.add("resizing");
      const move = (ev) => {
        const dx = ev.clientX - startX;
        widths[side] = Math.max(200, Math.min(window.innerWidth * 0.6, start + (side === "left" ? dx : -dx)));
        apply();
      };
      const up = () => {
        handle.removeEventListener("pointermove", move);
        document.body.classList.remove("resizing");
        save("layout", widths);
      };
      handle.addEventListener("pointermove", move);
      handle.addEventListener("pointerup", up, { once: true });
    });
  }
}

// ---------- Boot ----------

const editor = new CodeEditor($("#editor"), {
  onChange: (value) => {
    state.files[state.active].content = value;
    persist();
    refreshPreview();
  },
  onRun: run,
});

$("#btn-run").addEventListener("click", run);
$("#btn-format").addEventListener("click", format);
$("#btn-reset").addEventListener("click", reset);
$("#btn-reveal").addEventListener("click", reveal);
// Flush pending edits when the tab is hidden or reloaded so nothing typed is lost.
window.addEventListener("pagehide", flushDraft);
document.addEventListener("visibilitychange", () => document.hidden && flushDraft());
$("#btn-notes").addEventListener("click", () => {
  const notes = $("#notes");
  notes.hidden = !notes.hidden;
  $("#btn-notes").setAttribute("aria-expanded", String(!notes.hidden));
  if (!notes.hidden) $("#notes-text").focus();
});
$("#notes-text").addEventListener("input", debounce((e) => state.problem && save(`notes:${state.problem.id}`, e.target.value), 300));
document.addEventListener("keydown", (e) => {
  if ((e.metaKey || e.ctrlKey) && e.key === "Enter" && e.target !== editor.ta) {
    e.preventDefault();
    run();
  }
  if ((e.metaKey || e.ctrlKey) && e.key === "s") {
    e.preventDefault(); // drafts autosave; avoid the browser's save dialog
    if (state.problem?.language === "go") format();
  }
});

setupSplitters();
const groups = await loadProblemList();
state.reactAvailable = await api.status().then((s) => s.react).catch(() => false);
const wanted = decodeURIComponent(location.hash.slice(1)) || load("last", null);
const all = groups.flatMap((g) => g.problems.map((p) => p.id));
await selectProblem(all.includes(wanted) ? wanted : all[0]);
