package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chatgpt/internal/bot"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"github.com/zain-saqer/twitch-chatgpt/internal/db"
	"golang.org/x/oauth2"
	oauth2Twitch "golang.org/x/oauth2/twitch"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	config := getConfigs()
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDsn,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Fatal().Err(err).Msg("Sentry initialization failed")
	}
	defer sentry.Flush(2 * time.Second)

	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	twitchIrcClient := twitch.NewAnonymousClient()

	database, err := sql.Open("sqlite3", config.SqliteDbPath)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal().Err(err).Msg("Could not connect to database")
	}
	repo := db.NewRepository(database)
	if err := repo.PrepareDatabase(ctx); err != nil {
		sentry.CaptureException(err)
		log.Fatal().Err(err).Stack().Msg(`error while preparing clickhouse database`)
	}

	app := &bot.App{
		Repository:    repo,
		TwitchClient:  twitchIrcClient,
		Whitelist:     map[string]*chat.User{},
		WhitelistByID: map[uuid.UUID]*chat.User{},
	}
	if err := app.StartMessagePipeline(ctx); err != nil {
		sentry.CaptureException(err)
		log.Fatal().Err(err).Stack().Msg(`error starting the message pipeline`)
	}
	e := echo.New()
	e.Debug = config.Debug
	cookieStore := sessions.NewCookieStore([]byte(config.Secret))
	oauthRedirect, err := url.JoinPath(config.Domain, `/add-user/redirect`)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal().Err(err).Stack().Msg(`error creating oauth redirect URL`)
	}
	oauth2Config := &oauth2.Config{
		ClientID:     config.Oauth2ClientID,
		ClientSecret: config.Oauth2Secret,
		Scopes:       []string{"user:read:email"},
		Endpoint:     oauth2Twitch.Endpoint,
		RedirectURL:  oauthRedirect,
	}
	server := NewServer(app, e, config, cookieStore, oauth2Config)
	server.middlewares()
	server.setupRoutes()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := e.Start(config.ServerAddress)
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			sentry.CaptureException(err)
			log.Fatal().Err(err).Msg(`shutting down server error`)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg(`shutting down...`)
				if err := e.Shutdown(ctx); err != nil {
					sentry.CaptureException(err)
					log.Error().Err(err).Msg(`error while shutting down the web server`)
				}
				return
			}
		}
	}()

	wg.Wait()
}
