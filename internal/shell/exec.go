package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

var ErrNotExec = fmt.Errorf("the provided file is not executable")

func (sh *Shell) handleExec(bg bool, cmd *Command) error {
	path, err := exec.LookPath(cmd.Name)
	if err != nil {
		fmt.Fprintf(cmd.Stderr, "%s: command not found\n", cmd.Name)
		return ErrNotExec
	}

	sh.commandExec(path, bg, cmd.Stdin, cmd.Stdout, cmd.Stderr, cmd.Args...)
	return nil
}

func (sh *Shell) commandExec(path string, bg bool, stdin io.Reader, stdout, stderr io.Writer, args ...string) {
	c := exec.Command(path, args...)
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
