package commands

import "github.com/deltron-fr/gash/internal/shell"

func RegisterBuiltins(sh *shell.Shell) {
	sh.BuiltIn["echo"] = Echo
	sh.BuiltIn["exit"] = Exit
	sh.BuiltIn["history"] = HistoryCmd
	sh.BuiltIn["cd"] = Cd
	sh.BuiltIn["pwd"] = Pwd
	sh.BuiltIn["type"] = Type
	sh.BuiltIn["jobs"] = JobsCmd
	sh.BuiltIn["complete"] = Complete
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
		"jobs": {
			Name:        "jobs",
			Description: "displays the status of jobs in the current shell",
		},
		"complete": {
			Name:        "complete",
			Description: "displays the list of builtin commands",
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
