package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
)

// Command represents one executable unit with IO attached.
type Command struct {
	Name   string
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Pipeline is an ordered list of commands connected by pipes.
type Pipeline struct {
	Commands []Command
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		Commands: []Command{},
	}
}

type ProcessStatus int

const (
	StatusRunning = iota
	StatusDone
)

func (p ProcessStatus) String() string {
	switch p {
	case StatusRunning:
		return "Running"
	case StatusDone:
		return "Done"
	default:
		return ""
	}
}

type BackgroundJob struct {
	ID                int
	BackgroundProcess *exec.Cmd
	Status            ProcessStatus
}

// History stores one command entry plus bookkeeping for file persistence.
type History struct {
	Counter   int
	Name      string
	InFile    bool
	InFileArg bool
}

// CommandFunc is the signature for builtins.
type CommandFunc func(sh *Shell, cmd *Command) error

// Shell holds builtin handlers and session state.
type Shell struct {
	BuiltIn         map[string]CommandFunc
	History         []History
	BackgroundJobs  []*BackgroundJob
	JobUpdates      chan BackgroundJob
	CompleteScripts map[string]string
}

// NewShell creates a shell with initialized session state.
func NewShell() *Shell {
	return &Shell{
		BuiltIn:         make(map[string]CommandFunc),
		History:         make([]History, 0, 100),
		JobUpdates:      make(chan BackgroundJob, 15),
		CompleteScripts: make(map[string]string),
	}
}

func (sh *Shell) DrainJobUpdates() {
	select {
	case <-sh.JobUpdates:
		length := len(sh.BackgroundJobs)
		for i, job := range sh.BackgroundJobs {
			if job.Status != StatusDone {
				continue
			}

			cmdString := strings.Join(job.BackgroundProcess.Args, " ")
			marker := " "
			if i == length-2 {
				marker = "-"
			}

			if i == length-1 {
				marker = "+"
			}

			fmt.Fprintf(os.Stdout, "[%d]%s  %-24s%s\n", job.ID, marker, job.Status.String(), cmdString)
		}

		sh.BackgroundJobs = slices.DeleteFunc(sh.BackgroundJobs, func(job *BackgroundJob) bool {
			return job.Status == StatusDone
		})

	default:
		return
	}
}

// Executor wires pipelines and runs each command concurrently.
func (sh *Shell) Executor(p *Pipeline, bg bool) {
	for i := 0; i < len(p.Commands)-1; i++ {
		r, w := io.Pipe()
		p.Commands[i].Stdout = w
		p.Commands[i+1].Stdin = r
	}

	var wg sync.WaitGroup
	for i := range p.Commands {
		cmd := &p.Commands[i]
		wg.Add(1)
		go func(c *Command) {
			defer wg.Done()
			if builtInCmd, ok := sh.BuiltIn[c.Name]; ok {
				builtInCmd(sh, c)
			} else {
				sh.handleExec(bg, c)
			}

			if r, ok := c.Stdin.(*io.PipeReader); ok {
				r.Close()
			}

			if w, ok := c.Stdout.(*io.PipeWriter); ok {
				w.Close()
			}
		}(cmd)
	}
	wg.Wait()
}
