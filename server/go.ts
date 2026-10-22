// Go toolchain integration: run tests, gofmt, go doc.
import { mkdtemp, rm, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { run } from "./exec.ts";
import type { SourceFile } from "./problems.ts";

const GO_FILE_RE = /^[A-Za-z0-9_]+\.go$/;
const GO_MOD = "module practice\n\ngo 1.22\n";

export interface TestResult {
  name: string;
  status: "pass" | "fail" | "skip";
  elapsed: number;
  output: string;
}

export interface RunResult {
  ok: boolean;
  build: { ok: boolean; output: string };
  vet: { ok: boolean; output: string };
  gofmt: { ok: boolean; files: string[] };
  tests: TestResult[];
  race: boolean;
  timedOut: boolean;
  raw: string;
  durationMs: number;
}

interface TestEvent {
  Action: string;
  Test?: string;
  Elapsed?: number;
  Output?: string;
  FailedBuild?: string;
}

async function withTempModule<T>(files: SourceFile[], fn: (dir: string) => Promise<T>): Promise<T> {
  for (const f of files) {
    if (!GO_FILE_RE.test(f.name)) throw new Error(`invalid file name: ${f.name}`);
  }
  const dir = await mkdtemp(path.join(os.tmpdir(), "practice-go-"));
  try {
    await writeFile(path.join(dir, "go.mod"), GO_MOD);
    for (const f of files) await writeFile(path.join(dir, f.name), f.content);
    return await fn(dir);
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
}

export function runGo(files: SourceFile[], opts: { race?: boolean; count?: number } = {}): Promise<RunResult> {
  const race = opts.race ?? true;
  return withTempModule(files, async (dir) => {
    const started = Date.now();
    const gofmt = await run("gofmt", ["-l", "."], { cwd: dir, timeoutMs: 10_000 });
    const vet = await run("go", ["vet", "./..."], { cwd: dir, timeoutMs: 60_000 });
    const args = ["test", "-json", `-count=${opts.count ?? 1}`, "-timeout=20s", "./..."];
    if (race) args.splice(1, 0, "-race");
    const test = await run("go", args, { cwd: dir, timeoutMs: 90_000 });

    const tests = new Map<string, TestResult>();
    let buildOutput = "";
    let buildFailed = false;
    let raw = "";
    for (const line of test.stdout.split("\n")) {
      if (!line.trim()) continue;
      let ev: TestEvent;
      try {
        ev = JSON.parse(line) as TestEvent;
      } catch {
        raw += line + "\n";
        continue;
      }
      if (ev.Action === "build-output") buildOutput += ev.Output ?? "";
      if (ev.Action === "build-fail" || ev.FailedBuild) buildFailed = true;
      if (ev.Output && ev.Action !== "build-output") raw += ev.Output;
      if (!ev.Test) continue;
      const t = tests.get(ev.Test) ?? { name: ev.Test, status: "pass", elapsed: 0, output: "" };
      if (ev.Action === "output") t.output += ev.Output ?? "";
      if (ev.Action === "pass" || ev.Action === "fail" || ev.Action === "skip") {
        t.status = ev.Action;
        t.elapsed = ev.Elapsed ?? 0;
      }
      tests.set(ev.Test, t);
    }
    // A test that never reported a result (panic, timeout) counts as failed.
    for (const t of tests.values()) {
      if (t.status === "pass" && !/--- PASS/.test(t.output)) t.status = "fail";
    }
    buildOutput += test.stderr;
    if (test.code !== 0 && tests.size === 0) buildFailed = true;
    raw = (buildOutput ? buildOutput + "\n" : "") + raw;

    const fmtFiles = gofmt.stdout.split("\n").filter(Boolean);
    const allTests = [...tests.values()];
    const topLevelFailed = allTests.some((t) => t.status === "fail");
    const result: RunResult = {
      ok: test.code === 0 && vet.code === 0 && !buildFailed && !topLevelFailed && !test.timedOut,
      build: { ok: !buildFailed, output: buildFailed ? buildOutput || vet.stderr : "" },
      vet: { ok: vet.code === 0, output: vet.stderr.replace(/^# .*\n/gm, "") },
      gofmt: { ok: fmtFiles.length === 0, files: fmtFiles },
      tests: allTests,
      race: raw.includes("WARNING: DATA RACE"),
      timedOut: test.timedOut || /panic: test timed out/.test(raw),
      raw,
      durationMs: Date.now() - started,
    };
    return result;
  });
}

export async function gofmt(code: string): Promise<{ ok: boolean; code?: string; error?: string }> {
  const r = await run("gofmt", [], { cwd: os.tmpdir(), timeoutMs: 10_000, input: code });
  return r.code === 0 ? { ok: true, code: r.stdout } : { ok: false, error: r.stderr.replace(/<standard input>/g, "line") };
}

const DOC_QUERY_RE = /^[A-Za-z0-9_./]{1,100}$/;

export async function goDoc(query: string): Promise<{ ok: boolean; output: string }> {
  if (!DOC_QUERY_RE.test(query)) return { ok: false, output: "Query must look like: strings.Builder, sync.Mutex.Lock, errors" };
  const r = await run("go", ["doc", query], { cwd: os.tmpdir(), timeoutMs: 15_000 });
  return { ok: r.code === 0, output: (r.stdout || r.stderr).trim() };
}
