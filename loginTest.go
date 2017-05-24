package main

import (
	"flag"
	"fmt"
	"servercom"
)
func main() {
	username := flag.String("u","","Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	url := flag.Arg(0)
	client:= servercom.Startup(url, username, password)
	if servercom.Connected {
		fmt.Println("Connection Successful!")
	}
}