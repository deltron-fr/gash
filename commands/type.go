package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/deltron-fr/dshell/fs"
)

func handleType(cmdName, redirection string, pipeArgs []int, inputHistory *[]History, args ...string) {
	availableCmds := Commands()


	if len(args) == 0 {
		return
	}

	if redirection == "" && len(pipeArgs) == 0 {
		for _, arg := range args {
			_, exists := availableCmds[arg]
			if exists {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else {
				fs.CheckPath(nil, arg, "type")
			}
		}
		return
	}

	if len(pipeArgs) > 0 {
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		isExec := fs.CheckPath(nil, cmdName, "exec")
		if !isExec {
			fmt.Printf("%s: command not found\n", cmdName)
			return
		}

		idx := pipeArgs[0]
		for _, arg := range args {
			_, exists := availableCmds[arg]
			if exists {
				fmt.Fprintf(w, "%s is a shell builtin\n", arg)
			} else {
				fs.CheckPath(nil, arg, "type")
			}
		}
		w.Close()

		commands := Commands()
		if v, ok := commands[args[idx+1]]; ok {
			var cmdArgs []string
			if idx + 2 > len(args) {
				cmdArgs = []string{}
			} else {
				cmdArgs = args[idx+2:]
			}

			v.Callback(args[idx+1], "", nil, nil, cmdArgs...)
			return
		}

		cNew := exec.Command(args[idx+1], args[idx+2:]...)
		cNew.Stdin = r
		cNew.Stdout = os.Stdout
		cNew.Stderr = os.Stderr

		err = cNew.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		return
	}

	filepath := args[len(args)-1]
	args = args[:len(args)-2]

	switch redirection {
	case ">", "1>":
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		typeStdoutRedirect(file, args...)

	case "2>":
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		typeStderrRedirect(file, args...)

	case ">>", "1>>":
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		typeStdoutRedirect(file, args...)
	case "2>>":
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		typeStderrRedirect(file, args...)
	}
}

func typeStderrRedirect(f *os.File, args ...string) {
	availableCmds := Commands()

	for _, arg := range args {
		_, exists := availableCmds[arg]
		if exists {
			_, err := fmt.Fprintf(f, "%s is a shell builtin\n", arg)
			if err != nil {
				fmt.Fprintln(f, err)
				return
			}
			fmt.Fprintf(f, "\n")
		} else {
			fs.CheckPath(f, arg, "type")
		}
	}
	f.Close()
}

func typeStdoutRedirect(f *os.File, args ...string) {
	availableCmds := Commands()

	for _, arg := range args {
		_, exists := availableCmds[arg]
		if exists {
			fmt.Fprintf(f, "%s is a shell builtin\n", arg)
			fmt.Fprintf(f, "\n")
		} else {
			fs.CheckPath(f, arg, "type")
		}
	}
	f.Close()
}
