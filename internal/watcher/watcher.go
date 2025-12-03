package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"auto-git/internal/config"
	"auto-git/internal/git"
	"auto-git/internal/logger"
)

var (
	quietPeriodMap = make(map[string]*time.Timer)
)

// WatchGitDir starts watching a Git directory for changes
func WatchGitDir(gitDir string) error {
	absPath, err := filepath.Abs(gitDir)
	if err != nil {
		return err
	}

	// Check if it's a Git repository
	if !git.IsGitRepo(absPath) {
		logger.Info("Warning: %s is not a Git repository, skipping", absPath)
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Recursively add directories
	if err := addDirRecursive(watcher, absPath); err != nil {
		return err
	}

	logger.Info("Started watching: %s", absPath)

	// If there are already local changes when the service starts, schedule a quiet-period sync
	if dirty, err := git.HasWorkingTreeChanges(absPath); err == nil && dirty {
		logger.Debug("Detected existing uncommitted changes at startup for %s, scheduling quiet-period sync", absPath)
		ScheduleQuietSync(absPath, "Startup dirty state", "")
	}

	// Handle events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			handleEvent(event, absPath, watcher)
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			logger.Error("Watch error: %v", err)
		}
	}
}

func addDirRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignore errors, continue walking
		}

		// Skip .git directories
		if strings.Contains(path, "/.git/") || strings.HasSuffix(path, "/.git") {
			return nil
		}

		// Only watch directories
		if info.IsDir() {
			if err := watcher.Add(path); err != nil {
				logger.Error("Failed to add watch: %s, %v", path, err)
			}
		}

		return nil
	})
}

func handleEvent(event fsnotify.Event, gitDir string, watcher *fsnotify.Watcher) {
	// Ignore .git directory changes
	if strings.Contains(event.Name, "/.git/") {
		return
	}

	// For Remove events, we can't check if the path exists, but we can still check .gitignore
	// For other events, check if the path is ignored by .gitignore
	relPath, err := filepath.Rel(gitDir, event.Name)
	if err == nil {
		// git check-ignore works even for removed paths (checks Git index)
		if git.IsIgnoredByGit(gitDir, relPath) {
			logger.Debug("Ignored by .gitignore, skipping change: %s (in %s)", relPath, gitDir)
			return
		}
	}

	// Ignore temporary files (only for non-Remove events, as removed files won't have these suffixes)
	if event.Op&fsnotify.Remove != fsnotify.Remove {
		if strings.HasSuffix(event.Name, "~") || strings.HasSuffix(event.Name, ".swp") {
			return
		}
	}

	// Handle directory creation: add new directories to watch list
	if event.Op&fsnotify.Create == fsnotify.Create {
		info, err := os.Stat(event.Name)
		if err == nil && info.IsDir() {
			// New directory created, add it to watch list recursively
			if err := addDirRecursive(watcher, event.Name); err != nil {
				logger.Error("Failed to add new directory to watch: %s, %v", event.Name, err)
			} else {
				logger.Debug("Added new directory to watch: %s", event.Name)
			}
		}
	}

	// Handle directory removal: fsnotify automatically removes deleted directories from watch list
	// We just need to trigger the sync, which will be handled by ScheduleQuietSync below
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		logger.Debug("Directory or file removed: %s", event.Name)
		// Note: fsnotify automatically removes deleted directories from the watch list,
		// so we don't need to manually remove them
	}

	// Use quiet-period mechanism: reset timer on each change (including Create, Write, Remove, Rename)
	ScheduleQuietSync(gitDir, "File change", event.Name)
}

// ScheduleQuietSync sets or resets the quiet-period timer for a Git directory.
// If file is empty, it is treated as a generic trigger (e.g., startup state).
func ScheduleQuietSync(gitDir, source, file string) {
	// Reset timer on each trigger, only sync when quiet period ends without new triggers
	if timer, exists := quietPeriodMap[gitDir]; exists {
		timer.Stop()
		if file != "" {
			logger.Debug("%s: resetting quiet period timer: %s (in %s)", source, file, gitDir)
		} else {
			logger.Debug("%s: resetting quiet period timer for %s", source, gitDir)
		}
	} else {
		if file != "" {
			logger.Debug("%s: starting quiet period timer: %s (in %s)", source, file, gitDir)
		} else {
			logger.Debug("%s: starting quiet period timer for %s", source, gitDir)
		}
	}

	timer := time.AfterFunc(time.Duration(config.QuietPeriodMinutes)*time.Minute, func() {
		logger.Info("Quiet period ended (%d minutes without new file changes), starting sync: %s", config.QuietPeriodMinutes, gitDir)
		git.SyncGit(gitDir)
		delete(quietPeriodMap, gitDir)
	})

	quietPeriodMap[gitDir] = timer
}
