package chat

import (
	"context"
	"github.com/rs/zerolog/log"
	"slices"
	"strings"
)

type GetMessageStream func(ctx context.Context, messageTypes []uint8) (<-chan *Message, error)

func FilterMessageStream(ctx context.Context, messageStream <-chan *Message, allowedTypes []uint8) <-chan *Message {
	filteredMessageStream := make(chan *Message)

	go func() {
		defer close(filteredMessageStream)
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-messageStream:
				if !ok {
					return
				}
				if slices.Contains(allowedTypes, message.MessageType) && strings.HasPrefix(message.Message, `!!!`) {
					filteredMessageStream <- message
				}
			}
		}
	}()

	return filteredMessageStream
}

type FindUser func(username string) *User
type FindChannel func(user *User, channelName string) *Channel
type SendMessage func(ctx context.Context, user *User, channel *Channel, message string) error
type GPT func(ctx context.Context, query string) (string, error)

func ServeMessageStream(ctx context.Context, messagesStream <-chan *Message, findUser FindUser, findChannel FindChannel, sendMessage SendMessage, gpt GPT) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-messagesStream:
				user := findUser(message.Username)
				channel := findChannel(user, message.ChannelName)
				if user == nil || channel == nil {
					continue
				}
				answer, err := gpt(ctx, strings.TrimPrefix(message.Message, "!!!"))
				if err != nil {
					log.Err(err).Msg("gpt query failed")
				}
				if err := sendMessage(ctx, user, channel, answer); err != nil {
					log.Err(err).Msg(`error while sending a twitch message`)
				}
			}
		}
	}()
}
