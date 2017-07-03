package main

import (
    "strings"
	"flag"
	"fmt"
	client "github.com/LinuxAtApp/matterleast/servercom"
	"github.com/mattermost/platform/model"
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
	fmt.Println("Welcome to matterleast!")
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
		fatal(err)
		fmt.Println("No team specified, use the -t flag with one of:")
                for _, team := range teams {
    			fmt.Println("\t",team.Name)
		}
		return
	}
	err = serverCom.SetTeam(*team)
	fatal(err)
	if *channel == "" {
    		channels, err := serverCom.GetChannels()
    		fatal(err)
    		fmt.Println("No channel specified, use the -c flag with one of:")
    		for _, channel := range *channels {
        		fmt.Println("\t",channel.Name)
    		}
    		return
	}
	err = serverCom.SetChannel(*channel)
	fatal(err)
	fmt.Println("Channel [", serverCom.Channel.DisplayName, "]:")
	for event := range serverCom.Events {
    		if event.Broadcast.ChannelId != serverCom.Channel.Id {
        		continue // ignore events in a different channel
    		} else if event.Event != model.WEBSOCKET_EVENT_POSTED {
        		continue // ignore events that aren't messages
    		}
    		post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
    		fmt.Printf("%s: %s\n", event.Data["sender_name"], post.Message)
	}
}
