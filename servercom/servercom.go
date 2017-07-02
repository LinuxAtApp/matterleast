package servercom

import (
	"fmt"
	mm "github.com/mattermost/platform/model"
)

type Client interface {
    Login(loginId string, password string) (*mm.Result, *mm.AppError)
    GetTeamByName(teamName string) (*mm.Result, *mm.AppError)
    SetTeamId(teamId string)
    GetAllTeamListings() (*mm.Result, *mm.AppError)
    GetChannelByName(channelName string) (*mm.Result, *mm.AppError)
    GetChannels(etag string) (*mm.Result, *mm.AppError)
    CreatePost(post *mm.Post) (*mm.Result, *mm.AppError)
    GetPostsSince(channelId string, time int64) (*mm.Result, *mm.AppError)
}

// test that *mm.Client satisfies the Client interface at compile time
var _ Client = &mm.Client{}

/*
ServerCom acts as a mitigator between the frontend and the mattermost model API.
*/
type ServerCom struct {
	Client  Client
	Team    mm.Team
	Channel *mm.Channel
}

/*
Startup accepts the url and login credentials for a user, and returns a new serverCom struct.
*/
func Startup(url string, username string, password string) (*ServerCom, error) {
	ServerCom := &ServerCom{Client: mm.NewClient(url)}
	_, err := ServerCom.Client.Login(username, password)
	if err != nil {
		return nil, err
	}
	return ServerCom, nil
}

/*
SetTeam set ServerCom's client team ID using the team name.
*/
func (sc *ServerCom) SetTeam(teamName string) error {
	team, err := sc.Client.GetTeamByName(teamName)
	if err != nil {
		return err
	}
	sc.Team = *team.Data.(*mm.Team)
	sc.Client.SetTeamId(sc.Team.Id)
	return nil
}

func (sc *ServerCom) GetTeams() (map[string]*mm.Team, error) {
	teamListResult, teamListAppError := sc.Client.GetAllTeamListings()
	teamMap := teamListResult.Data.(map[string]*mm.Team)
	if teamListAppError != nil {
		return nil, teamListAppError
	}
	return teamMap, nil
}

func (sc *ServerCom) SetChannel(channelName string) error {
	channelResult, err := sc.Client.GetChannelByName(channelName)
	if err != nil {
		return err
	}
	sc.Channel = channelResult.Data.(*mm.Channel)
	return nil
}

//GetChannels returns a map with all availible channels in team.
//**Must have successfully run SetTeam()**
func (sc *ServerCom) GetChannels() (*mm.ChannelList, error) {
	channelResult, channelErr := sc.Client.GetChannels(sc.Team.Etag())
	if channelErr != nil {
		return nil, channelErr
	}
	return channelResult.Data.(*mm.ChannelList), nil
}

//GetChannelData returns a slice containing every post in a channel.
//The posts are in order, so the newest post is at the [0] index.
func (sc *ServerCom) GetChannelData() ([]*mm.Post, error) {
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
func (sc *ServerCom) GetSelectPosts(offset int, postCount int) ([]*mm.Post, error) {
	postList, err := sc.GetChannelData()
	if err != nil {
		return nil, err
	}
	if offset > len(postList)  {
    		return nil, fmt.Errorf("Index out of bounds, len(postList)=%d, offset=%d, postCount=%d", len(postList), offset, postCount)
	}
	if offset + postCount > len(postList) {
        	postCount = len(postList) - offset
	}
	return postList[offset:offset+postCount], nil
}

/*
NewPost creates and pushes a post to the channel in channelId.
*/
func (sc *ServerCom) NewPost(message string) error {
	post := &mm.Post{}
	post.ChannelId = sc.Channel.Id
	post.Message = message
	_, err := sc.Client.CreatePost(post)
	if err != nil {
		return err
	}
	return nil
}
