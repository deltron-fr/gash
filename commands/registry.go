package commands

import (
	"io"
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

// Shell holds builtin handlers and session state.
type Shell struct {
	BuiltIn map[string]CommandFunc
	History []History
}

// Executor wires pipelines and runs each command concurrently.
func (sh *Shell) Executor(p *Pipeline) {
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
			if builtInCmds, ok := sh.BuiltIn[c.Name]; ok {
				builtInCmds(sh, c)
			} else {
				sh.handleExec(c)
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


// CommandFunc is the signature for builtins.
type CommandFunc func(sh *Shell, cmd *Command) error

// NewShell creates a shell with builtin commands registered.
func NewShell() *Shell {
	sh := &Shell{
		BuiltIn: make(map[string]CommandFunc),
		History: make([]History, 0, 100),
	}

	sh.BuiltIn["echo"] = (*Shell).Echo
	sh.BuiltIn["exit"] = (*Shell).Exit
	sh.BuiltIn["history"] = (*Shell).HistoryCmd
	sh.BuiltIn["cd"] = (*Shell).Cd
	sh.BuiltIn["pwd"] = (*Shell).Pwd
	sh.BuiltIn["type"] = (*Shell).Type
	return sh
}

// BuiltInCommands documents builtin metadata for help output.
type BuiltInCommands struct {
	Name        string
	Description string
}

func Commands() map[string]BuiltInCommands {
	commands := map[string]BuiltInCommands{
		"exit": {
			Name:        "exit",
			Description: "Exit the shell",
		},
		"echo": {
			Name:        "echo",
			Description: "display a line of text",
		},
		"type": {
			Name:        "type",
			Description: "display information about command type",
		},
		"pwd": {
			Name:        "pwd",
			Description: "displays the current working directory",
		},
		"cd": {
			Name:        "cd",
			Description: "changes the shell working directory",
		},
		"history": {
			Name:        "history",
			Description: "displays the history list",
		},
	}

	return commands
}

type historyOptions struct {
	Name        string
	Description string
}

func historyArgs() map[string]historyOptions {
	options := map[string]historyOptions{
		"-r": {
			Name:        "-r",
			Description: "read the history file and append the contents to the history list",
		},
		"-w": {
			Name:        "-w",
			Description: "write the current history to the history file",
		},
		"-a": {
			Name:        "-a",
			Description: "append history lines from this session to the history file",
		},
	}
	return options
}
