package main

import (
	"flag"
	"fmt"
	mm "github.com/mattermost/platform/model"
)
/*
main -u <username> -p <password> <server URL> [team name]
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
	clientResult, clientAppError := client.GetAllTeamListings()
	teamMap := clientResult.Data.(map[string]*mm.Team)
	if clientAppError != nil {
		fmt.Println(clientAppError)
		return
	}
	//Validates input team name
	teamObjMap, teamError := client.GetTeamByName(teamName)
	if teamError != nil {
		fmt.Println( teamError )
		return
	}
	//Prints availible teams
	fmt.Println("teams:")
	for _, value := range teamMap {
		fmt.Println("\t", value.Name)
	}
	//Creates team map that can be accessed without string key, then assigns team ID
	lclTeamMap := make(map[int]*mm.Team)
	i := 0
	for _, value := range teamMap {
		lclTeamMap[i] = value
	}
	client.SetTeamId( lclTeamMap[0].Id )
	//Gather map of channels availible
	chnlResult,chnlErr := client.GetChannels( teamObjMap.Etag )
	if chnlErr != nil {
		fmt.Println( "Channel Error" )
		fmt.Println ( chnlErr )
		return
	}
	//List availible channels (direct messages appear as address string, still in progress)
	chnlSlice := chnlResult.Data.(*mm.ChannelList)
	fmt.Print( "\nChannels:\n" )
	for _, channel := range *chnlSlice {
		fmt.Println("\t", channel.Name)
	}
	
}
