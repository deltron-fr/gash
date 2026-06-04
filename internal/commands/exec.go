package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/deltron-fr/gash/internal/fs"
)

var ErrNotExec = fmt.Errorf("the provided file is not executable")

func (sh *Shell) handleExec(bg bool, cmd *Command) error {
	// handleExec runs an external program when the given command
	// is not a shell builtin. It supports optional redirection of
	// stdout/stderr by opening the destination file and wiring the
	// command's output streams accordingly.
	isExec := fs.CheckPath(nil, cmd.Name, "exec")
	if !isExec {
		fmt.Fprintf(cmd.Stderr, "%s: command not found\n", cmd.Name)
		return ErrNotExec
	}

	sh.commandExec(bg, cmd.Stdin, cmd.Stdout, cmd.Stderr, cmd.Name, cmd.Args...)
	return nil
}

func (sh *Shell) commandExec(bg bool, stdin io.Reader, stdout, stderr io.Writer, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdin = stdin
	c.Stdout = stdout
	c.Stderr = stderr

	if !bg {
		_ = c.Run()
		return
	}

	_ = c.Start()
	job := &BackgroundJob{
		ID:                len(sh.BackgroundJobs) + 1,
		BackgroundProcess: c,
		Status:            StatusRunning,
	}
	sh.BackgroundJobs = append(sh.BackgroundJobs, job)
	fmt.Fprintf(os.Stdout, "[%d] %d\n", len(sh.BackgroundJobs), c.Process.Pid)

	go func(c *exec.Cmd, job *BackgroundJob) {
		_ = c.Wait()
		job.Status = StatusDone
		sh.JobUpdates <- *job
	}(c, job)
}
