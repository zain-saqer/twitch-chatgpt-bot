package chat

import (
	"context"
	"fmt"
	"slices"
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
				if slices.Contains(allowedTypes, message.MessageType) {
					filteredMessageStream <- message
				}
			}
		}
	}()

	return filteredMessageStream
}

func ServeMessageStream(ctx context.Context, messagesStream <-chan *Message, getWhiteList func() map[string]*User) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-messagesStream:
				usernames := getWhiteList()
				if _, ok := usernames[message.Username]; !ok {
					continue
				}
				fmt.Println(message.Message)
			}
		}
	}()
}
