package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type getUserResponse struct {
	Data []*User `json:"data"`
}

type API struct {
	accessToken  string
	refreshToken string
	clientId     string
	client       *http.Client
}

func NewApi(accessToken, refreshToken, clientId string, client *http.Client) *API {
	return &API{accessToken: accessToken, refreshToken: refreshToken, clientId: clientId, client: client}
}

func (api *API) do(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+api.accessToken)
	r.Header.Set("Client-Id", api.clientId)
	r.Header.Set("Content-Type", "application/json")
	return api.client.Do(r)
}

func (api *API) GetUser() (*User, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, err
	}
	resp, err := api.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch: invalid status code: %d", resp.StatusCode)
	}
	userData, err := io.ReadAll(resp.Body)
	var users getUserResponse
	err = json.Unmarshal(userData, &users)
	if err != nil {
		return nil, err
	}
	if len(users.Data) == 0 {
		return nil, fmt.Errorf("twitch: no users found")
	}
	return users.Data[0], nil
}
