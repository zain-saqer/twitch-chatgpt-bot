package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zain-saqer/twitch-chatgpt/internal/bot"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"github.com/zain-saqer/twitch-chatgpt/internal/db"
	"net/http"
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
		Whitelist:     map[string]*chat.Username{},
		WhitelistByID: map[uuid.UUID]*chat.Username{},
	}
	if err := app.StartMessagePipeline(ctx); err != nil {
		sentry.CaptureException(err)
		log.Fatal().Err(err).Stack().Msg(`error starting the message pipeline`)
	}
	e := echo.New()
	e.Debug = config.Debug
	server := NewServer(app, e, config)
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
