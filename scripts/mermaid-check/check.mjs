// Check every ```mermaid block in the markdown files given as arguments.
//
// GitHub renders mermaid in the browser, so a diagram this repository generates
// can be wrong and still ship: nothing in the Go test suite looks at it, and the
// failure only appears on the rendered page. This script puts the committed
// markdown through the same three failures a reader can hit:
//
//  1. parse    - mermaid rejects the source outright.
//  2. render   - mermaid parses the source but throws while drawing it. GitHub
//                reports this as "Unable to render rich display"; the parser
//                alone never sees it, so the render runs in a real browser.
//  3. meaning  - mermaid renders the source without complaining, but draws
//                something other than what the builder asked for. A title is the
//                case this repository has already shipped twice, so a declared
//                title has to show up in the drawing.
import { readFileSync } from "node:fs";
import { createServer } from "node:http";
import { dirname, join, normalize, resolve, sep } from "node:path";
import { fileURLToPath } from "node:url";
import { JSDOM } from "jsdom";
import { load as parseYAML } from "js-yaml";
import puppeteer from "puppeteer";

const here = dirname(fileURLToPath(import.meta.url));

const dom = new JSDOM("<!doctype html><html><body></body></html>", {
  pretendToBeVisual: true,
});
global.window = dom.window;
global.document = dom.window.document;
Object.defineProperty(global, "navigator", {
  value: dom.window.navigator,
  configurable: true,
});

const mermaid = (await import("mermaid")).default;
// No securityLevel override: the default is "strict", which is what GitHub
// renders these diagrams with, and the point of this script is to agree with it.
mermaid.initialize({ startOnLoad: false });

