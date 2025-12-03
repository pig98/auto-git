package main

import (
	"flag"
	"fmt"
	"os"

	"auto-git/internal/config"
	"auto-git/internal/logger"
	"auto-git/internal/watcher"
)

var (
	version          = "dev"
	buildTime        = "unknown"
	gitCommit        = "unknown"
	showVersion      = flag.Bool("version", false, "Show version information")
	showVersionShort = flag.Bool("v", false, "Show version information")
)

func main() {
	flag.Parse()

	// Show version and exit
	if *showVersion || *showVersionShort {
		fmt.Printf("auto-git version %s\n", version)
		fmt.Printf("Build time: %s\n", buildTime)
		fmt.Printf("Git commit: %s\n", gitCommit)
		os.Exit(0)
	}

	// Load configuration from environment
	logger.SetLevelFromEnv()
	config.LoadFromEnv()

	// Get Git directories from environment
	gitDirs, err := config.GetGitDirs()
	if err != nil {
		logger.Fatal("%v", err)
	}

	logger.Info("Auto Git Sync service started")
	logger.Info("Quiet period: %d minutes (sync only after %d minutes without file changes)", config.QuietPeriodMinutes, config.QuietPeriodMinutes)
	logger.Info("Configured Git directories (%d):", len(gitDirs))
	for i, dir := range gitDirs {
		logger.Info("  [%d] %s", i+1, dir)
	}

	// Start watching each Git directory
	for _, gitDir := range gitDirs {
		go func(dir string) {
			if err := watcher.WatchGitDir(dir); err != nil {
				logger.Error("Failed to watch directory %s: %v", dir, err)
			}
		}(gitDir)
	}

	// Keep the program running
	select {}
}
