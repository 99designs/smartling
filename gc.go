package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/99designs/api-sdk-go"
	"github.com/99designs/smartling/git"
)

func GcBranch(client *FaultTolerantClient, dryRun bool) error {
	return gcBranch(client, git.CurrentBranch(), true, dryRun)
}

func GcProject(client *FaultTolerantClient, dryRun bool) error {
	smartlingBranches := SmartlingBranches(client)
	branches := git.RemoteBranches()
	mergedBranches := git.MergedRemoteBranches()

	for _, smartlingBranch := range smartlingBranches {
		isDeleted := branches[smartlingBranch] == ""
		isMerged := mergedBranches[smartlingBranch] != ""

		if isDeleted || isMerged {
			err := gcBranch(client, smartlingBranch, false, dryRun)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func gcBranch(client *FaultTolerantClient, branch string, ignoreLatest bool, dryRun bool) error {
	log.Printf("Delete smartling files referencing %q...\n", branch)

	files := FilesPrefixedWithGitBranch(client, branch)
	latestFiles := latestBranchFiles(files.Items)
	for _, file := range files.Items {
		filePath := filePathFromBranchURI(file.FileURI)
		if filePath == "" {
			continue
		}
		if ignoreLatest && latestFiles[filePath].FileURI == file.FileURI {
			continue
		}

		fmt.Printf("Delete %s uploaded at %s\n", file.FileURI, file.LastUploaded)

		if dryRun {
			continue
		}
		err := client.Delete(file.FileURI)
		if err != nil {
			return err
		}
	}

	return nil
}

func latestBranchFiles(branchFiles []smartling.File) map[string]smartling.File {
	latestFiles := map[string]smartling.File{}

	for _, file := range branchFiles {
		filePath := filePathFromBranchURI(file.FileURI)
		if filePath == "" {
			continue
		}

		latestFile, found := latestFiles[filePath]
		if !found || file.LastUploaded.After(latestFile.LastUploaded.Time) {
			latestFiles[filePath] = file
		}
	}

	return latestFiles
}

func filePathFromBranchURI(uri string) string {
	var branchRegex = regexp.MustCompile(`^/branch/[^/]+/[^/]+/(.*)$`)
	matches := branchRegex.FindAllStringSubmatch(uri, 1)
	if len(matches) != 1 || len(matches[0]) != 2 {
		return ""
	}

	return matches[0][1]
}
