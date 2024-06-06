package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"net/http"
	"time"
)

const (
	stateCallbackKey = "oauth-state-callback"
	oauthSessionName = "oauth-session"
)

func (s *Server) getAddUser(c echo.Context) (err error) {
	sess, err := session.Get(oauthSessionName, c)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		return err
	}

	state := hex.EncodeToString(tokenBytes[:])

	sess.AddFlash(state, stateCallbackKey)

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return
	}

	return c.Redirect(http.StatusTemporaryRedirect, s.Oauth2Config.AuthCodeURL(state))
}

func (s *Server) getOAuth2Callback(c echo.Context) (err error) {
	sess, err := session.Get(oauthSessionName, c)
	if err != nil {
		log.Err(err).Msgf("corrupted session %s -- generated new", err)
		err = nil
	}

	// ensure we flush the csrf challenge even if the request is ultimately unsuccessful
	defer func() {
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			log.Err(err).Msgf("error saving session: %s", err)
		}
	}()

	switch stateChallenge, state := sess.Flashes(stateCallbackKey), c.FormValue("state"); {
	case state == "", len(stateChallenge) < 1:
		err = errors.New("missing state challenge")
	case state != stateChallenge[len(stateChallenge)-1]:
		err = fmt.Errorf("invalid oauth state, expected '%s', got '%s'\n", state, stateChallenge[0])
	}
	if err != nil {
		return
	}

	token, err := s.Oauth2Config.Exchange(c.Request().Context(), c.FormValue("code"))
	if err != nil {
		return
	}
	twitchUser, err := s.App.TwitterAPI.GetCurrentUser(c.Request().Context(), token.AccessToken)
	if err != nil {
		return err
	}
	user := &chat.User{
		ID:           twitchUser.ID,
		Username:     twitchUser.Login,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
		CreatedAt:    time.Now(),
	}
	err = s.App.Repository.SaveUser(c.Request().Context(), user)
	if err != nil {
		return err
	}
	s.App.AddUser(user)
	err = c.Redirect(http.StatusTemporaryRedirect, "/")
	if err != nil {
		log.Err(err).Msg("")
	}

	return
}
