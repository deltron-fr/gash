package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/deltron-fr/gash/internal/shell"
)

func Exit(sh *shell.Shell, cmd *shell.Command) error {
	histFile := os.Getenv("HISTFILE")
	loadMemoryToHistFile(sh, histFile)

	os.Exit(0)
	return nil
}

func loadMemoryToHistFile(sh *shell.Shell, path string) {
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
	h := sh.History
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
