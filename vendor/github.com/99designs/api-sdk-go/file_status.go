package smartling

import (
	"fmt"
)

// FileStatus describes file translation status obtained by GetFileStatus
// method.
type FileStatus struct {
	File

	TotalStringCount int
	TotalWordCount   int
	TotalCount       int

	Items []FileStatusTranslation
}

func (fs FileStatus) GetFileStatusTranslation(locale string) (*FileStatusTranslation, error) {
	for i := range fs.Items {
		if fs.Items[i].LocaleID == locale {
			return &fs.Items[i], nil
		}
	}

	return nil, fmt.Errorf(
		"failed to get file status translation for locale: %s", locale,
	)
}

func (fs FileStatus) AwaitingAuthorizationStringCount() int {
	c := 0

	for i := range fs.Items {
		c += fs.Items[i].AwaitingAuthorizationStringCount(fs.TotalStringCount)
	}

	return c
}
