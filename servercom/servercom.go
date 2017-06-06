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
func Startup(url string, username string, password string) (*ServerCom, *mm.AppError) {
	ServerCom := ServerCom{Client: *mm.NewClient(url)}
	_, err := ServerCom.Client.Login(username, password)
	if err != nil {
		return nil, err
	}
	return ServerCom, nil
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
func (sc *ServerCom) SetTeam(teamName string) *mm.AppError {
	team, err := sc.Client.GetTeamByName(teamName)
	if err != nil {
		return err
	}
	sc.Team = *team.Data.(*mm.Team)
	sc.Client.SetTeamId(sc.Team.Id)
	return nil
}

func (sc *ServerCom) GetTeams() (map[string]*mm.Team, *mm.AppError) {
	teamListResult, teamListAppError := sc.Client.GetAllTeamListings()
	teamMap := teamListResult.Data.(map[string]*mm.Team)
	if teamListAppError != nil {
		return nil, teamListAppError
	}
	return teamMap, nil
}

func (sc *ServerCom) PrintTeams() {
	fmt.Println("Teams:")
	for name, value := range sc.GetTeams() {
		fmt.Println("\tName:", value.Name, "TeamID:", name)
	}
}

func (sc *ServerCom) SetChannel(channelName string) *mm.AppError {
	channelResult, err := sc.Client.GetChannelByName(channelName)
	if err != nil {
		return err
	}
	sc.Channel = channelResult.Data.(*mm.Channel)
	return nil
}

//GetChannels returns a map with all availible channels in team.
//**Must have successfully run SetTeam()**
func (sc *ServerCom) GetChannels() (*mm.ChannelList, *mm.AppError) {
	channelResult, channelErr := sc.Client.GetChannels(sc.Team.Etag())
	if channelErr != nil {
		return nil, channelErr
	}
	return channelResult.Data.(*mm.ChannelList), nil
}

func (sc *ServerCom) PrintChannels() {
	fmt.Println("Channels:")
	for id, channel := range *sc.GetChannels() {
		fmt.Print("\tChannelID: ", id, " ChannelName: ", channel.Name)
		fmt.Println()
	}
}
//GetChannelData returns a slice containing every post in a channel.
//The posts are in order, so the newest post is at the [0] index.
func (sc *ServerCom) GetChannelData() ([]*mm.Post, *mm.AppError) {
	//TownSquare Channel ID: "d5gpjz3k3fyd7fhzqrafrxg6zr"
	//Gets mm.PostList since begining of time
	postSinceDateResult, postsErr := sc.Client.GetPostsSince(sc.Channel.Id, 0)
	if postsErr != nil {
		return nil, postsErr
	}
	//Extracts Posts in order into postsSLice
	postSinceDate := postSinceDateResult.Data.(*mm.PostList)
	postsSlice := make([]*mm.Post, len(postSinceDate.Order))
	for index, key := range postSinceDate.Order {
		postsSlice[index] = postSinceDate.Posts[key]
	}
	return postsSlice, nil
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
func (sc *ServerCom) NewPost(message string) *mm.AppError {
	post := &mm.Post{}
	post.ChannelId = sc.Channel.Id
	post.Message = message
	_, err := sc.Client.CreatePost(post)
	if err != nil {
		return err
	}
	return nil
}
