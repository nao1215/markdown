// Run check.mjs over the fixtures in testdata/ and assert it reports what it
// should. A checker that stops catching anything looks exactly like a checker
// with nothing to catch, so the guard needs a guard.
//
// The fixtures use the .mdfixture extension on purpose: they are broken by
// design, and a *.md name would sweep them into the repository-wide run.
import { spawnSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const fixture = (name) => join(here, "testdata", name);

const cases = [
  {
    file: "block-title-statement.mdfixture",
    expect: 'a "block" diagram has no title statement',
  },
  {
    file: "kanban-title-statement.mdfixture",
    expect: 'a "kanban" diagram has no title statement',
  },
  {
    file: "class-standalone-annotation.mdfixture",
    expect: "standalone class annotation",
  },
  {
    file: "unparsable.mdfixture",
    expect: "mermaid parse error",
  },
  {
    file: "frontmatter-unquoted-colon.mdfixture",
    expect: "mermaid parse error",
  },
  {
    file: "frontmatter-comment-title.mdfixture",
    expect: "YAML reads its value as nothing",
  },
  {
    file: "c4-title-line-break.mdfixture",
    expect: "does not appear in the rendered diagram",
  },
  {
    file: "drawn-character-missing.mdfixture",
    expect: "never reached the rendered SVG text",
  },
  {
    file: "valid.mdfixture",
    expect: null,
  },
  {
    file: "c4-escaped-title.mdfixture",
    expect: null,
  },
];

let failed = 0;
for (const c of cases) {
  const run = spawnSync(process.execPath, [join(here, "check.mjs"), fixture(c.file)], {
    encoding: "utf8",
  });
  const output = `${run.stdout}${run.stderr}`;

  if (c.expect === null) {
    if (run.status !== 0) {
      failed++;
      console.error(`${c.file}: expected a clean run, got:\n${output}`);
    }
    continue;
  }
  if (run.status === 0) {
    failed++;
    console.error(`${c.file}: expected a failure mentioning ${JSON.stringify(c.expect)}, got a clean run`);
    continue;
  }
  if (!output.includes(c.expect)) {
    failed++;
    console.error(`${c.file}: expected ${JSON.stringify(c.expect)} in:\n${output}`);
  }
}

console.log(`self test: ${cases.length} fixture(s), ${failed} failed`);
process.exit(failed === 0 ? 0 : 1);
