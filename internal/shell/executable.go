package shell

import "golang.org/x/sys/unix"

func (sh *Shell) IsExecutable(path string) bool {
	return unix.Access(path, unix.X_OK) == nil
}
