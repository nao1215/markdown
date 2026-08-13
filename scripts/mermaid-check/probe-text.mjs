// Render every ```mermaid block in the given markdown files and print the SVG
// text content, so a human can see what the drawing actually says. Rendering
// path copied from scripts/mermaid-check/check.mjs.
import { readFileSync } from "node:fs";
import { createServer } from "node:http";
import { join, normalize, resolve, sep } from "node:path";
import puppeteer from "puppeteer";

const root = resolve("/home/nao/ghq/github.com/nao1215/markdown/scripts/mermaid-check");

function mermaidBlocks(text) {
  const blocks = [];
  const lines = text.split("\n");
  let fence = null;
  let body = [];
  for (let i = 0; i < lines.length; i++) {
    const m = lines[i].match(/^\s*(`{3,})\s*([A-Za-z0-9_+-]*)\s*$/);
    if (fence === null) {
      if (m) {
        fence = { ticks: m[1], info: m[2] };
        body = [];
      }
      continue;
    }
    if (m && m[1].length >= fence.ticks.length && m[2] === "") {
      if (fence.info === "mermaid") blocks.push(body.join("\n"));
      fence = null;
      continue;
    }
    body.push(lines[i]);
  }
  return blocks;
}

const server = createServer((req, res) => {
  const path = normalize(decodeURIComponent(req.url.split("?")[0]));
  if (path === "/") {
    res.setHeader("content-type", "text/html");
    res.end("<!doctype html><html><body></body></html>");
    return;
  }
  const file = join(root, path);
  if (!file.startsWith(root + sep)) {
    res.statusCode = 403;
    res.end("");
    return;
  }
  try {
    res.setHeader("content-type", "text/javascript");
    res.end(readFileSync(file));
  } catch {
    res.statusCode = 404;
    res.end("");
  }
});
await new Promise((ready) => server.listen(0, "127.0.0.1", ready));
const origin = `http://127.0.0.1:${server.address().port}`;

const browser = await puppeteer.launch({
  headless: true,
  executablePath: process.env.PUPPETEER_EXECUTABLE_PATH || undefined,
  args: ["--no-sandbox", "--disable-dev-shm-usage"],
});
const page = await browser.newPage();
await page.goto(origin);
await page.addScriptTag({
  type: "module",
  content: `
    import mermaid from "${origin}/node_modules/mermaid/dist/mermaid.esm.min.mjs";
    mermaid.initialize({ startOnLoad: false });
    window.draw = async (id, source) => {
      try {
        const { svg } = await mermaid.render(id, source);
        const host = document.createElement("div");
        host.innerHTML = svg;
        document.body.appendChild(host);
        const text = host.textContent;
        const nodes = [...host.querySelectorAll("svg text, svg tspan, svg span, svg div")]
          .map((n) => n.textContent.trim())
          .filter((t) => t !== "");
        host.remove();
        return { text, nodes };
      } catch (e) {
        return { error: String((e && e.message) || e) };
      }
    };
    window.ready = true;`,
});
await page.waitForFunction("window.ready === true", { timeout: 60000 });

let id = 0;
for (const file of process.argv.slice(2)) {
  for (const source of mermaidBlocks(readFileSync(file, "utf8"))) {
    const out = await page.evaluate(
      (name, text) => window.draw(name, text),
      `diagram-${id++}`,
      source,
    );
    console.log(`===== ${file} =====`);
    if (out.error) {
      console.log(`ERROR: ${out.error}`);
    } else {
      console.log(JSON.stringify([...new Set(out.nodes)]));
    }
  }
}
await browser.close();
server.close();
