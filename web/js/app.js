import { CodeEditor } from "./editor.js";
import { renderMarkdown } from "./markdown.js";
import { GoPanel, openReveal, ReactPanel } from "./panels.js";
import { api, debounce, h, load, remove, save } from "./util.js";

const $ = (sel) => document.querySelector(sel);

const state = {
  problem: null, // { id, language, minutes, readme, files: [{name, content, readonly}] }
  files: [], // working copy
  active: 0,
  panel: null,
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
  const problem = await api.problem(id);
  state.problem = problem;
  const draft = load(draftKey(id), null);
  state.files = problem.files.map((f) => ({ ...f, content: draft?.files?.[f.name] ?? f.content }));
  state.active = Math.min(draft?.active ?? 0, state.files.length - 1);
  save("last", id);
  location.hash = id;

  document.querySelectorAll(".problem-item").forEach((b) => b.classList.toggle("active", b.dataset.id === id));
  $("#problem-title").textContent = problem.title;
  document.title = `${problem.title} · Interview Practice`;
  renderReadme();
  renderTabs();
  showActiveFile();
  setupPanel();
  timer.reset(problem.minutes);
  $("#notes-text").value = load(`notes:${id}`, "");
  $("#btn-format").hidden = problem.language !== "go";
  $("#react-hint").hidden = problem.language !== "react" || state.reactAvailable;
  $("#btn-run").textContent = problem.language === "go" ? "▶ Run tests" : "↻ Refresh preview";
}

const persist = debounce(() => {
  if (!state.problem) return;
  save(draftKey(state.problem.id), {
    active: state.active,
    files: Object.fromEntries(state.files.map((f) => [f.name, f.content])),
  });
}, 300);

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
}

function setupPanel() {
  state.panel?.destroy?.();
  const host = $("#right");
  host.replaceChildren();
  state.panel = state.problem.language === "go" ? new GoPanel(host) : new ReactPanel(host);
  if (state.problem.language === "react") state.panel.update(state.files);
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
  try {
    state.panel.show(await api.goRun(state.files));
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
  remove(draftKey(state.problem.id));
  selectProblem(state.problem.id);
}

async function reveal() {
  if (!state.problem) return;
  if (!confirm("Reveal the reference solution? Try to finish first — in the interview you won't have one.")) return;
  openReveal(await api.reveal(state.problem.id), state.files);
}

// ---------- Timer ----------

const timer = {
  total: 0,
  remaining: 0,
  endsAt: 0,
  handle: 0,
  reset(minutes) {
    this.stop();
    this.total = minutes * 60;
    this.remaining = this.total;
    $("#timer-minutes").value = String(minutes);
    this.render();
  },
  start() {
    this.endsAt = Date.now() + this.remaining * 1000;
    this.handle = setInterval(() => {
      this.remaining = Math.round((this.endsAt - Date.now()) / 1000);
      this.render();
    }, 250);
    $("#timer-toggle").textContent = "Pause";
  },
  stop() {
    clearInterval(this.handle);
    this.handle = 0;
    $("#timer-toggle").textContent = "Start";
  },
  toggle() {
    this.handle ? this.stop() : this.start();
  },
  render() {
    const r = this.remaining;
    const abs = Math.abs(r);
    $("#timer").textContent = `${r < 0 ? "-" : ""}${String(Math.floor(abs / 60)).padStart(2, "0")}:${String(abs % 60).padStart(2, "0")}`;
    $("#timer").className = "timer" + (r <= 0 ? " red" : r <= 300 ? " amber" : "");
  },
};

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
$("#timer-toggle").addEventListener("click", () => timer.toggle());
$("#timer-reset").addEventListener("click", () => timer.reset(Number($("#timer-minutes").value)));
$("#timer-minutes").addEventListener("change", (e) => timer.reset(Number(e.target.value)));
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
