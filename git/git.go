package git

import (
	"bytes"
	"os/exec"
	"strings"
)

func CurrentBranch() string {
	return gitExec("symbolic-ref", "--short", "HEAD")
}

func Branch(args ...string) string {
	return gitExecSubCommand("branch", args)
}

func MergedRemoteBranches() []string {
	return RemoteBranches("--merged", "master")
}

func RemoteBranches(args ...string) []string {
	defaultArgs := []string{"-r"}
	allRemoteBranches := Branch(append(defaultArgs, args...)...)
	branches := strings.Split(allRemoteBranches, "\n")

	var result []string
	for _, branchPath := range branches {
		if strings.Contains(branchPath, "HEAD") || strings.Contains(branchPath, "master") {
			continue
		}

		trimmed := strings.TrimSpace(branchPath)
		branch := strings.Replace(trimmed, "origin/", "", 1)
		if branch != "" {
			result = append(result, branch)
		}
	}
	return result
}

func gitExecSubCommand(name string, args []string) string {
	return gitExec(append([]string{name}, args...)...)
}

func gitExec(args ...string) string {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()

	return strings.TrimSpace(out.String())
}
