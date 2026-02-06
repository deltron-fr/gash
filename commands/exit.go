package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func (sh *Shell) Exit(cmd *Command) error {
	histFile := os.Getenv("HISTFILE")
	loadMemoryToHistFile(histFile, sh.History)

	os.Exit(0)
	return nil
}

func loadMemoryToHistFile(path string, inputHistory *[]History) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}

		fmt.Fprintf(os.Stderr, "history: could not read %s: %v\n", path, err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	h := *inputHistory
	for i := 0; i < len(h); i++ {
		if !h[i].InFile {
			fmt.Fprintf(w, "%s\n", h[i].Name)
			h[i].InFile = true
		}
	}

	err = w.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %s: unable to flush buffer\n", path)
		return
	}

}
