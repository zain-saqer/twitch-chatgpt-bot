package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var ErrUnauthorized = errors.New("invalid credentials")

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type RefreshAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type API struct {
	clientId string
	client   *http.Client
}

func NewApi(clientId string, client *http.Client) *API {
	return &API{clientId: clientId, client: client}
}

func (api *API) RefreshAccessToken(refreshToken string) (*RefreshAccessTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", api.clientId)
	data.Set("client_secret", refreshToken)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 401 {
		return nil, ErrUnauthorized
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("refresh access token: non ok response " + resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := &RefreshAccessTokenResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *API) GetCurrentUser(ctx context.Context, accessToken string) (*User, error) {
	return api.GetUser(ctx, accessToken, "")
}

func (api *API) GetUser(ctx context.Context, accessToken, username string) (*User, error) {
	endpointUrl, err := url.ParseRequestURI("https://api.twitch.tv/helix/users")
	if err != nil {
		return nil, err
	}
	if username != "" {
		q := endpointUrl.Query()
		q.Add("login", username)
		endpointUrl.RawQuery = q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, "GET", endpointUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", api.clientId)
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch: invalid status code: %d", resp.StatusCode)
	}
	userData, err := io.ReadAll(resp.Body)
	var users struct {
		Data []*User `json:"data"`
	}
	err = json.Unmarshal(userData, &users)
	if err != nil {
		return nil, err
	}
	if len(users.Data) == 0 {
		return nil, fmt.Errorf("twitch: no users found")
	}
	return users.Data[0], nil
}

type sendMessageRequest struct {
	BroadcasterId string `json:"broadcaster_id"`
	SenderId      string `json:"sender_id"`
	Message       string `json:"message"`
}

type SendMessageResponse struct {
	MessageId string `json:"message_id"`
	IsSent    bool   `json:"is_sent"`
}

func (api *API) SendMessage(ctx context.Context, user *chat.User, broadcasterId string, message string) (*SendMessageResponse, error) {
	sendMessageReq := &sendMessageRequest{
		BroadcasterId: broadcasterId,
		SenderId:      user.ID,
		Message:       message,
	}
	reqBodyStr, err := json.Marshal(sendMessageReq)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewReader(reqBodyStr)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.twitch.tv/helix/chat/messages", reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch: invalid status code: %d", resp.StatusCode)
	}
	respData, err := io.ReadAll(resp.Body)
	var sendMessageResponse struct {
		Data []*SendMessageResponse `json:"data"`
	}
	err = json.Unmarshal(respData, &sendMessageResponse)
	if err != nil {
		return nil, err
	}
	if len(sendMessageResponse.Data) == 0 {
		return nil, fmt.Errorf("twitch: no users found")
	}
	return sendMessageResponse.Data[0], nil
}
