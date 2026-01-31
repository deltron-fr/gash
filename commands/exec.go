package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/deltron-fr/dshell/fs"
)

func HandleExec(cmd, redirection string, pipeArgs []int, args ...string) {
	// HandleExec runs an external program when the given command
	// is not a shell builtin. It supports optional redirection of
	// stdout/stderr by opening the destination file and wiring the
	// command's output streams accordingly.
	if redirection == "" && len(pipeArgs) == 0 {
		isExec := fs.CheckPath(nil, cmd, "exec")
		if !isExec {
			fmt.Printf("%s: command not found\n", cmd)
			return
		}
		commandExec(os.Stdout, os.Stderr, cmd, args...)
		return
	}

	if len(pipeArgs) > 0 {
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		isExec := fs.CheckPath(nil, cmd, "exec")
		if !isExec {
			fmt.Printf("%s: command not found\n", cmd)
			return
		}

		idx := pipeArgs[0]
		c := exec.Command(cmd, args[:idx]...)
		c.Stdout = w
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		err = c.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		w.Close()

		go func() {
			err = c.Wait()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}()

		commands := Commands()
		if v, ok := commands[args[idx+1]]; ok {
			var cmdArgs []string
			if idx + 2 > len(args) {
				cmdArgs = []string{}
			} else {
				cmdArgs = args[idx+2:]
			}

			v.Callback(args[idx+1], "",  nil, nil, cmdArgs...)
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
		defer file.Close()

		commandExec(file, os.Stderr, cmd, args...)

	case "2>":
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer file.Close()

		commandExec(os.Stdout, file, cmd, args...)

	case ">>", "1>>":
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer file.Close()

		commandExec(file, os.Stderr, cmd, args...)

	case "2>>":
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer file.Close()

		commandExec(os.Stdout, file, cmd, args...)
	}
}

func commandExec(stdout, stderr io.Writer, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdout = stdout
	c.Stderr = stderr

	err := c.Run()
	if err != nil {
		return
	}
}
