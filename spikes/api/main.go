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

	fmt.Println("Auth successful! Token: ", client.AuthToken)

	//Gathers all availible teams in a map,
	teamListResult, teamListAppError := client.GetAllTeamListings()
	teamMap := teamListResult.Data.(map[string]*mm.Team)
	if teamListAppError != nil {
		fmt.Println(teamListAppError)
		return
	}
	//Validates input team name
	teamObjMap, teamError := client.GetTeamByName(teamName)
	if teamError != nil {
		fmt.Println(teamError)
		return
	}
	//Prints availible teams
	fmt.Println("teams:")
	for _, value := range teamMap {
		fmt.Println("\t", value.Name)
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
	channelSlice := channelResult.Data.(*mm.ChannelList)
	fmt.Print("\nChannels:\n")
	index := 0
	for _, channel := range *channelSlice {
		fmt.Print("\t", index, ": ")
		fmt.Println(channel.DisplayName)
		index++
	}

}
