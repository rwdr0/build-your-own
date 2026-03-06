# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go implementation of Git, built as a [CodeCrafters](https://codecrafters.io) challenge. The program implements git commands (init, cat-file, etc.) incrementally across stages.

## Build & Run

```sh
# Build and run (compiles to /tmp/codecrafters-build-git-go)
./your_program.sh <command> [args...]

# Build only
go build -o /tmp/codecrafters-build-git-go app/*.go

# Run tests
go test ./app/...

# Run a single package's tests
go test ./app/catfile/
```

**Important:** Always test in a separate directory (e.g., `/tmp/testing`) to avoid corrupting this repo's `.git` folder:

```sh
mkdir -p /tmp/testing && cd /tmp/testing
/path/to/repo/your_program.sh init
```

## Architecture

- **`app/main.go`** — Entry point with command dispatch (`switch` on `os.Args[1]`)
- **`app/<command>/`** — Each git command lives in its own package (e.g., `app/init/`, `app/catfile/`)
- **`app/utils/`** — Shared helpers: `GetArgumentsForStage()` for positional arg extraction, `RunCmd()` test helper

### Adding a new git command

1. Create a new package under `app/<command>/`
2. Export a single entry function (e.g., `HashObject()`)
3. Add a `case` to the switch in `app/main.go`

### Conventions

- Package comments reference the CodeCrafters stage number (e.g., `// Package catfile => stage #4`)
- The `init` package is imported as `_init` to avoid collision with Go's builtin
- Git objects are read by decompressing zlib data from `.git/objects/<hash[:2]>/<hash[2:]>`
- Go 1.24, no external dependencies

### Testing

- Read all tests before writing a new one
- Write tests by comparing against real git command
- Write tests in the dedicated package folder with <stage-name>\_test pattern.
- use helpers from `app/utils/testhelpers.go`
