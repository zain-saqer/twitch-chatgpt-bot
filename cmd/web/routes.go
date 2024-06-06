package main

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"github.com/zain-saqer/twitch-chatgpt/web"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"
)

func (s *Server) setupRoutes() {
	route := s.Echo.Group(`/`, authMiddleware(s.Config))
	route.GET(``, s.getIndex)
	route.GET(`:userId/channels`, s.getAdminChannels)
	route.GET(`:userId/add-channel`, s.getAdminAddChannel)
	route.POST(`:userId/add-channel`, s.postAdminAddChannel)
	route.DELETE(`channels/:id`, s.deleteAdminDeleteChannel)
	route.DELETE(`users/:id`, s.deleteAdminDeleteUser)

	route.GET(`add-user`, s.getAddUser)
	route.GET(`add-user/redirect`, s.getOAuth2Callback)
}

func (s *Server) getIndex(c echo.Context) error {
	var t *template.Template
	sync.OnceFunc(func() {
		var err error
		t, err = template.ParseFS(web.F, `templates/layout.gohtml`, `templates/nav.gohtml`, `templates/index.gohtml`)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal().Err(err).Stack().Msg(`error parsing templates`)
		}
	})()
	usernames, err := s.App.Repository.GetUsers(c.Request().Context())
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(c.Response(), `base`, IndexView{Users: usernames})
}

func (s *Server) getAdminChannels(c echo.Context) error {
	var t *template.Template
	sync.OnceFunc(func() {
		var err error
		t, err = template.ParseFS(web.F, `templates/layout.gohtml`, `templates/nav.gohtml`, `templates/channels.gohtml`)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal().Err(err).Stack().Msg(`error parsing templates`)
		}
	})()
	channels, err := s.App.Repository.GetChannelsByUser(c.Request().Context(), c.Param(`userId`))
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(c.Response(), `base`, UserView{UserID: c.Param(`userId`), Channels: channels})
}

func (s *Server) getAdminAddChannel(c echo.Context) error {
	var t *template.Template
	sync.OnceFunc(func() {
		var err error
		t, err = template.ParseFS(web.F, `templates/layout.gohtml`, `templates/nav.gohtml`, `templates/add_channel.gohtml`)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal().Err(err).Stack().Msg(`error parsing templates`)
		}
	})()
	return t.ExecuteTemplate(c.Response(), `base`, AddChannel{})
}

func (s *Server) postAdminAddChannel(c echo.Context) error {
	var t *template.Template
	sync.OnceFunc(func() {
		var err error
		t, err = template.ParseFS(web.F, `templates/layout.gohtml`, `templates/nav.gohtml`, `templates/add_channel.gohtml`)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal().Err(err).Stack().Msg(`error parsing templates`)
		}
	})()
	addChannel := &AddChannel{}
	err := c.Bind(addChannel)
	if err != nil {
		return err
	}
	addChannel.Trim()
	if !addChannel.Validate() {
		return t.ExecuteTemplate(c.Response(), `base`, addChannel)
	}
	user, err := s.App.Repository.GetUser(c.Request().Context(), addChannel.UserId)
	if err != nil {
		return err
	}

	twitchChannel, err := s.App.TwitterAPI.GetUser(c.Request().Context(), user.AccessToken, addChannel.Username)
	if err != nil {
		return err
	}
	channel := &chat.Channel{ID: twitchChannel.ID, UserId: addChannel.UserId, Name: addChannel.Username, CreatedAt: time.Now()}
	if err = s.App.Repository.SaveChannel(c.Request().Context(), channel); err != nil {
		return err
	}
	s.App.AddChannel(user, channel)
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf(`/%s/channels`, c.Param(`userId`)))
}

func (s *Server) deleteAdminDeleteChannel(c echo.Context) error {
	deleteChannel := &DeleteChannel{}
	err := c.Bind(deleteChannel)
	if err != nil {
		return err
	}
	id := strings.TrimSpace(deleteChannel.ID)
	if id == `` {
		return errors.New(`invalid request`)
	}
	channel, err := s.App.Repository.GetChannel(c.Request().Context(), id)
	if err != nil {
		return err
	}
	err = s.App.Repository.DeleteChannel(c.Request().Context(), id)
	if err != nil {
		return err
	}
	s.App.Depart(channel.Name)
	c.Response().Header().Add(`HX-Refresh`, `true`)
	return c.String(http.StatusOK, ``)
}

func (s *Server) deleteAdminDeleteUser(c echo.Context) error {
	userChannel := &DeleteUser{}
	err := c.Bind(userChannel)
	if err != nil {
		return err
	}
	id := strings.TrimSpace(userChannel.ID)
	if id == `` {
		return errors.New(`invalid request`)
	}
	user, err := s.App.Repository.GetUser(c.Request().Context(), id)
	if err != nil {
		return err
	}
	err = s.App.Repository.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return err
	}
	s.App.RemoveUser(user)
	c.Response().Header().Add(`HX-Refresh`, `true`)
	return c.String(http.StatusOK, ``)
}
