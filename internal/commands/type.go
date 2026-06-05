package commands

import (
	"fmt"
	"os/exec"

	"github.com/deltron-fr/gash/internal/shell"
)

func Type(_ *shell.Shell, cmd *shell.Command) error {

	availableCmds := Commands()
	for _, arg := range cmd.Args {
		_, exists := availableCmds[arg]
		if exists {
			_, err := fmt.Fprintf(cmd.Stdout, "%s is a shell builtin\n", arg)
			if err != nil {
				fmt.Fprint(cmd.Stderr, err.Error())
				return err
			}
		} else {
			path, err := exec.LookPath(arg)
			if err != nil {
				fmt.Fprintf(cmd.Stderr, "%s: not found\n", arg)
				continue
			}

			_, err = fmt.Fprintf(cmd.Stdout, "%s is %s\n", arg, path)
			if err != nil {
				fmt.Fprint(cmd.Stderr, err.Error())
				return err
			}
		}
	}

	return nil
}
