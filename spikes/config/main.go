package main

import (
	"fmt"
	"os"
	"path"
)

const BUFFER_SIZE = 1024

// fatal crashes the program if the given error is non-nil
// This isn't a good way to perform production error-handling,
// but it will serve for this demo.
func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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
	configFile, err := os.Open(path.Join(xdgConfigPath, "matterleast.conf"))
	fatal(err)

	defer configFile.Close()
	data := make([]byte, BUFFER_SIZE)
	bytesRead, err := configFile.Read(data)
	fatal(err)

}
