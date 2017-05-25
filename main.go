package main

import (
	"flag"
	"fmt"
	"github.com/LinuxAtApp/matterleast/servercom"
)
/*
Main usage: `go run main.go -u <username> -p <password> <url>
Package demontrates simple login functions using the servercom package's methods.
*/
func main() {
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	url := flag.Arg(0)
	//Creates client and logs user in.
	client:= *servercom.Startup(url, *username, *password)
	//Tests if login was successful. 
	if servercom.Connected(client) {
		fmt.Println("*Connection Successful!*")
	} else {
		fmt.Println("*You are not connected.*")
	}
}