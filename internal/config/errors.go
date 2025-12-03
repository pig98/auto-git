package config

import "errors"

var (
	ErrGitDirsNotSet  = errors.New("environment variable GIT_DIRS is not set, please configure it in plist file")
	ErrNoValidGitDirs = errors.New("no valid Git directories configured")
)
