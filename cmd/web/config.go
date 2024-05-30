package main

import (
	"github.com/zain-saqer/twitch-chatgpt/internal/env"
	"os"
)

type Config struct {
	Debug         bool
	AuthUser      string
	AuthPass      string
	ServerAddress string
	SqliteDbPath  string
	SentryDsn     string
}

func getConfigs() *Config {
	_, debug := os.LookupEnv(`DEBUG`)
	return &Config{
		Debug:         debug,
		ServerAddress: env.MustGetEnv(`SERVER_ADDRESS`),
		SqliteDbPath:  env.MustGetEnv(`SQLITE_DB_PATH`),
		AuthUser:      env.MustGetEnv(`AUTH_USER`),
		AuthPass:      env.MustGetEnv(`AUTH_PASS`),
		SentryDsn:     env.MustGetEnv(`SENTRY_DSN`),
	}
}
