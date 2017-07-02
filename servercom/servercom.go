package servercom

import (
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

// WSClient is a websocket client to the mattermost server
type WSClient interface {
    Close()
    Listen()
    Events() chan *mm.WebSocketEvent
}

// wsWrapper wraps an existing websocket client implementation so that the event
// channel can be accessed via a method (rather than as a field). This allow us
// to hide the WSClient completely behind an interface and swap the implementation
// in the future if need be.
type wsWrapper struct {
    *mm.WebSocketClient
}

// Events returns the stream of realtime events from the mattermost server.
func (w *wsWrapper) Events() chan *mm.WebSocketEvent {
    return w.WebSocketClient.EventChannel
}

// test that wsWrapper implements WSClient at compile-time
var _ WSClient = &wsWrapper{}

// test that *mm.Client satisfies the Client interface at compile time
var _ Client = &mm.Client{}

/*
ServerCom acts as a mitigator between the frontend and the mattermost model API.
*/
type ServerCom struct {
	Client  Client
	Team    mm.Team
	Channel *mm.Channel
	Events  chan *mm.WebSocketEvent
	Responses chan *mm.WebSocketResponse
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
