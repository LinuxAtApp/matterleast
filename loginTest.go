package main

import (
	"flag"
	"fmt"
	"github.com/LinuxAtApp/matterleast/servercom"
)
func main() {
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	url := flag.Arg(0)
	client:= *servercom.Startup(url, *username, *password)
	if servercom.Connected(client) {
		fmt.Println("*Connection Successful!*")
	} else {
		fmt.Println("*You are not connected.*")
	}
}