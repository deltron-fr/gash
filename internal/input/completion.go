package input

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/deltron-fr/gash/internal/commands"
	"github.com/deltron-fr/gash/internal/shell"
)

// autoCompletion tries the completion sources and returns the first non-empty result.
// If no matches are found, it returns an empty slice.
func autoCompletion(sh shell.Shell, input, commandName string, hasCommand bool) [][]byte {
	if !hasCommand {
		out := autoCompleteCmds(input)
		if len(out) != 0 {
			return out
		}

		out = autoCompleteCmdPath(input)
		if len(out) != 0 {
			return out
		}

		return [][]byte{}
	}

	path, ok := sh.CompleteScripts[commandName]
	if ok {
		return singleMatchHelper(runCompleteScript(path), input)
	}

	out := autoCompleteFiles(input)
	if len(out) != 0 {
		return out
	}

	return [][]byte{}
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
		return singleMatchHelper(string(matches[0]), input)
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
func autoCompleteFiles(input string) [][]byte {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current working directory: %v", err)
		return [][]byte{}
	}

	matches := make([][]byte, 0, 30)

	if strings.ContainsRune(input, os.PathSeparator) {
		parts := strings.Split(input, string(os.PathSeparator))
		input = parts[len(parts)-1]

		filepath := strings.Join(parts[:len(parts)-1], string(os.PathSeparator))
		pwd = filepath
		files, err := os.ReadDir(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading directory: %v", err)
			return [][]byte{}
		}

		for _, f := range files {
			if strings.HasPrefix(f.Name(), input) {
				if !f.IsDir() {
					matches = append(matches, []byte(f.Name()+" "))
				} else {
					matches = append(matches, []byte(f.Name()+string(os.PathSeparator)))
				}
			}
		}
	} else {

		files, err := os.ReadDir(pwd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading directory: %v", err)
			return [][]byte{}
		}

		for _, f := range files {
			if strings.HasPrefix(f.Name(), input) {
				if !f.IsDir() {
					matches = append(matches, []byte(f.Name()+" "))
				} else {
					matches = append(matches, []byte(f.Name()+string(os.PathSeparator)))
				}
			}
		}
	}

	if len(matches) == 0 {
		if file, _ := os.ReadDir(pwd); len(file) == 1 {
			return singleMatchHelper(file[0].Name(), input)
		}

		return [][]byte{}
	} else if len(matches) == 1 {
		return singleMatchHelper(string(matches[0]), input)
	}

	return matches
}

func runCompleteScript(path string) string {
	cmd := exec.Command(path)

	var buf bytes.Buffer
	tempWriter := bufio.NewWriter(&buf)

	cmd.Stdout = tempWriter
	cmd.Stderr = tempWriter

	err := cmd.Run()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(buf.String()) + " "
}

func singleMatchHelper(singleMatch, input string) [][]byte {
	match := make([][]byte, 0, 1)
	match = append(match, []byte(singleMatch[len(input):]))

	return match
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
