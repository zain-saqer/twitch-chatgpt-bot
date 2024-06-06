package main

import (
	"github.com/zain-saqer/twitch-chatgpt/internal/env"
	"os"
)

type Config struct {
	Debug                bool
	AuthUser             string
	AuthPass             string
	ServerAddress        string
	SqliteDbPath         string
	SentryDsn            string
	Secret               string
	Domain               string
	Oauth2ClientID       string
	Oauth2Secret         string
	OpenAIAPIKey         string
	ChatGPTSystemMessage string
	ChatGPTModel         string
}

func getConfigs() *Config {
	_, debug := os.LookupEnv(`DEBUG`)
	return &Config{
		Debug:                debug,
		ServerAddress:        env.MustGetEnv(`SERVER_ADDRESS`),
		SqliteDbPath:         env.MustGetEnv(`SQLITE_DB_PATH`),
		AuthUser:             env.MustGetEnv(`AUTH_USER`),
		AuthPass:             env.MustGetEnv(`AUTH_PASS`),
		SentryDsn:            env.MustGetEnv(`SENTRY_DSN`),
		Secret:               env.MustGetEnv(`SECRET`),
		Domain:               env.MustGetEnv(`DOMAIN`),
		Oauth2ClientID:       env.MustGetEnv(`OAUTH2_CLIENT_ID`),
		Oauth2Secret:         env.MustGetEnv(`OAUTH2_CLIENT_SECRET`),
		OpenAIAPIKey:         env.MustGetEnv(`OPENAI_API_KEY`),
		ChatGPTSystemMessage: env.MustGetEnv(`CHAT_GPT_SYSTEM_MESSAGE`),
		ChatGPTModel:         env.MustGetEnv(`CHAT_GPT_MODEL`),
	}
}
