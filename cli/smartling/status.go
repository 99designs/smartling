package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/tabwriter"

	"github.com/99designs/smartling"
)

func MustStatus(remotefile, locale string) smartling.FileStatus {
	fs, err := client.Status(remotefile, locale)
	logAndQuitIfError(err)

	return fs
}

type ProjectStatus map[string]map[string]smartling.FileStatus

func (statuses ProjectStatus) Add(remotefile, locale string, fs smartling.FileStatus) {
	_, ok := statuses[remotefile]
	if !ok {
		mm := make(map[string]smartling.FileStatus)
		statuses[remotefile] = mm
	}
	statuses[remotefile][locale] = fs
}

func (statuses ProjectStatus) AwaitingAuthorizationCount() int {
	c := 0
	for _, s := range statuses {
		for _, status := range s {
			c += status.AwaitingAuthorizationStringCount()
			break
		}
	}

	return c
}

func (statuses ProjectStatus) TotalStringsCount() int {
	c := 0
	for _, s := range statuses {
		for _, status := range s {
			c += status.StringCount
			break
		}
	}

	return c
}

func GetProjectStatus(prefix string, locales []string) ProjectStatus {
	var wg sync.WaitGroup
	statuses := ProjectStatus{}
	remoteFiles := fetchRemoteFileList()

	for _, projectFilepath := range ProjectConfig.Files() {
		prefixedProjectFilepath := filepath.Clean("/" + prefix + "/" + projectFilepath)
		if !remoteFiles.contains(prefixedProjectFilepath) {
			prefixedProjectFilepath = filepath.Clean("/" + projectFilepath)
		}

		for _, l := range locales {
			wg.Add(1)
			go func(remotefile, locale string) {
				defer wg.Done()
				statuses.Add(remotefile, locale, MustStatus(remotefile, locale))
			}(prefixedProjectFilepath, l)
		}
	}
	wg.Wait()

	return statuses
}

func PrintProjectStatusTable(statuses ProjectStatus, locales []string) {
	// Format in columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprint(w, "Awaiting\t  In Progress -> Completed\n")
	fmt.Fprint(w, "    Auth\t\n")
	for _, locale := range locales {
		fmt.Fprint(w, "\t  ", locale)
	}
	fmt.Fprint(w, "\n")

	for projectFilepath, _ := range statuses {
		aa := false
		for _, locale := range locales {
			status := statuses[projectFilepath][locale]
			if !aa {
				fmt.Fprintf(w, "%7d", status.AwaitingAuthorizationStringCount())
				aa = true
			}
			fmt.Fprintf(w, "\t%3d->%-3d", status.InProgressStringCount(), status.CompletedStringCount)
		}
		fmt.Fprint(w, "\t", projectFilepath, "\n")
	}
	w.Flush()
}
