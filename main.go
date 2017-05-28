package main

import (
	"flag"
	"fmt"
	client "github.com/LinuxAtApp/matterleast/servercom"
)

/*
Main usage: `go run main.go -u <username> -p <password> <url>
Package demontrates simple login functions using the servercom package's methods.
*/
func main() {
	fmt.Println("matterleast")
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	url := flag.Arg(0)
	//Creates client and logs user in.
	serverCom := client.Startup(url, *username, *password)
	//Tests if login was successful.
	if serverCom.Connected() {
		fmt.Println("*Connection Successful!*")
	} else {
		fmt.Println("*You are not connected.*")
	}
	serverCom.PrintTeams()
	serverCom.SetTeam("linuxapp")
	serverCom.PrintChannels()
	serverCom.SetChannel("town-square")
}
