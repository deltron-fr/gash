package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/deltron-fr/gash/internal/commands"
	"golang.org/x/term"
)

func RawModeHandler(currentBuffer string, history []commands.History) (string, []string) {
	// RawModeHandler enters the terminal's raw mode and reads user input
	// byte-by-byte. It supports basic line editing (left/right/delete),
	// tab completion, and returns once the user presses Enter.
	//
	// Parameters:
	//   currentBuffer - an optional string that will be prefilled into
	//                   the input buffer before reading further keys.
	//	 history - contains the history of previous commands for up and
	//             down arrow navigation.
	//
	// Returns:
	//   string - the full input line the user typed (without trailing newline)
	//   []string - when tab completion listing is requested this slice
	//              contains the matches to be printed by the caller.
	//
	// - This function is synchronous and blocks until a terminating key
	//   (Enter) is pressed.
	// - The function temporarily puts the TTY into raw mode and restores
	//   the previous state on return.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	r := bufio.NewReader(os.Stdin)

	var buffer []byte
	var cursorPos int
	var tabPressed bool

	var historyTracker int
	if len(history) > 0 {
		historyTracker = len(history)
	}

	for {
		if currentBuffer != "" {
			fmt.Fprint(os.Stdout, currentBuffer)
			cursorPos += len(currentBuffer)
			buffer = append(buffer, []byte(currentBuffer)...)
			currentBuffer = ""
			continue
		}

		b, err := r.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			break
		}

		if b == 0x0D {
			fmt.Fprintf(os.Stdout, "\r\n")
			return string(buffer), []string{}
		}

		if b <= 0x1f || b == 0x7f {
			switch b {
			case 0x1b:
				tabPressed = false
				key := handleKeys(r)
				switch key {
				case "":
					continue
				case "Left":
					if cursorPos > 0 {
						fmt.Fprintf(os.Stdout, "\x1b[D")
						cursorPos--
					}
				case "Right":
					if cursorPos < len(buffer) {
						fmt.Fprintf(os.Stdout, "\x1b[C")
						cursorPos++
					}
				case "Up":
					if len(history) == 0 {
						continue
					}

					if historyTracker <= 0 {
						continue
					}

					for range buffer {
						fmt.Fprintf(os.Stdout, "\x1b[D")
						fmt.Fprintf(os.Stdout, " ")
						fmt.Fprintf(os.Stdout, "\x1b[D")
						cursorPos--
					}

					historyTracker--
					buffer = []byte{}
					fmt.Print(history[historyTracker].Name)
					cursorPos = 0
					cursorPos += len(history[historyTracker].Name)
					buffer = append(buffer, []byte(history[historyTracker].Name)...)
				case "Down":
					if len(history) == 0 {
						continue
					}

					if historyTracker >= len(history)-1 {
						if len(buffer) != 0 {
							for range buffer {
								fmt.Fprintf(os.Stdout, "\x1b[D")
								fmt.Fprintf(os.Stdout, " ")
								fmt.Fprintf(os.Stdout, "\x1b[D")
								cursorPos--
							}
							buffer = []byte{}
						}
						continue
					}

					for range buffer {
						fmt.Fprintf(os.Stdout, "\x1b[D")
						fmt.Fprintf(os.Stdout, " ")
						fmt.Fprintf(os.Stdout, "\x1b[D")
						cursorPos--
					}

					historyTracker++
					buffer = []byte{}
					fmt.Print(history[historyTracker].Name)
					cursorPos = 0
					cursorPos += len(history[historyTracker].Name)
					buffer = append(buffer, []byte(history[historyTracker].Name)...)
				}
			case 0x0A, 0x0C:
				fmt.Fprintf(os.Stdout, "\r\n")
				return string(buffer), []string{}

			case 0x7f, 0x08:
				tabPressed = false
				if len(buffer) == 0 {
					continue
				}
				fmt.Fprintf(os.Stdout, "\x1b[D")
				fmt.Fprintf(os.Stdout, " ")
				fmt.Fprintf(os.Stdout, "\x1b[D")
				cursorPos--
				buffer = buffer[:len(buffer)-1]

			case 0x09:
				if len(buffer) == 0 {
					fmt.Fprintf(os.Stdout, "\x07")
					continue
				}

				parts := strings.Split(string(buffer), " ")
				hasCommand := len(parts) > 1
				targetInput := parts[len(parts)-1]
				restOfInput := autoCompletion(targetInput, hasCommand)
				if len(restOfInput) == 0 {
					fmt.Fprintf(os.Stdout, "\x07")
					continue
				}

				if len(restOfInput) > 1 {
					matches := buildMatches(restOfInput)

					if tabPressed {
						fmt.Fprintf(os.Stdout, "\r\n")
						return string(buffer), matches
					}

					lcp := checkLongestCommonPrefix(matches)
					if len(lcp) > len(targetInput) && lcp != targetInput {
						fmt.Fprint(os.Stdout, lcp[len(targetInput):])
						cursorPos += len(lcp[len(targetInput):])
						buffer = append(buffer, []byte(lcp[len(targetInput):])...)
					}

					tabPressed = true
					fmt.Fprintf(os.Stdout, "\x07")
					continue
				}

				for _, b := range restOfInput[0] {
					fmt.Fprintf(os.Stdout, "%c", b)
					buffer = append(buffer, b)
					cursorPos++
				}
			}
		} else {
			if cursorPos == len(buffer) {
				fmt.Fprintf(os.Stdout, "%c", b)
				buffer = append(buffer, b)
				cursorPos++
			} else {
				buffer = append(buffer, 0)
				copy(buffer[cursorPos+1:], buffer[cursorPos:len(buffer)-1])
				buffer[cursorPos] = b

				for i := cursorPos; i < len(buffer); i++ {
					fmt.Fprintf(os.Stdout, "%c", buffer[i])
				}
				for i := 0; i < len(buffer[cursorPos:])-1; i++ {
					fmt.Fprintf(os.Stdout, "\x1b[D")
				}
			}
		}
	}
	return "", []string{}
}
