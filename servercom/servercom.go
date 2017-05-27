package serverCom

import (
	"fmt"
	mm "github.com/mattermost/platform/model"
)

type ServerCom struct {
	Client    *mm.Client
	channelId string
}

/*
Startup accepts the url and login credentials for a user, and returns a new serverCom struct.
*/
func Startup(url string, username string, password string) ServerCom {
	ServerCom := ServerCom{Client: mm.NewClient(url)}
	_, err := ServerCom.Client.Login(username, password)
	if err != nil {
		fmt.Println(err)
	}
	return ServerCom
}

/*
Connected returns true if the client has a login Authentication Token.
*/
func (sc ServerCom) Connected() bool {
	if sc.Client.AuthToken != "" {
		return true
	}
	return false
}

/*
SetTeam set ServerCom's client team ID using the team name.
*/
func (sc ServerCom) SetTeam(teamName string) *mm.AppError {
	team, err := sc.Client.GetTeamByName(teamName)
	if err != nil {
		return err
	}
	sc.Client.SetTeamId(team.Data.(*mm.Team).Id)
	return nil
}

/*
NewPost creates and pushes a post to the channel in channelId.
*/
func (sc ServerCom) NewPost(message string) *mm.AppError {
	post := &mm.Post{}
	post.ChannelId = sc.channelId
	post.Message = message
	_, err := sc.Client.CreatePost(post)
	if err != nil {
		return (err)
	}
	return nil
}
