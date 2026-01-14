package repl

import (
	"fmt"
	"os"

	"github.com/deltron-fr/dshell/commands"
	"github.com/deltron-fr/dshell/input"
	"github.com/deltron-fr/dshell/parser"
)

func StartRepl() {
	// StartRepl runs the main read-eval-print loop. It prints a prompt,
	// reads a line using the raw-mode input handler, runs tab-completion
	// listing when requested, parses the input, checks for redirections,
	// and sends the command to builtins or external commands.
	//
	// `exit` builtin will terminate this process.
	var buffer string
	inputHistory := make([]commands.History, 0, 400)

	var HistFile = os.Getenv("HISTFILE")
	commands.LoadHistoryToMemory(HistFile, &inputHistory)

	for {
		fmt.Print("$ ")

		input, tabMatches := input.RawModeHandler(buffer, inputHistory)

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

		h := commands.AddEntry(input, inputHistory)
		inputHistory = append(inputHistory, *h)

		var cmd string
		var extraArgs []string
		args := parser.ParseInput(input)
		if args == nil {
			continue
		}

		cmd = args[0]
		if len(args) > 1 {
			extraArgs = args[1:]
		}

		invalid := false
		var redCmd parser.RedirectionCommands

		redirectCommands := parser.Redirection()
		for i, arg := range extraArgs {
			if c, ok := redirectCommands[arg]; ok {
				if i+1 >= len(extraArgs) {
					fmt.Println("invaid command input")
					invalid = true
					break
				} else {
					if i+1 != len(extraArgs)-1 {
						fmt.Println("invalid command input, too many destination arguments")
						break
					}
					redCmd = c
				}
			}
		}

		if invalid {
			continue
		}

		builtinCmds := commands.Commands()

		if command, exists := builtinCmds[cmd]; exists {
			command.Callback(command.Name, redCmd.Name, &inputHistory, extraArgs...)
		} else {
			commands.HandleExec(cmd, redCmd.Name, extraArgs...)
		}
	}
}
