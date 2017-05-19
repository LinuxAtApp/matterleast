package main

import (
	"flag"
	"fmt"
	mm "github.com/mattermost/platform/model"
	"time"
)

func printLine() {
	fmt.Println("---------------------------------------------------------")
}

/*
main
Usage: go run main.go -u <username> -p <password> <server-url> [team-name]
Authenticates your login information, then gives you your AuthToken.
If the team name is unentered or invalid main shows valid team names.
*/
func main() {
	//Adds a  little clarity to the display
	printLine()
	defer printLine()
	//Sets up login
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	url := flag.Arg(0)
	teamName := flag.Arg(1)
	channelId := "d5gpjz3k3fyd7fhzqrafrxg6zr"
	client := mm.NewClient(url)
	_, err := client.Login(*username, *password)
	if err != nil {
		fmt.Println(err)
		return
	}
	if teamName == "" {
		//Gathers all availible teams in a map,
		teamListResult, teamListAppError := client.GetAllTeamListings()
		teamMap := teamListResult.Data.(map[string]*mm.Team)
		if teamListAppError != nil {
			fmt.Println(teamListAppError)
			return
		}
		fmt.Println("Teams:")
		for name, value := range teamMap {
			fmt.Println("\tName:", value.Name, "TeamID:", name)
		}
		return
	}
	//Validates input team name
	selectedTeam, teamError := client.GetTeamByName(teamName)
	if teamError != nil {
		fmt.Println(teamError)
		return
	}

	//Use the team provided as an argument
	client.SetTeamId(selectedTeam.Data.(*mm.Team).Id)
	//Gather map of channels availible
	channelResult, channelErr := client.GetChannels(selectedTeam.Etag)
	if channelErr != nil {
		fmt.Println("Channel Error")
		fmt.Println(channelErr)
		return
	}
	//List availible channels (direct messages appear as address string, still in progress)
	channelMap := channelResult.Data.(*mm.ChannelList)
	channelSlice := make([]*mm.Channel, len(*channelMap))
	fmt.Println("Channels:")
	index := 0
	for _, channel := range *channelMap {
		fmt.Print("\t", index, ": ")
		channelSlice[index] = channel
		fmt.Print(channelSlice[index].DisplayName)
		if channelSlice[index].Id == channelId {
			fmt.Print("*")
		}
		index++
		fmt.Println()
	}
	//Add a little clarity
	printLine()
	//Makes a new post then adds it to the server
	newPost := makePost(client, channelId, "Ping")
	_, createPostErr := client.CreatePost(newPost)
	if createPostErr != nil {
		println(err)
	}
	//displays last four posts
	printLastFourPosts(client, channelId)

}

func printLastFourPosts(client *mm.Client, channelId string) {
	//TownSquare Channel ID: "d5gpjz3k3fyd7fhzqrafrxg6zr"
	//Gets mm.PostList since begining of time (?)
	postSinceDateResult, postsErr := client.GetPostsSince(channelId, 0)
	if postsErr != nil {
		fmt.Println(postsErr)
	}
	//Extracts PostList Object
	postSinceDate := postSinceDateResult.Data.(*mm.PostList)
	//Parses of 4 most recent messages in selected Channel
	for i := 0; i < 4; i++ {
		// PostList.Order contains keys to the order of the posts. The most recent post gets stored at position 0
		postKey := postSinceDate.Order[i]
		post := postSinceDate.Posts[postKey]
		//Gets\Extracts username of each post.
		userResult, userErr := client.GetUser(post.UserId, client.Etag)
		if userErr != nil {
			fmt.Println(userErr)
		}
		//Prints username and message.
		user := userResult.Data.(*mm.User)
		fmt.Println(user.Username, user.Id, time.Unix(post.UpdateAt, 0))
		fmt.Println("\t", post.Message)
	}
}

//makePost Returns the address of a new Post
func makePost(client *mm.Client, channelId string, message string) *mm.Post {
	post := &mm.Post{}
	post.ChannelId = channelId
	post.Message = message
	return post
}
