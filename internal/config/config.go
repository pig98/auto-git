package config

import (
	"os"
	"strconv"
	"strings"
)

var (
	QuietPeriodMinutes  = 10
	NotificationsEnable = true
)

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() {
	// Load quiet period
	if quietPeriodEnv := os.Getenv("QUIET_PERIOD_MINUTES"); quietPeriodEnv != "" {
		if minutes, err := strconv.Atoi(quietPeriodEnv); err == nil && minutes > 0 {
			QuietPeriodMinutes = minutes
		}
	}

	// Load notification setting
	v := strings.ToLower(strings.TrimSpace(os.Getenv("DISABLE_NOTIFICATIONS")))
	if v == "1" || v == "true" || v == "yes" {
		NotificationsEnable = false
	}
}

// GetGitDirs parses GIT_DIRS environment variable and returns valid directories
func GetGitDirs() ([]string, error) {
	gitDirsEnv := os.Getenv("GIT_DIRS")
	if gitDirsEnv == "" {
		return nil, ErrGitDirsNotSet
	}

	gitDirs := strings.Split(gitDirsEnv, ":")
	validDirs := make([]string, 0)
	for _, gitDir := range gitDirs {
		gitDir = strings.TrimSpace(gitDir)
		if gitDir != "" {
			validDirs = append(validDirs, gitDir)
		}
	}

	if len(validDirs) == 0 {
		return nil, ErrNoValidGitDirs
	}

	return validDirs, nil
}
