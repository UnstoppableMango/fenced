# Agents

This document provides guidelines for AI agents contributing to the fenced project.

## About Fenced

Fenced is a Go library and CLI tool for parsing code fences from text (e.g., Markdown files). It extracts code blocks delimited by backticks (` ``` `) or asterisks (`***`), preserving the language identifier and content.

## Repository Structure

```
.
├── cmd/              # CLI command implementation
│   ├── root.go       # Main CLI command
│   └── version.go    # Version command
├── pkg/              # Core library code
│   └── fence.go      # Fence parsing logic
├── testdata/         # Test fixtures
├── main.go           # Entry point
└── *_test.go         # Test files (Ginkgo/Gomega)
```

## Development Environment

- **Language**: Go (see go.mod for version)
- **Testing**: Ginkgo & Gomega
- **Build System**: Make + Nix
- **CI/CD**: GitHub Actions
- **Container**: Built with Nix, uses Podman

## Building and Testing

```shell
# Install dependencies
make deps

# Build
go build

# Run tests
go tool ginkgo run -r --race --trace --randomize-all

# Run CLI
go run main.go testdata/markdown.md
```

## Code Guidelines

### Core Parsing Logic (`pkg/fence.go`)

- The parser uses `bufio.Scanner` to read line-by-line
- Supports both backtick (` ``` `) and asterisk (`***`) fence delimiters
- Language identifiers come from text after opening fence
- Maintains a state machine: `inBlock` tracks whether currently inside a fence

### CLI (`cmd/`)

- Uses Cobra for command structure
- Uses Charm Bracelet's `log` for logging (respects `DEBUG` env var)
- Reads from stdin if no file path provided
- Writes extracted blocks to stdout

### Testing

- All tests use Ginkgo/Gomega BDD framework
- Test files end with `_test.go`
- Suite files: `*_suite_test.go`
- Place test fixtures in `testdata/`
- Aim for comprehensive coverage (current CI includes coverage reporting via Codecov)

## Making Changes

### Before Submitting

1. **Format code**: `nix fmt`
2. **Run tests**: `ginkgo run -r`
3. **Check git status**: `git diff` should be clean (see CI's clean job)
4. **Build succeeds**: `go build`

### Commit Messages

Follow conventional commit format when possible:
- `feat:` new features
- `fix:` bug fixes
- `test:` test additions/changes
- `docs:` documentation
- `refactor:` code refactoring
- `chore:` maintenance tasks

### Pull Requests

- PRs trigger CI pipeline that runs:
  - Build and test with race detection
  - Docker container build (via Nix)
  - Git cleanliness check
- All CI checks must pass before merging

## Common Tasks

### Adding Support for New Fence Delimiters

Edit `pkg/fence.go`:
1. Add delimiter constant (e.g., `tildes = []byte("~~~")`)
2. Update `cutPrefix()` to check for new delimiter
3. Add tests in `pkg/fence_test.go`

### Adding New CLI Commands

1. Create new file in `cmd/` (e.g., `cmd/newcommand.go`)
2. Define command using `cobra.Command`
3. Register in `cmd/root.go`'s `init()` function
4. Add corresponding tests

### Updating Dependencies

```shell
go get -u <package>
go mod tidy
```

For Nix dependencies, update `flake.lock`:
```shell
nix flake update
```

## CI/CD Pipeline

The GitHub Actions workflow (`.github/workflows/ci.yml`) runs:

1. **Build & Test Job**
   - Checks out code
   - Sets up Go
   - Runs Ginkgo tests with race detection
   - Uploads coverage to Codecov

2. **Container Job**
   - Builds Docker image using Nix
   - Tags with metadata (branch, PR, semver, SHA)
   - Pushes to GHCR on main branch

3. **Clean Job**
   - Verifies `make deps` doesn't cause git changes
   - Ensures reproducible builds
