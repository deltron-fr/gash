package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type History struct {
	Counter   int
	Name      string
	InFile    bool
	InFileArg bool
}

func AddEntry(cmd string, history []History) *History {
	return &History{
		Counter: len(history) + 1,
		Name:    cmd,
	}
}

var ErrInvalidOptions = fmt.Errorf("invalid option")

// handleHistory implements the `history` builtin. Without args it
// prints the full list. With a single numeric arg it prints the last
// N entries.n  With two args it accepts an option(-r, -w, -a) and a
// filename to read/write/append the history file.
func (sh *Shell) HistoryCmd(cmd *Command) error {
	h := sh.History
	if len(h) <= 0 {
		return nil
	}

	switch len(cmd.Args) {
	case 0:
		for i := 0; i < len(h); i++ {
			_, err := fmt.Fprintf(cmd.Stdout, "    %d  %s\n", h[i].Counter, h[i].Name)
			if err != nil {
				fmt.Fprint(cmd.Stderr, err)
				return err
			}
		}
		return nil
	case 1:
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			fmt.Fprintln(cmd.Stderr, err)
			return err
		}

		if n > len(h) {
			return nil
		}

		newSlice := h[len(h)-n:]
		for i := 0; i < len(newSlice); i++ {
			_, err = fmt.Fprintf(cmd.Stdout, "    %d  %s\n", newSlice[i].Counter, newSlice[i].Name)
			if err != nil {
				fmt.Fprint(cmd.Stderr, err)
				return err
			}
		}
	case 2:
		options := historyArgs()
		if opt, exists := options[cmd.Args[0]]; !exists {
			fmt.Fprintf(cmd.Stderr, "%s: invalid option\n", cmd.Args[0])
			return ErrInvalidOptions
		} else {
			switch opt.Name {
			case "-r":
				entries := readHistoryFromFile(cmd.Args[1])
				if entries == nil {
					return nil
				}

				for _, line := range entries {
					sh.History = append(sh.History, History{
						Name:      line,
						Counter:   len(sh.History) + 1,
						InFileArg: true,
					},
					)
				}
			case "-w":
				err := sh.writeHistoryToFile(cmd.Args[1])
				if err != nil {
					return err
				}
			case "-a":
				err := sh.appendHistoryToFile(cmd.Args[1])
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func readHistoryFromFile(path string) []string {
	// readHistoryFromFile returns non-empty lines from the provided
	// file path or nil on error.
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

func (sh *Shell) writeHistoryToFile(path string) error {

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to write to this file\n", path)
		return err
	}
	defer f.Close()

	h := sh.History
	w := bufio.NewWriter(f)
	for i := 0; i < len(h); i++ {
		fmt.Fprintf(w, "%s\n", h[i].Name)
		h[i].InFileArg = true
	}

	err = w.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to flush buffer\n", path)
		return err
	}

	return nil
}

func (sh *Shell) appendHistoryToFile(path string) error {
	// appendHistoryToFile appends any in-memory entries that are not
	// already saved (InFileArg == false) to `path`.

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to write to this file\n", path)
		return err
	}
	defer f.Close()

	h := sh.History
	w := bufio.NewWriter(f)
	for i := 0; i < len(h); i++ {
		if !h[i].InFileArg {
			fmt.Fprintf(w, "%s\n", h[i].Name)
			h[i].InFileArg = true
		}
	}

	err = w.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to flush buffer\n", path)
		return err
	}

	return nil

}

func (sh *Shell) LoadHistoryToMemory(path string) {
	// LoadHistoryToMemory reads `path` and appends its non-empty
	// lines to the provided history slice.

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}

		fmt.Fprintf(os.Stderr, "history: could not read %s: %v\n", path, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			sh.History = append(sh.History, History{
				Name:    line,
				Counter: len(sh.History) + 1,
				InFile:  true,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading history: %v\n", err)
	}
}
