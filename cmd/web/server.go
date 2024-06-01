package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/zain-saqer/twitch-chatgpt/internal/bot"
	"golang.org/x/oauth2"
)

type Server struct {
	App          *bot.App
	Echo         *echo.Echo
	Config       *Config
	CookieStore  *sessions.CookieStore
	Oauth2Config *oauth2.Config
}

func NewServer(app *bot.App, e *echo.Echo, config *Config, cookieStore *sessions.CookieStore, oauth2Config *oauth2.Config) *Server {
	return &Server{App: app, Echo: e, Config: config, CookieStore: cookieStore, Oauth2Config: oauth2Config}
}
