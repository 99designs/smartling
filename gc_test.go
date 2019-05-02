package main

import (
	"testing"
	"time"

	"github.com/99designs/api-sdk-go"

	"github.com/stretchr/testify/require"
)

func TestGc(t *testing.T) {
	t.Run("filePathFromBranchURI", func(t *testing.T) {
		testCases := map[string]string{
			"/branch/gc-smartling-branches/9489628/test-config/locales/en.yml":  "test-config/locales/en.yml",
			"/invalid/gc-smartling-branches/9489628/test-config/locales/en.yml": "",
		}

		for subject, expected := range testCases {
			require.Equal(t, expected, filePathFromBranchURI(subject))
		}
	})

	t.Run("latestBranchFiles", func(t *testing.T) {
		latest := time.Now()
		latestFile1 := smartling.File{FileURI: "/branch/gc-smartling-branches/latest/test-config/1/1.yml", LastUploaded: smartling.UTC{latest}}
		latestFile2 := smartling.File{FileURI: "/branch/gc-smartling-branches/latest/test-config/2/2.yml", LastUploaded: smartling.UTC{latest}}

		old := time.Now().AddDate(0, 0, -1)
		old1 := smartling.File{FileURI: "/branch/gc-smartling-branches/old/test-config/1/1.yml", LastUploaded: smartling.UTC{old}}
		old2 := smartling.File{FileURI: "/branch/gc-smartling-branches/old/test-config/2/2.yml", LastUploaded: smartling.UTC{old}}

		branchFiles := []smartling.File{
			latestFile1,
			old1,
			old2,
			latestFile2,
		}

		require.Equal(t, map[string]smartling.File{
			"test-config/1/1.yml": latestFile1,
			"test-config/2/2.yml": latestFile2,
		}, latestBranchFiles(branchFiles))
	})
}
