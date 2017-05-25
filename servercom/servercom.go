package servercom

import (
	mm "github.com/mattermost/platform/model"
	"fmt"
)

func Startup(url string, username string, password string) *mm.Client {
	client := mm.NewClient(url)
	fmt.Println("URL:", url)
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
	fmt.Println(*client)
	_, err := client.Login(username, password)
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func Connected(client mm.Client) bool {
	if client.AuthToken != "" {
		return true
	}
	return false
}
