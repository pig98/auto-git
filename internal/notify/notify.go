package notify

import (
	"fmt"
	"os/exec"

	"auto-git/internal/config"
)

// Send sends a macOS system notification
func Send(title, message, level string) {
	if !config.NotificationsEnable {
		return
	}
	// Use macOS osascript to send notification
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", script)
	cmd.Run()
}
