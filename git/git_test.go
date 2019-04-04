package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrentBranch(t *testing.T) {
	b := CurrentBranch()
	require.NotEqual(t, b, "")
}

func TestRemoteBranches(t *testing.T) {
	remoteBranches := RemoteBranches()
	require.NotEmpty(t, remoteBranches)
}
