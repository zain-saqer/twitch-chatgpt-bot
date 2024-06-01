package main

import (
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"strings"
)

type IndexView struct {
	Channels []*chat.Channel
	Users    []*chat.User
}

type AddChannel struct {
	Errors []string
	Name   string `form:"name"`
}

func (c *AddChannel) Trim() {
	c.Name = strings.TrimSpace(c.Name)
}

func (c *AddChannel) Validate() bool {
	errors := make([]string, 0)
	if c.Name == "" {
		errors = append(errors, "Username is required")
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
