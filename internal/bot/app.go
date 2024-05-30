package bot

import (
	"context"
	twitchirc "github.com/gempir/go-twitch-irc/v4"
	"github.com/google/uuid"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"github.com/zain-saqer/twitch-chatgpt/internal/irc"
)

type App struct {
	Repository    chat.Repository
	TwitchClient  *twitchirc.Client
	WhitelistByID map[uuid.UUID]*chat.Username
	Whitelist     map[string]*chat.Username
}

func (a *App) JoinChannel(channel ...string) {
	a.TwitchClient.Join(channel...)
}

func (a *App) Depart(channel string) {
	a.TwitchClient.Depart(channel)
}

func (a *App) AddUsername(username *chat.Username) {
	a.Whitelist[username.Name] = username
	a.WhitelistByID[username.ID] = username
}

func (a *App) RemoveUsername(username *chat.Username) {
	delete(a.Whitelist, username.Name)
	delete(a.WhitelistByID, username.ID)
}

func (a *App) StartMessagePipeline(ctx context.Context) error {
	usernames, err := a.Repository.GetUsernames(ctx)
	if err != nil {
		return err
	}
	for _, username := range usernames {
		a.AddUsername(username)
	}
	channels, err := a.Repository.GetChannels(ctx)
	if err != nil {
		return err
	}
	for _, channel := range channels {
		a.JoinChannel(channel.Name)
	}
	messageTypes := []uint8{chat.PrivMsg}
	messageStream, err := irc.NewMessagePipeline(a.TwitchClient)(ctx, messageTypes)
	if err != nil {
		return err
	}
	filteredMessageStream := chat.FilterMessageStream(ctx, messageStream, messageTypes)
	chat.ServeMessageStream(ctx, filteredMessageStream, func() map[string]*chat.Username { return a.Whitelist })
	return nil
}
