package commands

import (
	"fmt"
	"os"
	"strconv"
)

type History struct {
    Counter int
	Name string
}

func AddEntry(cmd string, history []History) *History {
	return &History{
		Counter: len(history) + 1,
		Name: cmd,
	}
}

func handleHistory(cmdName, redirection string, inputHistory []History, args ...string) {
	if len(inputHistory) <= 0 {
		return
	}

	if len(args) <= 0 {
		for i := 0; i < len(inputHistory); i++ {
			fmt.Printf("    %d  %s\n", inputHistory[i].Counter, inputHistory[i].Name)
		}
		return
	}

	if len(args) == 1 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if n > len(inputHistory) {
			return
		}

		newSlice := inputHistory[len(inputHistory)-n:]
		for i := 0; i < len(newSlice); i++ {
			fmt.Printf("    %d  %s\n", newSlice[i].Counter, newSlice[i].Name)
		} 
	}
	
}