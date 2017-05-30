package main

import (
	"github.com/LinuxAtApp/matterleast/spikes/notify/notificator"
)

var notify *notificator.Notificator

func main() {
	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "My test App",
	})

	notify.Push("title", "text", true, "/home/user/icon.png", notificator.UR_NORMAL)
}
