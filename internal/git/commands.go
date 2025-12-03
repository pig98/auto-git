package git

import (
	"fmt"
	"os/exec"
	"time"
)

func gitPull(gitDir string) error {
	cmd := exec.Command("git", "pull", "--rebase")
	cmd.Dir = gitDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

func gitAddAll(gitDir string) error {
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = gitDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

func hasUncommittedChanges(gitDir string) (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = gitDir
	err := cmd.Run()
	if err != nil {
		// Exit code 1 means there are changes
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func gitCommit(gitDir string) error {
	commitMsg := fmt.Sprintf("Auto commit: %s", time.Now().Format("2006-01-02 15:04:05"))
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = gitDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

func gitPush(gitDir string) error {
	cmd := exec.Command("git", "push")
	cmd.Dir = gitDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}
