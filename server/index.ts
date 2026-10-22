// Local interview practice server. Binds to 127.0.0.1 only; no external network.
import { createReadStream } from "node:fs";
import { stat } from "node:fs/promises";
import http from "node:http";
import path from "node:path";
import { gofmt, goDoc, runGo } from "./go.ts";
import { getProblem, getReveal, listProblems, ROOT, type SourceFile } from "./problems.ts";
import { bundleReact, reactAvailable } from "./react.ts";

const PORT = Number(process.env.PORT ?? 5173);
const HOST = "127.0.0.1";
const WEB_DIR = path.join(ROOT, "web");
const MAX_BODY = 2 * 1024 * 1024;

const MIME: Record<string, string> = {
  ".html": "text/html; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".svg": "image/svg+xml",
  ".json": "application/json",
};

class HttpError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

function sendJson(res: http.ServerResponse, status: number, body: unknown) {
  res.writeHead(status, { "Content-Type": "application/json", "Cache-Control": "no-store" });
  res.end(JSON.stringify(body));
}

async function readJson(req: http.IncomingMessage): Promise<Record<string, unknown>> {
  let size = 0;
  const chunks: Buffer[] = [];
  for await (const chunk of req) {
    size += (chunk as Buffer).length;
    if (size > MAX_BODY) throw new HttpError(413, "request too large");
    chunks.push(chunk as Buffer);
  }
  try {
    return JSON.parse(Buffer.concat(chunks).toString("utf8")) as Record<string, unknown>;
  } catch {
    throw new HttpError(400, "invalid JSON");
  }
}

function asFiles(value: unknown): SourceFile[] {
  if (!Array.isArray(value)) throw new HttpError(400, "files must be an array");
  return value.map((f) => {
    const { name, content } = (f ?? {}) as Record<string, unknown>;
    if (typeof name !== "string" || typeof content !== "string") throw new HttpError(400, "bad file entry");
    return { name, content };
  });
}

async function serveStatic(urlPath: string, res: http.ServerResponse) {
  const rel = urlPath === "/" ? "index.html" : decodeURIComponent(urlPath).replace(/^\/+/, "");
  const file = path.resolve(WEB_DIR, rel);
  if (!file.startsWith(WEB_DIR + path.sep)) throw new HttpError(404, "not found");
  const info = await stat(file).catch(() => null);
  if (!info?.isFile()) throw new HttpError(404, "not found");
  res.writeHead(200, { "Content-Type": MIME[path.extname(file)] ?? "application/octet-stream", "Cache-Control": "no-cache" });
  createReadStream(file).pipe(res);
}

async function handle(req: http.IncomingMessage, res: http.ServerResponse) {
  const url = new URL(req.url ?? "/", `http://${HOST}`);
  const p = url.pathname;

  if (req.method === "GET" && !p.startsWith("/api/")) return serveStatic(p, res);

  if (req.method === "GET" && p === "/api/problems") return sendJson(res, 200, await listProblems());
  if (req.method === "GET" && p === "/api/status") return sendJson(res, 200, { react: await reactAvailable() });

  const m = p.match(/^\/api\/problems\/([a-z-]+)\/([a-z0-9-]+)(\/reveal)?$/);
  if (req.method === "GET" && m) {
    const [, type, name, reveal] = m;
    const body = reveal ? await getReveal(type, name) : await getProblem(type, name);
    if (!body) throw new HttpError(404, "unknown problem");
    return sendJson(res, 200, body);
  }

  if (req.method === "POST") {
    const body = await readJson(req);
    switch (p) {
      case "/api/go/run":
        return sendJson(res, 200, await runGo(asFiles(body.files)));
      case "/api/go/format":
        if (typeof body.code !== "string") throw new HttpError(400, "code required");
        return sendJson(res, 200, await gofmt(body.code));
      case "/api/go/doc":
        return sendJson(res, 200, await goDoc(String(body.query ?? "").trim()));
      case "/api/react/bundle":
        return sendJson(res, 200, await bundleReact(asFiles(body.files)));
    }
  }
  throw new HttpError(404, "not found");
}

const server = http.createServer((req, res) => {
  handle(req, res).catch((err: unknown) => {
    const status = err instanceof HttpError ? err.status : 500;
    if (status === 500) console.error(err);
    if (!res.headersSent) sendJson(res, status, { error: err instanceof Error ? err.message : String(err) });
    else res.end();
  });
});

server.listen(PORT, HOST, () => {
  console.log(`Interview practice running at http://${HOST}:${PORT}`);
});
