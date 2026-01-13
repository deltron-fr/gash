package input

import (
	"fmt"
	"os"
	"strings"

	"github.com/deltron-fr/dshell/commands"
)

func autoCompletion(input string) [][]byte {
	out := autoCompleteCmds(input)
	if len(out) != 0 {
		return out
	}

	out = autoCompleteFiles(input)
	if len(out) != 0 {
		return out
	}

	out = autoCompleteCmdPath(input)
	if len(out) != 0 {
		return out
	}

	return [][]byte{}
}

func autoCompleteCmds(input string) [][]byte {
	matches := make([][]byte, 0, 70)

	commands := commands.Commands()
	for _, v := range commands {
		if strings.HasPrefix(v.Name, input) {
			matches = append(matches, []byte(v.Name))
		}
	}

	if len(matches) == 0 {
		return [][]byte{}
	} else if len(matches) == 1 {
		singleMatch := make([][]byte, 0, 1)
		singleMatch = append(singleMatch, []byte(matches[0][len(input):]))
		return singleMatch
	}

	return matches
}

func autoCompleteCmdPath(input string) [][]byte {
	pathEnv := os.Getenv("PATH")
	separator := string(os.PathListSeparator)

	directories := strings.Split(pathEnv, separator)
	matches := make([][]byte, 0, 70)

	for _, dir := range directories {
		files, err := os.ReadDir(dir)
		if err == os.ErrPermission {
			fmt.Fprintf(os.Stderr, "insufficient permission to read directory: %v", dir)
			continue
		}

		for _, f := range files {
			if !f.Type().IsRegular() {
				continue
			}

			if strings.HasPrefix(f.Name(), input) {
				matches = append(matches, []byte(f.Name()))
			}
		}
	}

	if len(matches) == 0 {
		return [][]byte{}
	} else if len(matches) == 1 {
		singleMatch := make([][]byte, 0, 1)
		singleMatch = append(singleMatch, []byte(matches[0][len(input):]))
		return singleMatch
	} 

	return matches
}

func autoCompleteFiles(input string) [][]byte {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current working directory: %v", err)
		return [][]byte{}
	}

	files, err := os.ReadDir(pwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading directory: %v", err)
		return [][]byte{}
	}

	matches := make([][]byte, 0, 70)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), input) {
			matches = append(matches, []byte(f.Name()))
		}
	}

	if len(matches) == 0 {
		return [][]byte{}
	} else if len(matches) == 1 {
		singleMatch := make([][]byte, 0, 1)
		singleMatch = append(singleMatch, []byte(matches[0][len(input):]))
		return singleMatch
	}

	return matches
}
