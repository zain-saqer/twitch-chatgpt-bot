package main

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
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
	route.GET(`add-channel`, s.getAdminAddChannel)
	route.POST(`add-channel`, s.postAdminAddChannel)
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
	channels, err := s.App.Repository.GetChannels(c.Request().Context())
	if err != nil {
		return err
	}
	usernames, err := s.App.Repository.GetUsers(c.Request().Context())
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(c.Response(), `base`, IndexView{Channels: channels, Users: usernames})
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
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if err = s.App.Repository.SaveChannel(c.Request().Context(), &chat.Channel{ID: id, Name: addChannel.Name, CreatedAt: time.Now()}); err != nil {
		return err
	}
	s.App.JoinChannel(addChannel.Name)
	return c.Redirect(http.StatusSeeOther, `/`)
}

func (s *Server) deleteAdminDeleteChannel(c echo.Context) error {
	deleteChannel := &DeleteChannel{}
	err := c.Bind(deleteChannel)
	if err != nil {
		return err
	}
	idString := strings.TrimSpace(deleteChannel.ID)
	if idString == `` {
		return errors.New(`invalid request`)
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return err
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
	idString := strings.TrimSpace(userChannel.ID)
	if idString == `` {
		return errors.New(`invalid request`)
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return err
	}
	username, err := s.App.Repository.GetUser(c.Request().Context(), id)
	if err != nil {
		return err
	}
	err = s.App.Repository.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return err
	}
	s.App.RemoveUsername(username)
	c.Response().Header().Add(`HX-Refresh`, `true`)
	return c.String(http.StatusOK, ``)
}
