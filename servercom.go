package servercom

import (
	mm "github.com/mattermost/platform/model"
	"fmt"
)

func Startup(url string, username string, password string) *mm.Client {
	client := mm.NewClient(url)
	_, err := client.Login(*username, *password)
	if err != nil {
		fmt.Println(err)
	}
	return &client
}

func Connected(client mm.Client) boolean {
	if client.AuthToken != nil {
		return true
	}
	return false
}