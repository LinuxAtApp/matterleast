package serverCom

import (
	"fmt"
	mm "github.com/mattermost/platform/model"
)
/*
ServerCom acts as a mitigator between the frontend and the mattermost model API.
*/
type ServerCom struct {
	Client mm.Client
	Team mm.Team
	Channel *mm.Channel
}

/*
Startup accepts the url and login credentials for a user, and returns a new serverCom struct.
*/
func Startup(url string, username string, password string) ServerCom {
	ServerCom := ServerCom{Client: *mm.NewClient(url)}
	_, err := ServerCom.Client.Login(username, password)
	if err != nil {
		fmt.Println(err)
	}
	return ServerCom
}

/*
Connected returns true if the client has a login Authentication Token.
*/
func (sc *ServerCom) Connected() bool {
	if sc.Client.AuthToken != "" {
		return true
	}
	return false
}

/*
SetTeam set ServerCom's client team ID using the team name.
*/
func (sc *ServerCom) SetTeam(teamName string) {
	team, err := sc.Client.GetTeamByName(teamName)
	if err != nil {
		fmt.Println(err)
		return
	}
	sc.Team = *team.Data.(*mm.Team)
	sc.Client.SetTeamId(sc.Team.Id)
	return 
}

func (sc *ServerCom) GetTeams() map[string]*mm.Team {
	teamListResult, teamListAppError := sc.Client.GetAllTeamListings()
	teamMap := teamListResult.Data.(map[string]*mm.Team)
	if teamListAppError != nil {
		fmt.Println(teamListAppError)
		return nil
	}
	return teamMap
}

func (sc *ServerCom) PrintTeams() {
	fmt.Println("Teams:")
	for name, value := range sc.GetTeams() {
		fmt.Println("\tName:", value.Name, "TeamID:", name)
	}
}

func (sc *ServerCom) SetChannel(channelName string) {
	channelResult, err := sc.Client.GetChannelByName(channelName)
	if err != nil {
		fmt.Println(err)
	}
	sc.Channel = channelResult.Data.(*mm.Channel)
	return 
}

//GetChannels returns a map with all availible channels in team.
//**Must have successfully run SetTeam()**
func (sc *ServerCom) GetChannels() *mm.ChannelList {
	channelResult, channelErr := sc.Client.GetChannels(sc.Team.Etag())
	if channelErr != nil {
		fmt.Println(channelErr)
		return nil
	}
	return channelResult.Data.(*mm.ChannelList)
}

func (sc *ServerCom) PrintChannels() {
	fmt.Println("Channels:")
	for id, channel := range *sc.GetChannels() {
		fmt.Print("\tChannelID: ", id, " ChannelName: ", channel.Name)
		fmt.Println()
	}
}
//GetChannelData returns a slice containing every post in a channel.
//The posts are in reverse order, so the oldest post is at the [0] index.
func (sc *ServerCom) GetChannelData() []*mm.Post {
	//TownSquare Channel ID: "d5gpjz3k3fyd7fhzqrafrxg6zr"
	//Gets mm.PostList since begining of time
	postSinceDateResult, postsErr := sc.Client.GetPostsSince(sc.Channel.Id, 0)
	if postsErr != nil {
		fmt.Println(postsErr)
	}
	//Extracts Posts in order into postsSLice
	postSinceDate := postSinceDateResult.Data.(*mm.PostList)
	postsSlice := make([]*mm.Post, len(postSinceDate.Order))
	for index, key := range postSinceDate.Order {
		postsSlice[index] = postSinceDate.Posts[key]
	}
	return postsSlice
}

//GetSelectPosts returns a slice selection of posts.
//Offset (int) is how many posts back the newest post in the slice will be.
//Postcount is the number of posts before (and including) the offset that will be in the slice.
func (sc *ServerCom) GetSelectPosts(offset int, postCount int) []*mm.Post{
	postList := sc.GetChannelData()
	selectPosts := make([]*mm.Post, postCount)
	for i := 0; i < postCount; i++ {
		selectPosts[i] = postList[i + offset]
	}
	return selectPosts
}
	
/*
NewPost creates and pushes a post to the channel in channelId.
*/
func (sc *ServerCom) NewPost(message string) {
	post := &mm.Post{}
	post.ChannelId = sc.Channel.Id
	post.Message = message
	_, err := sc.Client.CreatePost(post)
	if err != nil {
		fmt.Println(err)
	}
	return 
}
