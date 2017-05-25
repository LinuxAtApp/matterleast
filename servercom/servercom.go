package servercom

import (
	mm "github.com/mattermost/platform/model"
	"fmt"
)
/*
Startup accepts the url and login credentials for a user, and returns a pointer to a logged in client.
*/
func Startup(url string, username string, password string) *mm.Client {
	client := mm.NewClient(url)
	_, err := client.Login(username, password)
	if err != nil {
		fmt.Println(err)
	}
	return client
}
/*
Connected accepts a mattermost client as an arguement, and returns true if the client has a
login Authentication Token.
*/
func Connected(client mm.Client) bool {
	if client.AuthToken != "" {
		return true
	}
	return false
}
