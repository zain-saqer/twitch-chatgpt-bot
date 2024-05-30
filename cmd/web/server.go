package main

import (
	"github.com/labstack/echo/v4"
	"github.com/zain-saqer/twitch-chatgpt/internal/bot"
)

type Server struct {
	App    *bot.App
	Echo   *echo.Echo
	Config *Config
}

func NewServer(app *bot.App, e *echo.Echo, config *Config) *Server {
	return &Server{App: app, Echo: e, Config: config}
}
