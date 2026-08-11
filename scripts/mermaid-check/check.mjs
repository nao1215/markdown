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

let checked = 0;
let failed = 0;

for (const file of process.argv.slice(2)) {
  for (const block of mermaidBlocks(readFileSync(file, "utf8"))) {
    checked++;
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
