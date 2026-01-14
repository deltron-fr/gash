package commands

type commandFunc func(string, string, *[]History, ...string)

type builtInCommands struct {
	Name        string
	Description string
	Callback    commandFunc
}

func Commands() map[string]builtInCommands {
	commands := map[string]builtInCommands{
		"exit": {
			Name:        "exit",
			Description: "Exit the shell",
			Callback:    handleExit,
		},
		"echo": {
			Name:        "echo",
			Description: "display a line of text",
			Callback:    handleEcho,
		},
		"type": {
			Name:        "type",
			Description: "display information about command type",
			Callback:    handleType,
		},
		"pwd": {
			Name:        "pwd",
			Description: "displays the current working directory",
			Callback:    handlePWD,
		},
		"cd": {
			Name:        "cd",
			Description: "changes the shell working directory",
			Callback:    handleCD,
		},
		"history": {
			Name:        "history",
			Description: "displays the history list",
			Callback:    handleHistory,
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
