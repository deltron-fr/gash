package commands

import (
	"fmt"
	"slices"
	"strings"
)

func (sh *Shell) JobsCmd(cmd *Command) error {
	length := len(sh.BackgroundJobs)
	for i, job := range sh.BackgroundJobs {
		cmdString := strings.Join(job.BackgroundProcess.Args, " ")
		marker := " "
		if i == length-2 {
			marker = "-"
		}

		if i == length-1 {
			marker = "+"
		}

		fmt.Fprintf(cmd.Stdout, "[%d]%s  %-24s%s\n", job.ID, marker, job.Status.String(), cmdString)
	}

	sh.BackgroundJobs = slices.DeleteFunc(sh.BackgroundJobs, func(job *BackgroundJob) bool {
		return job.Status == StatusDone
	})

	return nil
}
