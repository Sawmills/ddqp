READ ${CODE_ROOT:-$HOME/Code}/agent-scripts/AGENTS.md BEFORE ANYTHING (skip if missing). If missing, also try: $HOME/repos/agent-scripts/AGENTS.md

# AGENTS.md

## Project

- This repository provides DDQP, a DataDog query parser.
- The Go module is `github.com/Sawmills/ddqp`.
- The package parses and structures DataDog queries.
- It is not a DataDog client.
- The parser uses `github.com/alecthomas/participle/v2`.
- The module declares Go 1.26.0.
- The project uses Hermit to manage dependencies and Go tools.
- The project uses Trunk for linting and code quality checks.

## Map

- `parser.go` contains parser code.
- `metricexpression.go`, `metricfilter.go`, `metricmonitor.go`, `metricquery.go`, and `queryfilter.go` contain query parser units.
- Each parser unit has an associated `_test.go` file.
- `_examples/metrics/main.go` contains a metrics example.
- `bin/` contains Hermit files.
- `scripts/test` is the repository test script.

## Commands

- Install Go dependencies: `go mod download`
- Run repository tests: `go test ./...`
- Run CI tests in the Hermit environment: `gotestsum -f standard-verbose`
- Run lint checks: `trunk check`
- Run the CI linter in the Hermit environment: `golangci-lint run`
- Export the Hermit environment in CI: `./bin/hermit env -r >> $GITHUB_ENV`

## Rules

- Keep this package focused on parsing and structuring DataDog queries.
- Do not add DataDog client behavior.
- Keep each parser unit in its own file.
- Keep its validation in the associated test file.
- Build larger query parsers from smaller parser units.
- Preserve Go 1.26.0 compatibility.
- Use Hermit-managed dependencies and Go tools.
- Run tests and lint checks for code changes.
