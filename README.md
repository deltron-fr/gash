# dshell (dsh)

A minimal Unix-like shell implemented from scratch to explore parsing, terminal control, process execution, and command I/O. This repository is a systems-learning project rather than a replacement for existing shells. It focuses on the pieces developers need to understand to build a working command interpreter.

### Goal
This project is an exploration of Unix shell internals. Rather than relying on high-level abstractions, core features such as line editing, tab completion, and command parsing are implemented from scratch. The goal is to move past library-managed behavior and gain an understanding of how a REPL interacts with terminal raw mode, process execution and file descriptor-based I/O.

---

## Running locally
- Build / run:

```bash
# run the shell directly
go run ./...
```

---

## Core architecture
- REPL
  - `repl.StartRepl` drives the read–eval–print loop: prompt, read a line, parse, detect redirections, and dispatch to builtin handlers or external commands.
- Lexer / Parser
  - `parser.ParseInput` implements tokenization with support for single and double quotes and escape sequences. It returns nil on malformed input (eg. incomplete escapes/quotes).
- Execution
  - External commands are found via `fs.CheckPath` and executed in `commands.commandExec` / `commands.handleExec`, with optional stdout/stderr redirection.
  - Pipelines are wired in `commands.Executor`, which connects command stdout/stdin using pipes and runs pipeline stages concurrently.


### Supporting subsystems
- Terminal & Input
  - Raw-mode input and a small, custom readline implementation are in `input/raw.go`. Arrow key handling is in `input/keys.go`.

- Redirections
  - Recognized redirection tokens are provided by `parser/redirect.go`. The REPL detects redirection usage and applies file redirection on each command before execution.

- Completions
  - Tab completion is implemented in `input/completion.go`.
  - Sources: builtins (`commands.Commands()`), current directory entries, and executables found on `$PATH` (`fs.CheckPath` is used to test executability).
  - Behavior: single-match returns the suffix to append; multiple matches can be listed; partial/ambiguous matches are handled via a longest-common-prefix helper.
  - Note: multi-path completion for files (merging matches from different directories) is not implemented. The code currently enumerates `$PATH` entries and local files separately.

- History
  - In-memory history and persistence live in `commands/history.go`. There is a `history` builtin to list and read/write/append a history file; `repl` loads the history file on startup.

---

## Technical deep dive
This section briefly explains the important technical decisions and known limitations.

#### Parsing and quoting
- `parser.ParseInput` operates on runes to remain UTF-8 safe. It supports:
  - Double quotes: allow escaping `"` and `\\` inside.
  - Single quotes: treat content literally (backslashes are not interpreted inside single quotes).
  - Backslashes outside quotes escape the next rune.
- The parser returns `nil` on malformed input (for example, a backslash at end-of-input inside double quotes). See `parser/parser.go` for details.


#### Terminal raw mode and readline
- Raw-mode handling is implemented in `input/raw.go` using `golang.org/x/term` to switch the TTY into raw mode. The function reads bytes and implements minimal line-editing:
  - Left/Right cursor movement, delete/backspace, tab completion hook, and up/down history navigation.
- I implemented a custom readline to learn byte-level terminal interactions rather than relying on GNU Readline. This made it easier to explore how arrow sequences, control characters, and raw I/O behave.


#### Known limitations
- Multiline editing: the current readline doesn't support multiline cursor navigation. Lines that wrap or explicit multiline input are not fully supported.
- Job control: background/foreground jobs are not implemented yet.

#### Known Bugs
- Backspace when cursor is not at end-of-line: deleting a character in the middle of the buffer can leave a visual space or otherwise corrupt the display; this is a bug to be fixed.

---

## Features
| Category | Feature | Notes |
|:---|:---:|:---|
| Tab Completion | Executables, Files, Dirs | Tab completes from builtins, cwd files, and `$PATH` executables; single match appends suffix, multiple matches can be listed. |
| History | Up/Down navigation, `history` builtin, persistent store | History loads from `$HISTFILE` and supports write/append/read ops. |
| Raw Mode | Terminal raw input + custom readline | Byte-level control for keys and editing. No multiline navigation yet. |
| Logic | Multi-level quoting, basic escape rules | See `parser/parser.go` for exact rules. |
| System | Stdout/Stderr redirection, append | Redirections handled in `commands/exec.go`; appending supported. |
| System | Pipelining | Pipelines are wired in `commands.Executor` and run concurrently. |
| Builtins | `cd`, `history`, `exit`, `echo`, `type`, `pwd` | Implemented in `commands/*.go`. |
| Coming soon | Job control, multiline editing | Listed in roadmap. |

---

## Roadmap
- Support multiline commands and robust cursor navigation for wrapped lines.
- Add job control (background jobs).
- Improve completion to merge and deduplicate entries found across multiple `$PATH` directories when completing executables.
- Add more tests: unit tests for `parser.ParseInput` and completion helpers.
