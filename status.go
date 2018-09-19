package main

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/Smartling/api-sdk-go"
)

func mustStatus(remotefile string) smartling.FileStatus {
	fs, err := client.Status(remotefile)
	logAndQuitIfError(err)

	return *fs
}

type ProjectStatus struct {
	sync.RWMutex
	statuses map[string]smartling.FileStatus
}

func New() *ProjectStatus {
	return &ProjectStatus{
		statuses: make(map[string]smartling.FileStatus),
	}
}

func (ps *ProjectStatus) AwaitingAuthorizationCount() int {
	c := 0
	for _, s := range ps.statuses {
		c += s.AwaitingAuthorizationStringCount()
	}
	return c
}

func (ps *ProjectStatus) TotalStringsCount() int {
	c := 0
	for _, s := range ps.statuses {
		c += s.TotalStringCount
	}

	return c
}

func GetProjectStatus(prefix string, locales []string) *ProjectStatus {
	var wg sync.WaitGroup
	statuses := New()

	for _, projectFilepath := range ProjectConfig.Files() {
		remoteFilePath := findIdenticalRemoteFileOrPush(projectFilepath, prefix)

		wg.Add(1)
		go func(remoteFile string) {
			defer wg.Done()
			fs := mustStatus(remoteFile)
			statuses.statuses[remoteFile] = fs
		}(remoteFilePath)
	}
	wg.Wait()

	return statuses
}

func PrintProjectStatusTable(ps *ProjectStatus, locales []string) {
	// Format in columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprint(w, "Awaiting\t  In Progress -> Completed\n")
	fmt.Fprint(w, "    Auth\t\n")
	for _, locale := range locales {
		fmt.Fprint(w, "\t  ", locale)
	}
	fmt.Fprint(w, "\n")

	for projectFilepath, _ := range ps.statuses {
		aa := false
		for _, locale := range locales {
			status := ps.statuses[projectFilepath]
			fst, err := status.GetFileStatusTranslation(locale)
			if err != nil {
				logAndQuitIfError(err)
			}
			if !aa {
				fmt.Fprintf(w, "%7d", fst.AwaitingAuthorizationStringCount(status.TotalStringCount))
				aa = true
			}

			fmt.Fprintf(w, "\t%3d->%-3d", fst.AuthorizedStringCount, fst.CompletedStringCount)
		}
		fmt.Fprint(w, "\t", projectFilepath, "\n")
	}
	w.Flush()
}
