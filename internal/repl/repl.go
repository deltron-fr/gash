package repl

import (
	"fmt"
	"os"

	"github.com/deltron-fr/gash/internal/commands"
	"github.com/deltron-fr/gash/internal/input"
	"github.com/deltron-fr/gash/internal/parser"
	"github.com/deltron-fr/gash/internal/shell"
)

func StartRepl() {
	// StartRepl runs the main read-eval-print loop. It prints a prompt,
	// reads a line using the raw-mode input handler, runs tab-completion
	// listing when requested, parses the input, checks for redirections,
	// and sends the command to builtins or external commands.
	//
	// `exit` builtin will terminate this process.
	var buffer string
	HistFile := os.Getenv("HISTFILE")

	sh := shell.NewShell()
	commands.RegisterBuiltins(sh)
	commands.LoadHistoryToMemory(sh, HistFile)

	for {
		sh.DrainJobUpdates()
		fmt.Print("$ ")

		input, tabMatches := input.RawModeHandler(*sh, buffer, sh.History)

		if len(tabMatches) > 0 {
			for _, match := range tabMatches {
				fmt.Fprintf(os.Stdout, "%s  ", match)
			}
			fmt.Println()
			buffer = input
			continue
		}

		buffer = ""

		if input == "" {
			continue
		}

		h := commands.AddEntry(input, sh.History)
		sh.History = append(sh.History, *h)

		args := parser.ParseInput(input)
		if args == nil {
			continue
		}

		pipeline, file, bg := ParsePipeline(args)
		sh.Executor(pipeline, bg)
		if file != nil {
			file.Close()
		}
	}
}

// ParsePipeline builds a pipeline from parsed args and applies any redirections.
// It returns the pipeline plus the last redirection file opened (if any).
func ParsePipeline(args []string) (*shell.Pipeline, *os.File, bool) {
	pipeline := shell.NewPipeline()
	isFirstArg := true
	isBackgroundJob := false
	var f *os.File

	cmd := shell.Command{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if args[len(args)-1] == "&" {
		isBackgroundJob = true
		args = args[:len(args)-1]
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "|":
			pipeline.Commands = append(pipeline.Commands, cmd)
			cmd = shell.Command{
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			isFirstArg = true

		case isFirstArg:
			cmd.Name = arg
			isFirstArg = false

		case isRedirection(arg):
			if i+1 < len(args) {
				target := args[i+1]
				r := parser.Redirect{
					Operator: parser.Redirector(arg),
					Target:   target,
				}
				f, err := r.Apply(&cmd)
				if err != nil {
					fmt.Fprint(os.Stderr, err)
					return pipeline, f, isBackgroundJob
				}
				i++
			}

		default:
			cmd.Args = append(cmd.Args, arg)
		}
	}

	if cmd.Name != "" {
		pipeline.Commands = append(pipeline.Commands, cmd)
	}

	return pipeline, f, isBackgroundJob
}

// isRedirection reports whether a token is a supported redirection operator.
func isRedirection(token string) bool {
	redirectionOperators := parser.Redirection()
	if _, ok := redirectionOperators[token]; !ok {
		return false
	}

	return true
}
