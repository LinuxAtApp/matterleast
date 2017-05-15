package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	xdgConfigPath := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigPath == "" {
		home := os.Getenv("HOME")
		if home == "" {
			fmt.Println("$HOME undefined, aborting...")
			os.Exit(1)
		}
		xdgConfigPath = path.Join(home, ".config")
	}
	fmt.Println("Config Dir:", xdgConfigPath)
}
