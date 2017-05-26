package serverCom

import (
	"fmt"
	mm "github.com/mattermost/platform/model"
)

type ServerCom struct {
	Client *mm.Client
}

/*
Startup accepts the url and login credentials for a user, and returns a new serverCom struct.
*/
func Startup(url string, username string, password string) ServerCom {
	ServerCom := ServerCom{mm.NewClient(url)}
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
