package notify

import (
	"os/exec"
	"runtime"
)

//Display a notification.
func Display(header string, body string, isUrgent bool, iconPath string) {
	switch runtime.GOOS {
	case "darwin":
		//Display dialog on MacOS
        exec.Command("osascript", "-e", "display notification \""+body+"\" with title \""+header+"\"").Run()
	case "linux":
		if isUrgent {
			//Display on Linux if urgent(will open until manually closed)
			exec.Command("notify-send", "-i", iconPath, header, body, "-u", "critical").Run()
		} else {
			//Display on Linux if not urgent
			exec.Command("notify-send", "-i", iconPath, header, body).Run()
		}
	}
}
