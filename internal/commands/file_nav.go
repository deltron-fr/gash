package commands

import (
	"fmt"
	"os"

	"github.com/deltron-fr/gash/internal/shell"
)

var ErrTooManyArguments = fmt.Errorf("too many arguments")

func Pwd(_ *shell.Shell, cmd *shell.Command) error {
	path, err := os.Getwd()
	if err != nil {
		_, err = fmt.Fprint(cmd.Stderr, err)
		return err
	}

	fmt.Fprintln(cmd.Stdout, path)
	return nil
}

func Cd(_ *shell.Shell, cmd *shell.Command) error {
	if len(cmd.Args) > 1 {
		fmt.Fprint(cmd.Stderr, "too many arguments")
		return ErrTooManyArguments
	}

	filePath := cmd.Args[0]

	if filePath == "~" {
		cdHomeDir()
		return nil
	}

	err := os.Chdir(filePath)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", filePath)
		return err
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}

func cdHomeDir() {
	homePath := os.Getenv("HOME")
	err := os.Chdir(homePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
