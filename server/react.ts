// Bundles a React exercise for the preview iframe with esbuild.
// esbuild, react and react-dom are optional dependencies: without them the
// Go rounds still work and the preview explains how to enable itself.
import { mkdtemp, rm, writeFile } from "node:fs/promises";
import { createRequire } from "node:module";
import os from "node:os";
import path from "node:path";
import { ROOT, type SourceFile } from "./problems.ts";

const FILE_RE = /^[A-Za-z0-9_-]+\.(tsx|ts|jsx|js|css)$/;
const NODE_MODULES = path.join(ROOT, "node_modules");
const MISSING_HELP =
  "React preview needs local packages that are not installed yet.\n" +
  "Run `npm install` in the project folder (installs esbuild, react, react-dom), then restart `npm start`.";

const ENTRY = `import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
createRoot(document.getElementById("root")!).render(<StrictMode><App /></StrictMode>);
`;

export interface BundleResult {
  ok: boolean;
  js?: string;
  css?: string;
  errors?: string[];
}

type Esbuild = typeof import("esbuild");
let esbuildPromise: Promise<Esbuild | null> | null = null;

function loadEsbuild(): Promise<Esbuild | null> {
  esbuildPromise ??= (async () => {
    try {
      const require = createRequire(path.join(ROOT, "package.json"));
      require.resolve("react");
      require.resolve("react-dom/client");
      return (await import(require.resolve("esbuild"))) as Esbuild;
    } catch {
      esbuildPromise = null; // retry after the user installs
      return null;
    }
  })();
  return esbuildPromise;
}

export async function reactAvailable(): Promise<boolean> {
  return (await loadEsbuild()) !== null;
}

export async function bundleReact(files: SourceFile[]): Promise<BundleResult> {
  for (const f of files) {
    if (!FILE_RE.test(f.name)) return { ok: false, errors: [`invalid file name: ${f.name}`] };
  }
  if (!files.some((f) => /^App\.(tsx|jsx|ts|js)$/.test(f.name))) {
    return { ok: false, errors: ["An App.tsx (or App.jsx) file with a default export is required."] };
  }
  const esbuild = await loadEsbuild();
  if (!esbuild) return { ok: false, errors: [MISSING_HELP] };

  const dir = await mkdtemp(path.join(os.tmpdir(), "practice-react-"));
  try {
    for (const f of files) await writeFile(path.join(dir, f.name), f.content);
    await writeFile(path.join(dir, "__entry.tsx"), ENTRY);
    const result = await esbuild.build({
      entryPoints: [path.join(dir, "__entry.tsx")],
      absWorkingDir: dir,
      nodePaths: [NODE_MODULES],
      bundle: true,
      write: false,
      outdir: path.join(dir, "out"),
      format: "iife",
      jsx: "automatic",
      target: "es2020",
      sourcemap: "inline",
      define: { "process.env.NODE_ENV": '"development"' },
      logLevel: "silent",
    });
    let js = "";
    let css = "";
    for (const out of result.outputFiles) {
      if (out.path.endsWith(".js")) js += out.text;
      else if (out.path.endsWith(".css")) css += out.text;
    }
    return { ok: true, js, css };
  } catch (err) {
    const e = err as { errors?: { text: string; location?: { file: string; line: number; column: number; lineText: string } | null }[] };
    if (!e.errors) return { ok: false, errors: [String(err)] };
    return {
      ok: false,
      errors: e.errors.map((m) =>
        m.location
          ? `${path.basename(m.location.file)}:${m.location.line}:${m.location.column}: ${m.text}\n    ${m.location.lineText}`
          : m.text,
      ),
    };
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
}
