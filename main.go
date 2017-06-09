package main

import (
	"flag"
	"fmt"
	client "github.com/LinuxAtApp/matterleast/servercom"
	"os"
)

//fatal is a simple error handler that prints an error if it is not nil
//then continues with the program just to clean up code.
func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
Main usage: `go run main.go -u <username> -p <password> <url>
Package demontrates simple login functions using the servercom package's methods.
*/
func main() {
	fmt.Println("matterleast")
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	team := flag.String("t", "", "Team Name")
	channel := flag.String("c", "", "Channel Name")
	flag.Parse()
	url := flag.Arg(0)
	//if any login data is missing main.go cleanly exits.
	if url == "" || *username == "" || *password == "" {
		return
	}
	//Creates client and logs user in.
	serverCom, err := client.Startup(url, *username, *password, os.Stdout)
	fatal(err)
	//Tests if login was successful.
	if serverCom.Connected() {
		fmt.Println("*Connection Successful!*")
	} else {
		fmt.Println("*You are not connected.*")
	}
	if *team == "" {
		serverCom.PrintTeams()
		return
	}
	err = serverCom.SetTeam(*team)
	fatal(err)
	if *channel == "" {
		serverCom.PrintChannels()
		return
	}
	err = serverCom.SetChannel(*channel)
	fatal(err)
	fmt.Println("Channel [", serverCom.Channel.DisplayName, "] data:\n")
	currentPost := 0
	numPosts := 1
	for {
        	posts, err := serverCom.GetSelectPosts(currentPost, numPosts)
        	fatal(err)
        	for _, post := range posts {
        		fmt.Println(post.Message, "\n")
        	}
        	currentPost++
	}
}
