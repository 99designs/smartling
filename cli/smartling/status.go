package main

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	smartlingNew "github.com/Smartling/api-sdk-go"
)

func mustStatus(remotefile, locale string) smartlingNew.FileStatusExtended {
	fs, err := client.Status(remotefile, locale)
	logAndQuitIfError(err)

	return *fs
}

type ProjectStatus struct {
	sync.RWMutex
	statuses map[string]map[string]smartlingNew.FileStatusExtended
}

func New() *ProjectStatus {
	return &ProjectStatus{
		statuses: make(map[string]map[string]smartlingNew.FileStatusExtended),
	}
}

func (ps *ProjectStatus) Add(remotefile, locale string, fs smartlingNew.FileStatusExtended) {
	ps.Lock()
	defer ps.Unlock()

	_, ok := ps.statuses[remotefile]
	if !ok {
		mm := make(map[string]smartlingNew.FileStatusExtended)
		ps.statuses[remotefile] = mm
	}
	ps.statuses[remotefile][locale] = fs
}

func (ps *ProjectStatus) AwaitingAuthorizationCount() int {
	ps.RLock()
	defer ps.RUnlock()

	c := 0
	for _, s := range ps.statuses {
		for _, status := range s {
			c += status.AwaitingAuthorizationStringCount()
			break
		}
	}
	return c
}

func (ps *ProjectStatus) TotalStringsCount() int {
	ps.RLock()
	defer ps.RUnlock()

	c := 0
	for _, s := range ps.statuses {
		for _, status := range s {
			c += status.TotalStringCount
			break
		}
	}

	return c
}

func GetProjectStatus(prefix string, locales []string) *ProjectStatus {
	var wg sync.WaitGroup
	statuses := New()

	for _, projectFilepath := range ProjectConfig.Files() {
		remoteFilePath := findIdenticalRemoteFileOrPush(projectFilepath, prefix)

		for _, l := range locales {
			wg.Add(1)
			go func(remotefile, locale string) {
				defer wg.Done()
				statuses.Add(remotefile, locale, mustStatus(remotefile, locale))
			}(remoteFilePath, l)
		}
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
			status := ps.statuses[projectFilepath][locale]
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
