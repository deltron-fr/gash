# gash

`gash` is a small interactive Unix-like shell written in Go. It supports builtins, external commands, pipelines, redirection, tab completion, command history, and background jobs.

## Install

Install the binary with Go:

```bash
go install github.com/deltron-fr/gash/cmd/gash@latest
```

Or run it directly from the repository:

```bash
go run ./cmd/gash
```

## Quick start

Start the shell:

```bash
gash
```

Basic examples:

```bash
pwd
cd ~/projects
echo "hello"
ls -la | grep go
go test ./... > test.log
sleep 30 &
jobs
```

## Not yet supported

`gash` is intentionally small. These are not implemented yet:

- A shell scripting language
- Command substitution such as `$(...)`
- Job control signals such as `fg`, `bg`, `SIGTSTP`, and `SIGCONT` handling

## Architecture

Implementation notes live in [docs/](docs/README.md), including an overview of the lexer, parser, executor, and readline flow.

## Builtins

`gash` currently includes these builtins:

| Command | Description |
| --- | --- |
| `cd` | Change the current working directory |
| `pwd` | Print the current working directory |
| `echo` | Print text |
| `exit` | Exit the shell |
| `type` | Show whether a command is a builtin or an executable on `PATH` |
| `history` | Print history or read/write history files |
| `jobs` | Show background jobs tracked by the shell |
| `complete` | Register, inspect, or remove completion scripts |

## Shell features

### External commands

Commands that are not builtins are resolved through `PATH` and executed as child processes.

### Pipelines

Use `|` to connect commands:

```bash
cat go.mod | grep module
```

### Redirection

Supported redirection operators:

| Operator | Effect |
| --- | --- |
| `>` / `1>` | Redirect stdout |
| `>>` / `1>>` | Append stdout |
| `2>` | Redirect stderr |
| `2>>` | Append stderr |

Examples:

```bash
go test ./... > out.log
go test ./... 2> err.log
go test ./... >> out.log
```

### Background jobs

Append `&` to run a command in the background:

```bash
sleep 30 &
jobs
```

### Quoting and escaping

`gash` supports:

- Single quotes: `'literal text'`
- Double quotes: `"quoted text"`
- Backslash escaping outside quotes
- Escaping `"` and `\` inside double quotes

Examples:

```bash
echo 'hello world'
echo "hello world"
echo hello\ world
```

### Tab completion

Tab completion supports:

- Builtin commands
- Executables found on `PATH`
- Files and directories in the current working directory
- Custom completion scripts registered with `complete -C`

## History

If `HISTFILE` is set, `gash` loads history from that file on startup.

The `history` builtin supports:

| Command | Description |
| --- | --- |
| `history` | Print all in-memory history entries |
| `history N` | Print the last `N` entries |
| `history -r <file>` | Read entries from a file into history |
| `history -w <file>` | Write the current history to a file |
| `history -a <file>` | Append the current session history to a file |

Example:

```bash
export HISTFILE="$HOME/.gash_history"
gash
```

## Completion scripts

The `complete` builtin can register an executable script as a completer for a command:

```bash
complete -C ./scripts/my-completer git
complete -p git
complete -r git
```

The completer is executed with:

- argv: `<command-name> <current-token> <previous-token>`
- env: `COMP_LINE` and `COMP_POINT`

Each completion candidate should be written on its own line.

## Development

Run locally:

```bash
go run ./cmd/gash
```

Run tests:

```bash
go test ./...
```
