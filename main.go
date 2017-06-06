package main

import (
	"flag"
	"fmt"
	client "github.com/LinuxAtApp/matterleast/servercom"
)

//handle is a simple error handler that prints an error if it is not nil
//then continues with the program just to clean up code.
func handle(err *mm.AppError) {
	if err != nil {
		fmt.Println(err)
	}
	return
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
	//if any login data is missing main.go cleanly exits.
	if (url == "" || *username == "" || *password == "") {
		return
	}
	//Creates client and logs user in.
	serverCom, := client.Startup(url, *username, *password)
	//Tests if login was successful.
	if serverCom.Connected() {
		fmt.Println("*Connection Successful!*")
	} else {
		fmt.Println("*You are not connected.*")
	}
	serverCom.PrintTeams()
	handle( err := serverCom.SetTeam("linuxapp"))
	serverCom.PrintChannels()
	handle( err = serverCom.SetChannel("town-square"))
	fmt.Println("Channel [", serverCom.Channel.DisplayName, "] data:\n")
	for _, post := range serverCom.GetSelectPosts(5,3) {
		fmt.Println(post.Message,"\n")
	}
}

