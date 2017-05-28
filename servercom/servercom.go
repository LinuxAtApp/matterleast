package serverCom

import (
	"fmt"
	mm "github.com/mattermost/platform/model"
)

type ServerCom struct {
	Client mm.Client
	team mm.Team
	channel *mm.Channel
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
	sc.team = *team.Data.(*mm.Team)
	sc.Client.SetTeamId(sc.team.Id)
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
	teamListResult, teamListAppError := sc.Client.GetAllTeamListings()
	teamMap := teamListResult.Data.(map[string]*mm.Team)
	if teamListAppError != nil {
		fmt.Println(teamListAppError)
		return 
	}
	fmt.Println("Teams:")
	for name, value := range teamMap {
		fmt.Println("\tName:", value.Name, "TeamID:", name)
	}
}

func (sc *ServerCom) SetChannel(channelName string) *mm.AppError {
	channelResult, err := sc.Client.GetChannelByName(channelName)
	if err != nil {
		return err
	}
	sc.channel = channelResult.Data.(*mm.Channel)
	return nil
}

func (sc *ServerCom) GetChannels() *mm.ChannelList {
	channelResult, channelErr := sc.Client.GetChannels(sc.team.Etag())
	if channelErr != nil {
		fmt.Println(channelErr)
		return nil
	}
	return channelResult.Data.(*mm.ChannelList)
}

func (sc *ServerCom) PrintChannels() {
	channelResult, channelErr := sc.Client.GetChannels(sc.team.Etag())
	if channelErr != nil {
		fmt.Println(channelErr)
		return
	}
	channelMap := channelResult.Data.(*mm.ChannelList)
	fmt.Println("Channels:")
	for id, channel := range *channelMap {
		fmt.Print("\tChannelID: ", id, " ChannelName: ", channel.Name)
		fmt.Println()
	}
}

/*
NewPost creates and pushes a post to the channel in channelId.
*/
func (sc *ServerCom) NewPost(message string) {
	post := &mm.Post{}
	post.ChannelId = sc.channel.Id
	post.Message = message
	_, err := sc.Client.CreatePost(post)
	if err != nil {
		fmt.Println(err)
	}
	return 
}