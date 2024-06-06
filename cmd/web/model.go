package main

import (
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"strings"
)

type IndexView struct {
	Users []*chat.User
}

type UserView struct {
	Channels []*chat.Channel
	UserID   string
}

type AddChannel struct {
	Errors   []string
	Username string `form:"name"`
	UserId   string `param:"userId"`
}

func (c *AddChannel) Trim() {
	c.Username = strings.TrimSpace(c.Username)
	c.UserId = strings.TrimSpace(c.UserId)
}

func (c *AddChannel) Validate() bool {
	errors := make([]string, 0)
	if c.Username == "" {
		errors = append(errors, "Username is required")
	}
	if c.UserId == "" {
		errors = append(errors, "UserId is required")
	}
	c.Errors = errors
	return len(errors) == 0
}

type DeleteChannel struct {
	ID string `param:"id"`
}

type DeleteUser struct {
	ID string `param:"id"`
}
