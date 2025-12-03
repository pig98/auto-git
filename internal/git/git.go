package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"auto-git/internal/logger"
	"auto-git/internal/notify"
)

// IsGitRepo checks if a path is a Git repository
func IsGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

// IsIgnoredByGit checks whether a path is ignored by git (.gitignore, global excludes, etc.)
func IsIgnoredByGit(gitDir, relPath string) bool {
	cmd := exec.Command("git", "check-ignore", relPath)
	cmd.Dir = gitDir
	if err := cmd.Run(); err != nil {
		// Non-zero exit code means not ignored (or git version without check-ignore),
		// in either case we treat as not ignored to avoid missing real changes.
		return false
	}
	return true
}

// HasWorkingTreeChanges checks if there are any unstaged or uncommitted changes
// in the working tree (equivalent to checking `git status --porcelain`).
func HasWorkingTreeChanges(gitDir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = gitDir
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

// SyncGit performs the complete git sync operation
func SyncGit(gitDir string) {
	logger.Info("Starting sync: %s", gitDir)

	// 1. Stage local changes first so that working tree is clean before pulling
	if err := gitAddAll(gitDir); err != nil {
		logger.Error("Failed to add files: %v", err)
		notify.Send("Git Sync Failed", fmt.Sprintf("Failed to add files: %v", err), "Alert")
		return
	}

	// 2. Check if there is anything to commit
	hasChanges, err := hasUncommittedChanges(gitDir)
	if err != nil {
		logger.Error("Failed to check changes: %v", err)
		return
	}

	if !hasChanges {
		logger.Info("No local changes to sync, skipping pull and push")
		return
	}

	// 3. Commit local changes
	if err := gitCommit(gitDir); err != nil {
		logger.Error("Commit failed: %v", err)
		notify.Send("Git Sync Failed", fmt.Sprintf("Commit failed: %v", err), "Alert")
		return
	}
	logger.Info("Local changes committed successfully")

	// 4. Pull latest changes (rebasing our local commits on top)
	if err := gitPull(gitDir); err != nil {
		logger.Error("Pull failed: %v", err)
		notify.Send("Git Sync Failed", fmt.Sprintf("Pull failed: %v", err), "Alert")
		return
	}

	// 5. Check and resolve conflicts after pull/rebase
	hasConflict, err := checkConflict(gitDir)
	if err != nil {
		logger.Error("Failed to check conflicts: %v", err)
		return
	}

	if hasConflict {
		logger.Info("Conflict detected, attempting to resolve...")
		if err := resolveConflict(gitDir); err != nil {
			logger.Error("Failed to resolve conflict: %v", err)
			notify.Send("Git Conflict Unresolved", fmt.Sprintf("Directory: %s\nError: %v", gitDir, err), "Alert")
			return
		}
	}

	// 6. Push
	if err := gitPush(gitDir); err != nil {
		logger.Error("Push failed: %v", err)
		notify.Send("Git Sync Failed", fmt.Sprintf("Push failed: %v", err), "Alert")
		return
	}

	logger.Info("Sync successful: %s", gitDir)
	notify.Send("Git Sync Successful", fmt.Sprintf("Directory: %s", gitDir), "Info")
}
