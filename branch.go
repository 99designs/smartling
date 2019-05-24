package main

import (
	"strings"

	"github.com/99designs/api-sdk-go"
)

func FilesPrefixedWithGitBranch(client *FaultTolerantClient, branch string) *smartling.FilesList {
	URIMask := "/branch/"
	if branch != "" {
		URIMask = URIMask + branch + "/"
	}

	req := smartling.FilesListRequest{
		URIMask: URIMask,
	}

	files, err := client.List(req)
	if err != nil {
		panic(err)
	}
	return files
}

func SmartlingBranches(client *FaultTolerantClient) map[string]string {
	files := FilesPrefixedWithGitBranch(client, "")

	branches := map[string]string{}
	for _, f := range files.Items {
		uriParts := strings.Split(f.FileURI, "/")
		branch := uriParts[2]

		// Some apps are pushing using the master branch as production
		// We want to exclude it from the smartling branches
		if branch == "master" {
			continue
		}
		branches[branch] = branch
	}

	return branches
}
