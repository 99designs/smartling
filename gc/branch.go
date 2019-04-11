package gc

import (
	"log"

	"github.com/99designs/smartling/git"
)

func Branch() error {
	return gcBranch(git.CurrentBranch(), true)
}

func gcBranch(branch string, ignoreLatest bool) error {
	log.Printf("Delete smartling files referencing %q...\n", branch)

	// get all files referencing the branch
	// loop through and delete the file

	return nil
}
