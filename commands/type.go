package commands

import (
	"fmt"

	"github.com/deltron-fr/dshell/fs"
)

func (sh *Shell) Type(cmd *Command) error {

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
			fs.CheckPath(nil, arg, "type")
		}
	}

	return nil
}
