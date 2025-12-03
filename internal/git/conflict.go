package git

import (
	"fmt"
	"os/exec"
	"strings"

	"auto-git/internal/logger"
)

// checkConflict checks if there are any merge conflicts
func checkConflict(gitDir string) (bool, error) {
	// Check Git status
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = gitDir
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	status := string(output)
	// Check for conflict markers (UU=both modified, AA=both added, DD=both deleted)
	if strings.Contains(status, "UU") || strings.Contains(status, "AA") || strings.Contains(status, "DD") {
		return true, nil
	}

	// Check if in rebase/merge state
	cmd = exec.Command("git", "status")
	cmd.Dir = gitDir
	statusOutput, err := cmd.Output()
	if err == nil {
		statusStr := string(statusOutput)
		if strings.Contains(statusStr, "rebase in progress") ||
			strings.Contains(statusStr, "merge conflict") ||
			strings.Contains(statusStr, "Unmerged paths") {
			return true, nil
		}
	}

	// Check for conflict markers in files (<<<<<<, =======, >>>>>>>)
	cmd = exec.Command("git", "diff", "--check")
	cmd.Dir = gitDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		// Output usually means there are conflicts
		if len(output) > 0 {
			return true, nil
		}
	}

	// Use grep to find conflict markers
	cmd = exec.Command("grep", "-r", "-l", "<<<<<<<", ".")
	cmd.Dir = gitDir
	output, _ = cmd.Output()
	if len(output) > 0 {
		return true, nil
	}

	return false, nil
}

// resolveConflict attempts to automatically resolve merge conflicts
func resolveConflict(gitDir string) error {
	// Check if in rebase state
	cmd := exec.Command("git", "status")
	cmd.Dir = gitDir
	statusOutput, _ := cmd.Output()
	statusStr := string(statusOutput)

	isRebasing := strings.Contains(statusStr, "rebase") || strings.Contains(statusStr, "REBASE")

	// Strategy 1: Try using 'ours' strategy (keep local changes)
	logger.Info("Attempting to resolve conflict using 'ours' strategy (keep local changes)...")
	cmd = exec.Command("git", "checkout", "--ours", ".")
	cmd.Dir = gitDir
	if err := cmd.Run(); err == nil {
		// Ours strategy succeeded, add files and continue
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		if err := cmd.Run(); err == nil {
			if isRebasing {
				cmd = exec.Command("git", "rebase", "--continue")
				cmd.Dir = gitDir
				output, err := cmd.CombinedOutput()
				if err == nil {
					logger.Info("Conflict resolved successfully using 'ours' strategy")
					return nil
				}
				logger.Error("Rebase continue failed: %s", string(output))
			} else {
				logger.Info("Conflict resolved successfully using 'ours' strategy")
				return nil
			}
		}
	}

	// Strategy 2: If ours failed, try using 'theirs' strategy (use remote changes)
	logger.Info("Attempting to resolve conflict using 'theirs' strategy (use remote changes)...")
	cmd = exec.Command("git", "checkout", "--theirs", ".")
	cmd.Dir = gitDir
	if err := cmd.Run(); err == nil {
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		if err := cmd.Run(); err == nil {
			if isRebasing {
				cmd = exec.Command("git", "rebase", "--continue")
				cmd.Dir = gitDir
				output, err := cmd.CombinedOutput()
				if err == nil {
					logger.Info("Conflict resolved successfully using 'theirs' strategy")
					return nil
				}
				logger.Error("Rebase continue failed: %s", string(output))
			} else {
				logger.Info("Conflict resolved successfully using 'theirs' strategy")
				return nil
			}
		}
	}

	// If both failed, try to abort rebase
	if isRebasing {
		logger.Error("Conflict resolution failed, attempting to abort rebase...")
		abortCmd := exec.Command("git", "rebase", "--abort")
		abortCmd.Dir = gitDir
		abortCmd.Run()
	}

	return fmt.Errorf("unable to automatically resolve conflict, manual intervention required")
}
