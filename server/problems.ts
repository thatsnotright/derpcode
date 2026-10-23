// Loads the problem bank from problems/<type>/<name>/.
//
// Layout per problem:
//   meta.json, README.md       -> metadata + statement (shown)
//   REVIEW.md, reference/      -> hidden until "reveal"
//   everything else            -> starter files the candidate edits
import { readdir, readFile, stat } from "node:fs/promises";
import path from "node:path";

export const ROOT = path.resolve(import.meta.dirname, "..");
export const PROBLEMS_DIR = path.join(ROOT, "problems");

export type ProblemType = "go-implement" | "go-debug" | "react";

export const GROUPS: { type: ProblemType; label: string }[] = [
  { type: "go-implement", label: "Implement (Go)" },
  { type: "go-debug", label: "Debug a PR (Go)" },
  { type: "react", label: "Build a frontend (React)" },
];

const HIDDEN = new Set(["meta.json", "README.md", "REVIEW.md", "reference"]);

export interface Meta {
  title: string;
  difficulty: string;
  minutes: number;
  language: "go" | "react";
  readonly: string[];
  order: number;
}

export interface SourceFile {
  name: string;
  content: string;
}

export interface ProblemSummary extends Omit<Meta, "readonly"> {
  id: string;
  type: ProblemType;
}

export interface Problem extends ProblemSummary {
  readme: string;
  files: (SourceFile & { readonly: boolean })[];
}

const NAME_RE = /^[a-z0-9][a-z0-9-]*$/;

export function problemDir(type: string, name: string): string | null {
  if (!GROUPS.some((g) => g.type === type) || !NAME_RE.test(name)) return null;
  return path.join(PROBLEMS_DIR, type, name);
}

async function readMeta(dir: string): Promise<Meta> {
  return JSON.parse(await readFile(path.join(dir, "meta.json"), "utf8")) as Meta;
}

async function readFiles(dir: string, skip: Set<string> = new Set()): Promise<SourceFile[]> {
  const out: SourceFile[] = [];
  for (const entry of await readdir(dir, { withFileTypes: true })) {
    if (!entry.isFile() || skip.has(entry.name) || entry.name.startsWith(".")) continue;
    out.push({ name: entry.name, content: await readFile(path.join(dir, entry.name), "utf8") });
  }
  return out.sort((a, b) => fileRank(a.name) - fileRank(b.name) || a.name.localeCompare(b.name));
}

// Show the main source first, tests and styles after.
function fileRank(name: string): number {
  if (name === "App.tsx") return 0;
  if (name.endsWith("_test.go")) return 2;
  if (name.endsWith(".css")) return 3;
  return 1;
}

async function readOptional(file: string): Promise<string | null> {
  try {
    return await readFile(file, "utf8");
  } catch {
    return null;
  }
}

export async function listProblems() {
  const groups = [];
  for (const g of GROUPS) {
    const names = await readdir(path.join(PROBLEMS_DIR, g.type)).catch(() => [] as string[]);
    const problems: ProblemSummary[] = [];
    for (const name of names) {
      const dir = problemDir(g.type, name);
      if (!dir || !(await stat(dir)).isDirectory()) continue;
      const full = await readMeta(dir).catch(() => null);
      if (!full) continue; // incomplete problem folder (no/invalid meta.json)
      const { readonly: _ro, ...meta } = full;
      problems.push({ ...meta, id: `${g.type}/${name}`, type: g.type });
    }
    problems.sort((a, b) => a.order - b.order);
    groups.push({ ...g, problems });
  }
  return { groups };
}

export async function getProblem(type: string, name: string): Promise<Problem | null> {
  const dir = problemDir(type, name);
  if (!dir) return null;
  const meta = await readMeta(dir).catch(() => null);
  if (!meta) return null;
  const { readonly, ...rest } = meta;
  const files = (await readFiles(dir, HIDDEN)).map((f) => ({ ...f, readonly: readonly.includes(f.name) }));
  return {
    ...rest,
    id: `${type}/${name}`,
    type: type as ProblemType,
    readme: (await readOptional(path.join(dir, "README.md"))) ?? "",
    files,
  };
}

export async function getReveal(type: string, name: string) {
  const dir = problemDir(type, name);
  if (!dir) return null;
  const files = await readFiles(path.join(dir, "reference")).catch(() => []);
  return { files, review: await readOptional(path.join(dir, "REVIEW.md")) };
}

/** Starter files with reference/ overlaid (same name replaces, extra files added). */
export async function referenceFiles(type: string, name: string): Promise<SourceFile[]> {
  const problem = await getProblem(type, name);
  const reveal = await getReveal(type, name);
  if (!problem || !reveal) throw new Error(`unknown problem ${type}/${name}`);
  const byName = new Map(problem.files.map((f) => [f.name, { name: f.name, content: f.content }]));
  for (const f of reveal.files) byName.set(f.name, f);
  return [...byName.values()];
}
