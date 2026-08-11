## Contributing as a Developer
- When creating a bug report: Please follow the template and provide detailed information.
- When fixing a feature: Create a Pull Request (PR) with accompanying test code.
- When adding a feature: First, propose the feature in an Issue.
- When touching a mermaid builder: Run `make generate` to refresh the samples under `doc/`, then `make render-check`. The latter draws every diagram committed here with the real mermaid renderer, which is the only thing that catches a diagram that ships, parses, and still shows the reader something else.
- When writing a test: Put it in a test file that already exists. Each mermaid subpackage keeps its tests in the one file named after its builder (`treemap_test.go`) and its godoc examples in `examples_test.go`; the root package uses `markdown_test.go` and `index_test.go` for tests that reach inside it, `contract_test.go` for tests that go through the exported API, and `examples_test.go` for the examples. A new test viewpoint (a golden file, a fuzz target, an error case, a boundary) is a new function in one of those, not a new file.

## Contributing Outside of Coding
The following actions help boost my motivation:

- Giving a GitHub Star
- Promoting the application
- Becoming a GitHub Sponsor
