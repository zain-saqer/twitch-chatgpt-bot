package bot

import (
	"context"
	twitchirc "github.com/gempir/go-twitch-irc/v4"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"github.com/zain-saqer/twitch-chatgpt/internal/chatgpt"
	"github.com/zain-saqer/twitch-chatgpt/internal/twitch"
	"sync"
)

type App struct {
	Repository     chat.Repository
	TwitchClient   *twitchirc.Client
	lock           sync.Mutex
	Users          map[string]*chat.User
	ChannelsByUser map[string]map[string]*chat.Channel
	TwitterAPI     *TwitchApiCaller
	ChatGPTAPI     *chatgpt.API
}

func (a *App) JoinChannel(channel ...string) {
	a.TwitchClient.Join(channel...)
}

func (a *App) Depart(channel string) {
	a.TwitchClient.Depart(channel)
}

func (a *App) AddUser(user *chat.User) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, ok := a.Users[user.Username]; ok {
		return
	}
	a.Users[user.Username] = user
	a.ChannelsByUser[user.Username] = make(map[string]*chat.Channel)
}

func (a *App) RemoveUser(user *chat.User) {
	a.lock.Lock()
	defer a.lock.Unlock()
	delete(a.Users, user.Username)
	delete(a.ChannelsByUser, user.Username)
}

func (a *App) AddChannel(user *chat.User, channel *chat.Channel) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.ChannelsByUser[user.Username][channel.Name] = channel
	a.JoinChannel(channel.Name)
}

func (a *App) RemoveChannel(user *chat.User, channel *chat.Channel) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, ok := a.ChannelsByUser[user.Username]; !ok {
		return
	}
	delete(a.ChannelsByUser[user.Username], channel.Name)
	a.Depart(channel.Name)
}

func (a *App) findUser(username string) *chat.User {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.Users[username]
}

func (a *App) findChannel(user *chat.User, channelName string) *chat.Channel {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, ok := a.ChannelsByUser[user.Username]; !ok {
		return nil
	}
	return a.ChannelsByUser[user.Username][channelName]
}

func (a *App) sendTwitchMessage(ctx context.Context, user *chat.User, channel *chat.Channel, message string) error {
	_, err := a.TwitterAPI.SendMessage(ctx, user, channel.ID, message)
	return err
}

func (a *App) gpt(ctx context.Context, query string) (string, error) {
	answer, err := a.ChatGPTAPI.Completions(ctx, query)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func (a *App) StartMessagePipeline(ctx context.Context) error {
	users, err := a.Repository.GetUsers(ctx)
	if err != nil {
		return err
	}
	var channels []string
	for _, user := range users {
		a.AddUser(user)
		userChannels, err := a.Repository.GetChannelsByUser(ctx, user.ID)
		if err != nil {
			return err
		}
		for _, channel := range userChannels {
			channels = append(channels, channel.Name)
			a.AddChannel(user, channel)
		}
	}
	messageTypes := []uint8{chat.PrivMsg}
	messageStream, err := twitch.NewMessagePipeline(a.TwitchClient)(ctx, messageTypes)
	if err != nil {
		return err
	}
	filteredMessageStream := chat.FilterMessageStream(ctx, messageStream, messageTypes)
	chat.ServeMessageStream(ctx, filteredMessageStream, a.findUser, a.findChannel, a.sendTwitchMessage, a.gpt)
	return nil
}
