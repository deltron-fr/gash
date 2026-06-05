package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/deltron-fr/gash/internal/shell"
)

var (
	ErrNoCompletionSpec  = errors.New("no completion specification")
	ErrFileNotExecutable = errors.New("file is not executable")
)

func completeArgs() map[string]completeOptions {
	options := map[string]completeOptions{
		"-p": {
			Name:        "-p",
			Description: "prints the completion specification registered for a given command",
		},
		"-C": {
			Name:        "-C",
			Description: "registers a completer script for a command",
		},
	}
	return options
}

func Complete(sh *shell.Shell, cmd *shell.Command) error {
	options := completeArgs()

	switch len(cmd.Args) {
	case 2:
		if opt, exists := options[cmd.Args[0]]; !exists {
			fmt.Fprintf(cmd.Stderr, "%s: invalid option\n", cmd.Args[0])
			return ErrInvalidOptions
		} else {
			switch opt.Name {
			case "-p":
				path, ok := sh.CompleteScripts[cmd.Args[1]]

				if !ok {
					fmt.Fprintf(cmd.Stderr, "%s: %s: no completion specification\n", cmd.Name, cmd.Args[1])
					return ErrNoCompletionSpec
				}

				fmt.Fprintf(cmd.Stdout, "complete -C '%s' %s\n", path, cmd.Args[1])
			}
		}
	case 3:
		if opt, exists := options[cmd.Args[0]]; !exists {
			fmt.Fprintf(cmd.Stderr, "%s: invalid option\n", cmd.Args[0])
			return ErrInvalidOptions
		} else {
			switch opt.Name {
			case "-C":
				path, err := filepath.Abs(cmd.Args[1])
				if err != nil {
					fmt.Fprintf(cmd.Stderr, "%s: unable to normalize file path\n", path)
					return err
				}

				if !sh.IsExecutable(path) {
					fmt.Fprintf(cmd.Stderr, "%s: %s: file is not executable\n", cmd.Name, path)
					return ErrFileNotExecutable
				}

				sh.CompleteScripts[cmd.Args[2]] = path
			}
		}
	default:
		return nil
	}

	return nil
}

type completeOptions struct {
	Name        string
	Description string
}
