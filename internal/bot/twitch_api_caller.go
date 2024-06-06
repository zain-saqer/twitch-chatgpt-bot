package bot

import (
	"context"
	"errors"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"github.com/zain-saqer/twitch-chatgpt/internal/twitch"
	"time"
)

type TwitchApiCaller struct {
	api        *twitch.API
	repository chat.Repository
}

func NewTwitchApiCaller(api *twitch.API, repository chat.Repository) *TwitchApiCaller {
	return &TwitchApiCaller{
		api:        api,
		repository: repository,
	}
}

func (a *TwitchApiCaller) GetUser(ctx context.Context, accessToken, username string) (*twitch.User, error) {
	return a.api.GetUser(ctx, accessToken, username)
}

func (a *TwitchApiCaller) GetCurrentUser(ctx context.Context, accessToken string) (*twitch.User, error) {
	return a.api.GetCurrentUser(ctx, accessToken)
}

func (a *TwitchApiCaller) SendMessage(ctx context.Context, user *chat.User, broadcasterId, message string) (*twitch.SendMessageResponse, error) {
	response, err := a.api.SendMessage(ctx, user, broadcasterId, message)
	if err == nil {
		return response, nil
	}
	if !errors.Is(err, twitch.ErrUnauthorized) {
		return nil, err
	}
	refreshTokenResponse, err := a.api.RefreshAccessToken(user.RefreshToken)
	if err != nil {
		return nil, err
	}
	user.AccessToken = refreshTokenResponse.AccessToken
	user.RefreshToken = refreshTokenResponse.RefreshToken
	user.ExpiresAt = time.Now().Add(time.Duration(refreshTokenResponse.ExpiresIn) * time.Second)
	response, err = a.api.SendMessage(ctx, user, broadcasterId, message)
	return response, err
}
