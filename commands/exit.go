package commands

import "os"

func handleExit(cmdName, redirection string, inputHistory *[]History, args ...string) {
	os.Exit(0)
}
