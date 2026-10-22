// Self-check for the problem bank:
//   go-implement: starter fails its tests, reference passes (gofmt, vet, -race)
//   go-debug:     buggy code fails (vet/tests/-race), reference passes 3x without flakes
//   react:        starter and reference both bundle (skipped if esbuild/react aren't installed)
import { runGo, type RunResult } from "../server/go.ts";
import { getProblem, listProblems, referenceFiles } from "../server/problems.ts";
import { bundleReact, reactAvailable } from "../server/react.ts";

let failures = 0;

function report(id: string, ok: boolean, detail: string) {
  if (!ok) failures++;
  console.log(`${ok ? "PASS" : "FAIL"}  ${id.padEnd(36)} ${detail}`);
}

function summarize(r: RunResult): string {
  const failed = r.tests.filter((t) => t.status === "fail").length;
  const bits = [`${r.tests.length - failed}/${r.tests.length} tests`];
  if (!r.build.ok) bits.push("build error");
  if (!r.vet.ok) bits.push("vet");
  if (!r.gofmt.ok) bits.push("gofmt");
  if (r.race) bits.push("race");
  if (r.timedOut) bits.push("timeout");
  return bits.join(", ");
}

const { groups } = await listProblems();
const canBundle = await reactAvailable();

for (const group of groups) {
  for (const summary of group.problems) {
    const [type, name] = summary.id.split("/");
    const problem = (await getProblem(type, name))!;
    const starter = problem.files.map(({ name, content }) => ({ name, content }));
    const reference = await referenceFiles(type, name);

    if (type === "react") {
      if (!canBundle) {
        console.log(`SKIP  ${summary.id.padEnd(36)} esbuild/react not installed (npm install)`);
        continue;
      }
      const [s, r] = [await bundleReact(starter), await bundleReact(reference)];
      report(summary.id, s.ok && r.ok, `starter ${s.ok ? "bundles" : "BROKEN"}, reference ${r.ok ? "bundles" : "BROKEN"}`);
      for (const e of [...(s.errors ?? []), ...(r.errors ?? [])]) console.log("      " + e);
      continue;
    }

    const before = await runGo(starter);
    const after = await runGo(reference, { count: type === "go-debug" ? 3 : 1 });
    const starterShape = type === "go-implement" ? before.build.ok && !before.ok : before.build.ok && !before.ok;
    const ok = starterShape && after.ok && after.gofmt.ok;
    report(summary.id, ok, `starter: ${summarize(before)} | reference: ${summarize(after)}`);
    if (!after.ok) console.log(after.raw.split("\n").slice(0, 30).map((l) => "      " + l).join("\n"));
  }
}

console.log(failures ? `\n${failures} problem(s) failed verification` : "\nAll problems verified");
process.exit(failures ? 1 : 0);
