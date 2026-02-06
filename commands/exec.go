package commands

import (
	"fmt"
	"io"

	"os/exec"

	"github.com/deltron-fr/dshell/fs"
)

var ErrNotExec = fmt.Errorf("the provided file is not executable")

func (sh *Shell) handleExec(cmd *Command) error {
	// handleExec runs an external program when the given command
	// is not a shell builtin. It supports optional redirection of
	// stdout/stderr by opening the destination file and wiring the
	// command's output streams accordingly.
	isExec := fs.CheckPath(nil, cmd.Name, "exec")
	if !isExec {
		fmt.Fprintf(cmd.Stderr, "%s: command not found\n", cmd.Name)
		return ErrNotExec
	}

	commandExec(cmd.Stdin, cmd.Stdout, cmd.Stderr, cmd.Name, cmd.Args...)
	return nil

}

func commandExec(stdin io.Reader, stdout, stderr io.Writer, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdin = stdin
	c.Stdout = stdout
	c.Stderr = stderr

	err := c.Run()
	if err != nil {
		return
	}
}
