// DOM + storage helpers.

/** h("div.cls", {attrs/on*}, ...children) */
export function h(tag, props = {}, ...children) {
  const [name, ...classes] = tag.split(".");
  const el = document.createElement(name || "div");
  if (classes.length) el.className = classes.join(" ");
  for (const [k, v] of Object.entries(props ?? {})) {
    if (v == null || v === false) continue;
    if (k.startsWith("on")) el.addEventListener(k.slice(2).toLowerCase(), v);
    else if (k === "className") el.className += (el.className ? " " : "") + v;
    else if (k === "text") el.textContent = v;
    else if (k === "html") el.innerHTML = v;
    else if (k in el && typeof v !== "string") el[k] = v;
    else el.setAttribute(k, v === true ? "" : v);
  }
  for (const c of children.flat()) {
    if (c == null || c === false) continue;
    el.append(c instanceof Node ? c : String(c));
  }
  return el;
}

const PREFIX = "ip:";

export function load(key, fallback) {
  try {
    const raw = localStorage.getItem(PREFIX + key);
    return raw === null ? fallback : JSON.parse(raw);
  } catch {
    return fallback;
  }
}

export function save(key, value) {
  try {
    localStorage.setItem(PREFIX + key, JSON.stringify(value));
  } catch {
    /* storage unavailable: ignore */
  }
}

export function remove(key) {
  try {
    localStorage.removeItem(PREFIX + key);
  } catch {
    /* ignore */
  }
}

/** Debounced fn with .flush() (run now if pending) and .cancel(). */
export function debounce(fn, ms) {
  let t = 0;
  let pending = null;
  const debounced = (...args) => {
    clearTimeout(t);
    pending = args;
    t = setTimeout(debounced.flush, ms);
  };
  debounced.flush = () => {
    clearTimeout(t);
    if (!pending) return;
    const args = pending;
    pending = null;
    fn(...args);
  };
  debounced.cancel = () => {
    clearTimeout(t);
    pending = null;
  };
  return debounced;
}

async function request(url, body) {
  const res = await fetch(url, body === undefined ? undefined : {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error ?? `${res.status} ${res.statusText}`);
  return data;
}

const plain = (files) => files.map(({ name, content }) => ({ name, content }));

export const api = {
  status: () => request("/api/status"),
  problems: () => request("/api/problems"),
  problem: (id) => request(`/api/problems/${id}`),
  reveal: (id) => request(`/api/problems/${id}/reveal`),
  goRun: (files) => request("/api/go/run", { files: plain(files) }),
  goFormat: (code) => request("/api/go/format", { code }),
  goDoc: (query) => request("/api/go/doc", { query }),
  reactBundle: (files) => request("/api/react/bundle", { files: plain(files) }),
};

/** Line diff (LCS) for the reveal view. Returns [{type: "same"|"add"|"del", text}]. */
export function diffLines(a, b) {
  const x = a.split("\n");
  const y = b.split("\n");
  const n = x.length;
  const m = y.length;
  const dp = Array.from({ length: n + 1 }, () => new Uint32Array(m + 1));
  for (let i = n - 1; i >= 0; i--) {
    for (let j = m - 1; j >= 0; j--) {
      dp[i][j] = x[i] === y[j] ? dp[i + 1][j + 1] + 1 : Math.max(dp[i + 1][j], dp[i][j + 1]);
    }
  }
  const out = [];
  let i = 0;
  let j = 0;
  while (i < n && j < m) {
    if (x[i] === y[j]) out.push({ type: "same", text: x[i++] }), j++;
    else if (dp[i + 1][j] >= dp[i][j + 1]) out.push({ type: "del", text: x[i++] });
    else out.push({ type: "add", text: y[j++] });
  }
  while (i < n) out.push({ type: "del", text: x[i++] });
  while (j < m) out.push({ type: "add", text: y[j++] });
  return out;
}
