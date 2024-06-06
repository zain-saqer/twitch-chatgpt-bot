package chat

import (
	"context"
	"fmt"
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
type IsUserChannel func(username, channelName string) bool

func ServeMessageStream(ctx context.Context, messagesStream <-chan *Message, findUser FindUser, isUserChannel IsUserChannel) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-messagesStream:
				user := findUser(message.Username)
				if user == nil || !isUserChannel(user.Username, message.ChannelName) {
					continue
				}
				fmt.Println(message.Message)
			}
		}
	}()
}
