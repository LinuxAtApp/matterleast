package servercom_test

import (
	"fmt"
	mm "github.com/mattermost/platform/model"
	"io"
	"github.com/LinuxAtApp/matterleast/servercom"
)

// mockClient fakes a connection to the mattermost server and allows
// tests to control the return values of all mattermost client functions.
type mockClient struct {
	NextResult *mm.Result
	NextErr    *mm.AppError
}

var _ servercom.Client = &mockClient{}

func (m *mockClient) Login(loginId string, password string) (*mm.Result, *mm.AppError) {
	return m.NextResult, m.NextErr
}

func (m *mockClient) GetTeamByName(teamName string) (*mm.Result, *mm.AppError) {
	return m.NextResult, m.NextErr
}

func (m *mockClient) SetTeamId(teamId string) {
	return m.NextResult, m.NextErr
}

func (m *mockClient) GetAllTeamListings() (*mm.Result, *mm.AppError) {
	return m.NextResult, m.NextErr
}

func (m *mockClient) GetChannelByName(channelName string) (*mm.Result, *mm.AppError) {
	return m.NextResult, m.NextErr
}

func (m *mockClient) GetChannels(etag string) (*mm.Result, *mm.AppError) {
	return m.NextResult, m.NextErr
}

func (m *mockClient) CreatePost(post *mm.Post) (*mm.Result, *mm.AppError) {
	return m.NextResult, m.NextErr
}

func (m *mockClient) GetPostsSince(channelId string, time int64) (*mm.Result, *mm.AppError) {
	return m.NextResult, m.NextErr
}
