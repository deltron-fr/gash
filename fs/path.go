package fs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func CheckPath(f *os.File, cmdName, cmdType string) bool {
	// CheckPath searches directories in $PATH for `cmdName` and ensures
	// it is present and executable. If `cmdType` is "type" the function
	// will print the found path (optionally to `f`) and return true; if
	// `cmdType` is "exec" the function simply returns true when an
	// executable is found. Returns false when no suitable entry exists.
	pathEnv := os.Getenv("PATH")
	separator := string(os.PathListSeparator)

	directories := strings.Split(pathEnv, separator)
	for _, dir := range directories {
		cmdPath := dir + "/" + cmdName
		if !fileExists(cmdPath) {
			continue
		}

		if !isExecutable(cmdPath) {
			continue
		}

		switch cmdType {
		case "type":
			checkPathType(f, cmdName, cmdPath)
			return true
		case "exec":
			return true
		}
	}

	if cmdType == "type" {
		fmt.Fprintf(os.Stderr, "%s: not found\n", cmdName)
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return err == nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		log.Printf("error getting file info: %v", err)
		return false
	}

	mode := info.Mode()
	return mode&0111 != 0
}

func checkPathType(f *os.File, name, path string) {
	if f != nil {
		_, err := fmt.Fprintf(f, "%s is %s\n", name, path)
		if err != nil {
			fmt.Fprintln(f, err)
		}
		return
	}

	fmt.Printf("%s is %s\n", name, path)
}
