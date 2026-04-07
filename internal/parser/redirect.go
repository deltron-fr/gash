package parser

import (
	"os"

	"github.com/deltron-fr/gash/internal/commands"
)

// Redirector is a redirection operator token (e.g. ">", "2>>").
type Redirector string

// Redirect describes a single redirection operator and its target path.
type Redirect struct {
	Operator Redirector
	Target   string
}

// RedirectionCommands documents supported redirection operators.
type RedirectionCommands struct {
	Name        string
	Description string
}

// Redirection returns the supported redirection operators.
func Redirection() map[string]RedirectionCommands {
	commands := map[string]RedirectionCommands{
		">": {
			Name:        ">",
			Description: "Redirect standard output",
		},
		"1>": {
			Name:        "1>",
			Description: "Redirect standard output",
		},
		"2>": {
			Name:        "2>",
			Description: "Redirect standard error",
		},
		">>": {
			Name:        ">>",
			Description: "Appending redirect standard output",
		},
		"1>>": {
			Name:        "1>>",
			Description: "Appending redirect standard output",
		},
		"2>>": {
			Name:        "2>>",
			Description: "Appending redirect standard error",
		},
	}
	return commands

}

// Apply opens the target file and attaches it to the command's IO.
func (r *Redirect) Apply(cmd *commands.Command) (*os.File, error) {
	var file *os.File
	var err error

	switch r.Operator {
	case ">", "1>":
		file, err = os.Create(r.Target)
		if err != nil {
			return nil, err
		}

		cmd.Stdout = file
	case "2>":
		file, err := os.Create(r.Target)
		if err != nil {
			return nil, err
		}

		cmd.Stderr = file
	case ">>", "1>>":
		file, err := os.OpenFile(r.Target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		cmd.Stdout = file
	case "2>>":
		file, err := os.OpenFile(r.Target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		cmd.Stderr = file
	}

	return file, nil
}
