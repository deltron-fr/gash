# Architecture Overview

This document gives a high-level map of how `gash` turns a line of terminal input into command execution.

## Request flow

At a high level, the shell follows this loop:

1. `cmd/gash/main.go` starts the REPL with `repl.StartRepl()`.
2. `internal/input/raw.go` reads one interactive line in raw terminal mode.
3. `internal/parser/parser.go` turns that line into argument tokens.
4. `internal/repl/repl.go` builds a pipeline, detects redirections, and checks for background execution with `&`.
5. `internal/shell/shell.go` executes builtins directly or launches external commands.

## Lexer

`gash` does not currently have a standalone lexer package. Tokenization happens inside `parser.ParseInput` in [internal/parser/parser.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/parser/parser.go).

The lexer stage is simple and inline:

- It walks the input one rune at a time.
- It splits tokens on unquoted spaces.
- It tracks single-quote and double-quote state.
- It handles backslash escaping outside quotes and limited escaping inside double quotes.

This means lexing and parsing are currently combined into one pass rather than being modeled as separate compiler-style stages.

## Parser

The parser is responsible for turning a raw line into shell arguments that the REPL can interpret.

Current responsibilities:

- Preserve quoted text as a single argument
- Apply escape rules while building tokens
- Return a flat `[]string` argument list

After `ParseInput` returns, `repl.ParsePipeline` in [internal/repl/repl.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/repl/repl.go) performs the shell-specific structural pass:

- `|` splits commands into pipeline stages
- `>`, `>>`, `1>`, `1>>`, `2>`, and `2>>` are treated as redirection operators
- `&` at the end of the command marks the pipeline as a background job

That structural pass produces a `shell.Pipeline` made up of `shell.Command` values with the correct stdin/stdout/stderr wiring.

## Executor

Execution is split between pipeline orchestration and process launching.

### Pipeline orchestration

`Shell.Executor` in [internal/shell/shell.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/shell/shell.go) is responsible for:

- Creating `io.Pipe` pairs between adjacent commands
- Running each pipeline stage concurrently
- Dispatching builtins through the shell's builtin registry
- Dispatching non-builtins to the external execution path

### External command execution

`handleExec` and `commandExec` in [internal/shell/exec.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/shell/exec.go) handle external programs:

- Resolve the command name with `exec.LookPath`
- Create an `exec.Cmd`
- Attach stdin, stdout, and stderr
- Run the process in the foreground or start it in the background

For background jobs, the shell stores job metadata in `Shell.BackgroundJobs` and sends completion updates through `Shell.JobUpdates`.

## Readline

Interactive input is implemented in [internal/input/raw.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/input/raw.go).

The readline layer is custom and built directly on terminal raw mode:

- `golang.org/x/term` switches the terminal into raw mode
- Input is read byte-by-byte from stdin
- Printable bytes are inserted into the current buffer
- Escape sequences are interpreted for left, right, up, and down arrow keys
- Backspace deletes from the current buffer
- Enter submits the line back to the REPL

This layer also owns:

- History navigation using the shell's in-memory history
- Tab completion dispatch through `internal/input/completion.go`
- Double-tab listing behavior for ambiguous completions
- Prefilling the current buffer when the REPL needs to redraw a partially completed line

## Related files

- [cmd/gash/main.go](/home/deltron/workspace/github.com/deltron-fr/gash/cmd/gash/main.go)
- [internal/repl/repl.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/repl/repl.go)
- [internal/parser/parser.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/parser/parser.go)
- [internal/parser/redirect.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/parser/redirect.go)
- [internal/shell/shell.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/shell/shell.go)
- [internal/shell/exec.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/shell/exec.go)
- [internal/input/raw.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/input/raw.go)
- [internal/input/completion.go](/home/deltron/workspace/github.com/deltron-fr/gash/internal/input/completion.go)
