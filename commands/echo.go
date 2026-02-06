package commands

import (
	"bufio"
	"fmt"
	"strings"
)

func (sh *Shell) Echo(cmd *Command) error {

	buf := bufio.NewWriter(cmd.Stdout)
	_, err := fmt.Fprint(buf, strings.Join(cmd.Args, " "), "\n")
	if err != nil {
		fmt.Fprint(cmd.Stderr, err)
		return err
	}

	if err := buf.Flush(); err != nil {
		fmt.Fprint(cmd.Stderr, err)
		return err
	}

	return nil
}
