// Parse every ```mermaid block in the markdown files given as arguments.
//
// GitHub renders mermaid in the browser, so a diagram this repository generates
// can be syntactically wrong and still ship: nothing in the Go test suite looks
// at it, and the failure only appears as an error box on the rendered page.
// This script runs the real mermaid parser over the committed markdown so a
// broken diagram fails CI instead.
import { readFileSync } from "node:fs";
import { JSDOM } from "jsdom";

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
mermaid.initialize({ startOnLoad: false, securityLevel: "loose" });

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

let checked = 0;
let failed = 0;

for (const file of await files()) {
  for (const block of mermaidBlocks(readFileSync(file, "utf8"))) {
    checked++;
    for (const form of rejectedForms) {
      if (form.pattern.test(block.body)) {
        failed++;
        console.error(`${file}:${block.line}: ${form.reason}\n`);
      }
    }
    try {
      await mermaid.parse(block.body);
    } catch (e) {
      failed++;
      const message = (e && e.message ? e.message : String(e))
        .split("\n")
        .slice(0, 8)
        .join("\n");
      console.error(`${file}:${block.line}: mermaid parse error\n${message}\n`);
    }
  }
}

console.log(`checked ${checked} mermaid block(s), ${failed} failed`);
if (checked === 0) {
  console.error("no mermaid blocks found; the file list is probably wrong");
  process.exit(1);
}
process.exit(failed === 0 ? 0 : 1);
