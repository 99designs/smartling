package gc

import (
	"github.com/99designs/smartling/git"
)

func Project() error {
	smartlingBranches := []string{
		"add-dockerfile",
	}

	branches := git.RemoteBranches()
	mergedBranches := git.MergedRemoteBranches()

	for _, smartlingBranch := range smartlingBranches {
		isDeleted := branches[smartlingBranch] == ""
		isMerged := mergedBranches[smartlingBranch] != ""

		if isDeleted || isMerged {
			err := gcBranch(smartlingBranch, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
