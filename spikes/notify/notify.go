package notify

import (
	"os/exec"
	"runtime"
)

//Display a notification.
func Display(header string, body string, isUrgent bool, iconPath string) {
	switch runtime.GOSS {
	case "darwin":
		//Display dialog on MacOS
		exec.Command("osascript", "-e", "display dialog \""+body+"\" with title \""+summary+"\"").Run()
	case "linux":
		if isUrgent {
			//Display on Linux if urgent(will open until manually closed)
			exec.Command("notify-send", "-i", iconPath, summary, body, "-u", "critical").Run()
		} else {
			//Display on Linux if not urgent
			exec.Command("notify-send", "-i", iconPath, summary, body).Run()
		}
	}
}
