package main

import (
	"flag"
	"fmt"
	client "github.com/LinuxAtApp/matterleast/servercom"
)

//fatal is a simple error handler that prints an error if it is not nil
//then continues with the program just to clean up code.
func fatal(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}
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
	serverCom, err := client.Startup(url, *username, *password)
	fatal(err)
	//Tests if login was successful.
	if serverCom.Connected() {
		fmt.Println("*Connection Successful!*")
	} else {
		fmt.Println("*You are not connected.*")
	}
	serverCom.PrintTeams()
	err = serverCom.SetTeam("linuxapp")
	fatal(err)
	serverCom.PrintChannels()
	err = serverCom.SetChannel("town-square")
	fatal(err)
	fmt.Println("Channel [", serverCom.Channel.DisplayName, "] data:\n")
	posts, err := serverCom.GetSelectPosts(5,3)
	for _, post := range posts {
		fmt.Println(post.Message,"\n")
	}
}

