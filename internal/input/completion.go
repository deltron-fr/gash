package input

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/deltron-fr/gash/internal/commands"
)

// autoCompletion tries the completion sources in order and
// returns the first non-empty set of byte-slices to use for tab completion.
func autoCompletion(input string, hasCommand bool) ([][]byte, bool) {
	if !hasCommand {
		out := autoCompleteCmds(input)
		if len(out) != 0 {
			return out, false
		}

		out = autoCompleteCmdPath(input)
		if len(out) != 0 {
			return out, false
		}

		return [][]byte{}, false
	}

	out, isDir := autoCompleteFiles(input)
	if len(out) != 0 {
		return out, isDir
	}

	return [][]byte{}, false
}

func autoCompleteCmds(input string) [][]byte {
	// Return slices of bytes for matches of built-in commands.
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

// autoCompleteCmdPath looks for executable files on $PATH whose names start with
// the provided input string. Returns either a list of full matches
// or a single entry containing only the suffix to append if only
// one match is found.
func autoCompleteCmdPath(input string) [][]byte {
	pathEnv := os.Getenv("PATH")
	separator := string(os.PathListSeparator)

	directories := strings.Split(pathEnv, separator)
	matches := make([][]byte, 0, 70)
	seen := make(map[string]bool)

	for _, dir := range directories {
		files, err := os.ReadDir(dir)
		if os.IsPermission(err) {
			fmt.Fprintf(os.Stderr, "insufficient permission to read directory: %v", dir)
			continue
		}

		for _, f := range files {
			if !f.Type().IsRegular() {
				continue
			}

			if strings.HasPrefix(f.Name(), input) {
				if !seen[f.Name()] {
					matches = append(matches, []byte(f.Name()))
					seen[f.Name()] = true
				}
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

// autoCompleteFiles looks for files in the working directory whose names start
// the provided input string. Returns either a list of full matches
// or a single entry containing only the suffix to append if only
// one match is found.
func autoCompleteFiles(input string) ([][]byte, bool) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current working directory: %v", err)
		return [][]byte{}, false
	}

	matches := make([][]byte, 0, 70)

	if strings.ContainsRune(input, os.PathSeparator) {
		parts := strings.Split(input, string(os.PathSeparator))
		input = parts[len(parts)-1]

		filepath := strings.Join(parts[:len(parts)-1], string(os.PathSeparator))
		pwd = filepath
		files, err := os.ReadDir(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading directory: %v", err)
			return [][]byte{}, false
		}

		for _, f := range files {
			if strings.HasPrefix(f.Name(), input) {
				matches = append(matches, []byte(f.Name()))
			}
		}
	} else {

		files, err := os.ReadDir(pwd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading directory: %v", err)
			return [][]byte{}, false
		}

		for _, f := range files {
			if strings.HasPrefix(f.Name(), input) {
				matches = append(matches, []byte(f.Name()))
			}
		}
	}

	if len(matches) == 0 {
		if file, _ := os.ReadDir(pwd); len(file) == 1 {
			return singleMatchHelper(file[0].Name(), pwd, input)
		}

		return [][]byte{}, false
	} else if len(matches) == 1 {
		return singleMatchHelper(string(matches[0]), pwd, input)
	}

	return matches, false
}

func singleMatchHelper(singleMatch, pwd, input string) ([][]byte, bool) {
	info, err := os.Stat(filepath.Join(pwd, string(singleMatch)))
	if err != nil {
		return [][]byte{}, false
	}

	isDir := false
	if info.IsDir() {
		isDir = true
	}

	match := make([][]byte, 0, 1)
	match = append(match, []byte(singleMatch[len(input):]))

	return match, isDir
}

func checkLongestCommonPrefix(matches []string) string {
	// Returns the longest common prefix of the provided strings.
	if len(matches) == 0 {
		return ""
	}

	prefix := matches[0]

	for i := 0; i < len(prefix); i++ {
		b := prefix[i]
		for j := 1; j < len(matches); j++ {
			if i >= len(matches[j]) || matches[j][i] != b {
				return prefix[:i]
			}
		}
	}
	return prefix
}

func buildMatches(rest [][]byte) []string {
	matches := make([]string, 0, len(rest))
	for _, b := range rest {
		matches = append(matches, string(b))
	}
	sort.Strings(matches)
	return matches
}
