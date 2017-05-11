package main

import (
	"flag"
	"fmt"
	mm "github.com/mattermost/platform/model"
)
/*
main Run main -u <username> -p <password> <server URL> 
Authenticates your login information, giving you your AuthToken, then lists the Team Listings.
*/
func main() {
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	url := flag.Arg(0)
	client := mm.NewClient(url)
	_, err := client.Login(*username, *password)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, appError := client.GetAllTeamListings()
	//result.(map[string]string)
	teamMap := result.Data.(map[string]*mm.Team)
	fmt.Println("Auth successful! Token: ", client.AuthToken)
	if appError != nil {
		fmt.Println(appError)
		return
	}
	fmt.Println("teams:")

	for _, value := range teamMap {
		fmt.Println("\t", value.Name)
	}
}