// mermaidBlocks returns every fenced block tagged `mermaid`, with the line
// number of its first content line. Fences longer than three backticks are
// respected so that a ```mermaid block quoted inside a ````text block is not
// mistaken for a diagram.
function mermaidBlocks(text) {
  const blocks = [];
  const lines = text.split("\n");
  let fence = null;
  let start = 0;
  let body = [];

  for (let i = 0; i < lines.length; i++) {
    const m = lines[i].match(/^\s*(`{3,})\s*([A-Za-z0-9_+-]*)\s*$/);
    if (fence === null) {
      if (m) {
        fence = { ticks: m[1], info: m[2] };
        start = i + 2;
        body = [];
      }
      continue;
    }
    if (m && m[1].length >= fence.ticks.length && m[2] === "") {
      if (fence.info === "mermaid") {
        blocks.push({ line: start, body: body.join("\n") });
      }
      fence = null;
      continue;
    }
    body.push(lines[i]);
  }
  return blocks;
}

// files returns the paths to check. With --stdin0 the list is read as a
// NUL-delimited stream, which is the only form that survives a path containing
// whitespace; `git ls-files -z` produces exactly that.
async function files() {
  const args = process.argv.slice(2);
  if (!args.includes("--stdin0")) {
    return args;
  }
  const chunks = [];
  for await (const chunk of process.stdin) {
    chunks.push(chunk);
  }
  return Buffer.concat(chunks)
    .toString("utf8")
    .split("\0")
    .filter((name) => name !== "");
}

// rejectedForms are constructs the mermaid npm package parses but GitHub's
// renderer does not. The parser alone cannot catch these, so they are checked by
// pattern with the failure they cause spelled out.
const rejectedForms = [
  {
    // "<<Interface>> Name" on its own line. GitHub lexes the leading "<" as a
    // relationship token and fails the whole diagram with
    // "Expecting ... ANNOTATION_START ... got DEPENDENCY".
    // Put the annotation inside the class body instead:
    //   class Name {
    //       <<Interface>>
    //   }
    // [ \t] rather than \s so the match cannot run past the end of the
    // line: the body form has "}" on the next line and would otherwise hit.
    pattern: /^[ \t]*<<[^>]*>>[ \t]+\S/m,
    reason:
      "standalone class annotation; GitHub rejects it. Put <<...>> inside the class body.",
  },
];

// header returns the front matter of a diagram and the source below it.
function header(source) {
  const lines = source.split("\n");
  if (lines[0]?.trim() !== "---") {
    return { frontMatter: "", rest: lines };
  }
  const end = lines.indexOf("---", 1);
  if (end === -1) {
    return { frontMatter: "", rest: lines };
  }
  return {
    frontMatter: lines.slice(1, end).join("\n"),
    rest: lines.slice(end + 1),
  };
}

// keyword returns the diagram keyword, e.g. "block" or "stateDiagram-v2".
function keyword(source) {
  for (const line of header(source).rest) {
    const trimmed = line.trim();
    if (trimmed === "" || trimmed.startsWith("%%")) {
      continue;
    }
    return trimmed.split(/[\s{]/, 1)[0];
  }
  return "";
}

// unquote strips the quoting mermaid accepts around a title.
function unquote(value) {
  const trimmed = value.trim();
  if (trimmed.length >= 2 && trimmed.startsWith('"') && trimmed.endsWith('"')) {
    return trimmed.slice(1, -1);
  }
  return trimmed;
}

// declaredTitle returns the title a diagram asks for: front matter first, then
// the `title` statement the flat diagram types use.
//
// The front matter goes through a YAML parser rather than a regex, because that
// is what mermaid does with it and the two disagree on exactly the inputs worth
// catching: `title: Checkout # API` carries a comment, `title: 'Sprint Board'`
// is quoted, and `title: ~` is not a string at all.
function declaredTitle(source) {
  const { frontMatter, rest } = header(source);
  if (frontMatter !== "") {
    let parsed;
    try {
      parsed = parseYAML(frontMatter);
    } catch {
      // Unparsable front matter is mermaid's error to report, and the render
      // below reports it.
      parsed = null;
    }
    if (parsed && typeof parsed === "object" && "title" in parsed) {
      // A null title means YAML read the value as something that is not there:
      // "~", or a "#" that started a comment instead of the title.
      if (parsed.title === null || parsed.title === undefined) {
        return { text: null, frontMatter: true };
      }
      return { text: String(parsed.title), frontMatter: true };
    }
  }
  for (const line of rest) {
    const m = line.match(/^\s*title\s+(\S.*)$/);
    if (m) {
      return { text: unquote(m[1]), frontMatter: false };
    }
  }
  return null;
}

// untitledDiagrams are the diagram types whose renderer draws no title at all,
// in front matter or anywhere else. A `title` statement in one of these is not a
// title but content: `block` reads it as a row and draws stray blocks labeled
// "title" and whatever followed it, and `kanban` reads it as a column. Front
// matter is still allowed, because that form is inert metadata.
const untitledDiagrams = new Set([
  "block",
  "block-beta",
  "architecture-beta",
  "mindmap",
  "kanban",
]);

// renderer runs mermaid in a real browser, because rendering needs the layout
// and text measurement that jsdom does not implement.
async function renderer() {
  const root = resolve(here);
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

  // --no-sandbox is not a preference: Ubuntu 23.10 and later, which is what the
  // CI runner and most contributors are on, block the unprivileged user
  // namespaces Chrome's sandbox needs, and it refuses to start without it. The
  // exposure is bounded instead: nothing here is fetched from the network, the
  // browser only ever sees markdown already committed to this repository, and it
  // renders it at mermaid's default "strict" security level.
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
          host.remove();
          return { text };
        } catch (e) {
          return { error: String((e && e.message) || e) };
        }
      };
      window.ready = true;`,
  });
  await page.waitForFunction("window.ready === true", { timeout: 60000 });

  let id = 0;
  return {
    draw: (source) =>
      page.evaluate((name, text) => window.draw(name, text), `diagram-${id++}`, source),
    close: async () => {
      await browser.close();
      server.close();
    },
  };
}

const draw = await renderer();
let checked = 0;
let failed = 0;

const fail = (where, message) => {
  failed++;
  console.error(`${where}: ${message}\n`);
};

for (const file of await files()) {
  for (const block of mermaidBlocks(readFileSync(file, "utf8"))) {
    checked++;
    const where = `${file}:${block.line}`;

    for (const form of rejectedForms) {
      if (form.pattern.test(block.body)) {
        fail(where, form.reason);
      }
    }

    try {
      await mermaid.parse(block.body);
    } catch (e) {
      const message = (e && e.message ? e.message : String(e))
        .split("\n")
        .slice(0, 8)
        .join("\n");
      fail(where, `mermaid parse error\n${message}`);
      continue;
    }

    const title = declaredTitle(block.body);
    const type = keyword(block.body);
    if (title && title.text === null) {
      fail(
        where,
        "the front matter declares a title, but YAML reads its value as nothing. Quote it.",
      );
      continue;
    }
    if (title && !title.frontMatter && untitledDiagrams.has(type)) {
      fail(
        where,
        `a "${type}" diagram has no title statement; mermaid reads "title ${title.text}" as diagram content. Move it to the front matter.`,
      );
      continue;
    }

    const drawn = await draw.draw(block.body);
    if (drawn.error) {
      fail(where, `mermaid render error\n${drawn.error}`);
      continue;
    }
    if (title && !untitledDiagrams.has(type) && !drawn.text.includes(title.text)) {
      fail(
        where,
        `the title "${title.text}" is declared but does not appear in the rendered diagram`,
      );
    }
  }
}

await draw.close();

console.log(`checked ${checked} mermaid block(s), ${failed} failed`);
if (checked === 0) {
  console.error("no mermaid blocks found; the file list is probably wrong");
  process.exit(1);
}
process.exit(failed === 0 ? 0 : 1);
