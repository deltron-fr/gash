package input

import (
	"bufio"
)

func handleKeys(reader *bufio.Reader) string {
	seq := make([]byte, 2)
	var err error

	seq[0], err = reader.ReadByte()
	if err != nil {
		return ""
	}

	seq[1], err = reader.ReadByte()
	if err != nil {
		return ""
	}

	if seq[0] == '[' {
		switch seq[1] {
		case 'A':
			return "Up"
		case 'B':
			return "Down"
		case 'C':
			return "Right"
		case 'D':
			return "Left"
		}
	}

	return ""
}
