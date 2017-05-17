package main

import (
	"flag"
	"fmt"
	mm "github.com/mattermost/platform/model"
)

/*
main
Usage: go run main.go -u <username> -p <password> <server-url> [team-name]
Authenticates your login information, then gives you your AuthToken.
If the team name is unentered or invalid main shows valid team names.
*/
func main() {
	//Adds a  little clarity to the display
	fmt.Println("---------------------------------------------------------")
	//Sets up login
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	url := flag.Arg(0)
	teamName := flag.Arg(1)
	client := mm.NewClient(url)
	_, err := client.Login(*username, *password)
	if err != nil {
		fmt.Println(err)
		return
	}
	//Gathers all availible teams in a map,
	teamListResult, teamListAppError := client.GetAllTeamListings()
	teamMap := teamListResult.Data.(map[string]*mm.Team)
	if teamListAppError != nil {
		fmt.Println(teamListAppError)
		return
	}
	if teamName == nil {
		//Validates input team name
		teamObjMap, teamError := client.GetTeamByName(teamName)
		if teamError != nil {
			fmt.Println(teamError)
			return
		}
	} else {
		//Prints availible teams
		fmt.Println("Teams:")
		for _, value := range teamMap {
			fmt.Println("\t", value.Name)
		}
	}
	//Creates team map that can be accessed without string key, then assigns team ID
	localTeamSlice := make([]*mm.Team, len(teamMap))
	i := 0
	for _, value := range teamMap {
		localTeamSlice[i] = value
		i++
	}
	client.SetTeamId(localTeamSlice[0].Id)
	//Gather map of channels availible
	channelResult, channelErr := client.GetChannels(teamObjMap.Etag)
	if channelErr != nil {
		fmt.Println("Channel Error")
		fmt.Println(channelErr)
		return
	}
	//List availible channels (direct messages appear as address string, still in progress)
	channelMap := channelResult.Data.(*mm.ChannelList)
	channelSlice := make([]*mm.Channel, len(*channelMap))
	fmt.Print("\nChannels:\n")
	index := 0
	for _, channel := range *channelMap {
		fmt.Print("\t", index, ": ")
		channelSlice[index] = channel
		fmt.Println(channelSlice[index].DisplayName)
		index++
	}
	//TownSquare Channel ID: "d5gpjz3k3fyd7fhzqrafrxg6zr"
	//Gets mm.PostList since begining of time (?)
	postSinceDateResult, postsErr := client.GetPostsSince("d5gpjz3k3fyd7fhzqrafrxg6zr", 0)
	if postsErr != nil {
		fmt.Println(postsErr)
	}
	//Extracts PostList Object
	postSinceDate := postSinceDateResult.Data.(*mm.PostList)
	for _, post := range postSinceDate.Posts {
		//Gets\Extracts username of each post.
		userResult, userErr := client.GetUser(post.UserId, client.Etag)
		if userErr != nil {
			fmt.Println(userErr)
		}
		//Prints username and message.
		user := userResult.Data.(*mm.User)
		fmt.Println(user.Username)
		fmt.Println("\t", post.Message, "\n")
	}

	//Adds a  little clarity to the display
	fmt.Println("---------------------------------------------------------")

}
