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
	serverCom, err := client.Startup(url, *username, *password)
	fatal(err)
	if *team == "" {
		teams, err := serverCom.GetTeams()
                for _, team := range teams {
    			fmt.Println(team)
		}
		fatal(err)
	}
	err = serverCom.SetTeam(*team)
	fatal(err)
	if *channel == "" {
    		channels, err := serverCom.GetChannels()
    		for _, channel := range *channels {
        		fmt.Println(channel)
    		}
    		fatal(err)
	}
	err = serverCom.SetChannel(*channel)
	fatal(err)
	fmt.Println("Channel [", serverCom.Channel.DisplayName, "] data:\n")
	for event := range serverCom.Events {
    		fmt.Println(event)
	}
}
