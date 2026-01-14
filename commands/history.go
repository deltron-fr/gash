package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type History struct {
	Counter int
	Name    string
	InFile  bool
}

func AddEntry(cmd string, history []History) *History {
	return &History{
		Counter: len(history) + 1,
		Name:    cmd,
	}
}

func handleHistory(cmdName, redirection string, inputHistory *[]History, args ...string) {
	h := *inputHistory
	if len(h) <= 0 {
		return
	}

	if len(args) <= 0 {
		for i := 0; i < len(h); i++ {
			fmt.Printf("    %d  %s\n", h[i].Counter, h[i].Name)
		}
		return
	}

	if len(args) == 1 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if n > len(h) {
			return
		}

		newSlice := h[len(h)-n:]
		for i := 0; i < len(newSlice); i++ {
			fmt.Printf("    %d  %s\n", newSlice[i].Counter, newSlice[i].Name)
		}
	}

	if len(args) == 2 {
		options := historyArgs()
		if opt, exists := options[args[0]]; !exists {
			fmt.Fprintf(os.Stderr, "%s: invalid option\n", args[0])
			return
		} else {
			switch opt.Name {
			case "-r":
				entries := readHistoryFromFile(args[1])
				if entries == nil {
					return
				}

				for _, line := range entries {
					*inputHistory = append(*inputHistory, History{
						Name:    line,
						Counter: len(*inputHistory) + 1,
						InFile:  true,
					},
					)
				}
			case "-w":
				err := writeHistoryToFile(args[1], inputHistory)
				if err != nil {
					return
				}
			case "-a":
				err := appendHistoryToFile(args[1], inputHistory)
				if err != nil {
					return
				}
			}
		}
	}
}

type HistoryOptions struct {
	Name        string
	Description string
}

func historyArgs() map[string]HistoryOptions {
	options := map[string]HistoryOptions{
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

func readHistoryFromFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: No such file or directory\n", path)
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	history := make([]string, 0, 100)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			history = append(history, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading history: %v\n", err)
	}

	return history
}

func writeHistoryToFile(path string, inputHistory *[]History) error {

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to write to this file\n", path)
		return err
	}
	defer f.Close()

	h := *inputHistory
	w := bufio.NewWriter(f)
	for i := 0; i < len(h); i++ {
		fmt.Fprintf(w, "%s\n", h[i].Name)
		h[i].InFile = true
	}

	err = w.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to flush buffer\n", path)
		return err
	}

	return nil
}

func appendHistoryToFile(path string, inputHistory *[]History) error {

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to write to this file\n", path)
		return err
	}
	defer f.Close()

	h := *inputHistory
	w := bufio.NewWriter(f)
	for i := 0; i < len(h); i++ {
		if !h[i].InFile {
			fmt.Fprintf(w, "%s\n", h[i].Name)
			h[i].InFile = true
		}
	}

	err = w.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to flush buffer\n", path)
		return err
	}

	return nil

}
