package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/deltron-fr/dshell/fs"
)

func handleEcho(cmdName, redirection string, pipeArgs []int, inputHistory *[]History, args ...string) {
	if len(args) == 0 {
		return
	}
	
	if redirection == "" && len(pipeArgs) == 0 {
		for _, w := range args {
			fmt.Printf("%s ", w)
		}
		fmt.Println()
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
		buf := bufio.NewWriter(w)
		fmt.Fprint(buf, strings.Join(args[:idx], " "))
		fmt.Fprint(buf, "\n")

		if err := buf.Flush(); err != nil {
   	 		fmt.Fprintln(os.Stderr, err)
    		return
		}
		
		w.Close()

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

		w := bufio.NewWriter(file)
		for _, arg := range args {
			fmt.Fprintf(w, "%s ", arg)
		}
		fmt.Fprintf(w, "\n")

		err = w.Flush()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	case "2>":
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer file.Close()

		for _, arg := range args {
			_, err = fmt.Printf("%s ", arg)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		fmt.Println()

	case ">>", "1>>":
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		for _, arg := range args {
			fmt.Fprintf(w, "%s ", arg)
		}
		fmt.Fprintf(w, "\n")

		err = w.Flush()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	case "2>>":
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer file.Close()

		for _, arg := range args {
			_, err = fmt.Printf("%s ", arg)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		fmt.Println()
	}

}
