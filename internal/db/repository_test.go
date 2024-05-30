package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"path"
	"testing"
	"time"
)

func TestSqliteRepository(t *testing.T) {
	dbPath := path.Join(t.TempDir(), "sqlite.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	err = repo.PrepareDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	channel := &chat.Channel{ID: uuid.New(), Name: `Name`, CreatedAt: time.Now()}
	t.Run("Test SaveChannel", func(t *testing.T) {
		err := repo.SaveChannel(context.Background(), channel)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Test GetChannel", func(t *testing.T) {
		channel2, err := repo.GetChannel(context.Background(), channel.ID)
		if err != nil {
			t.Fatal(err)
		}
		if channel.ID != channel2.ID {
			t.Fatal("Channel IDs don't match")
		}
	})
	t.Run("Test GetChannels", func(t *testing.T) {
		channels, err := repo.GetChannels(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if len(channels) != 1 {
			t.Fatal("Expected 1 channel, got ", len(channels))
		}
		if channels[0].ID != channel.ID {
			t.Fatal("Expected channel id ", channel.ID, "got ", channels[0].ID)
		}
	})
	t.Run("Test DeleteChannel", func(t *testing.T) {
		err := repo.DeleteChannel(context.Background(), channel.ID)
		if err != nil {
			t.Fatal(err)
		}
		channels, err := repo.GetChannels(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if len(channels) != 0 {
			t.Fatal("Expected 0 channel, got ", len(channels))
		}
	})

	username := &chat.Username{ID: uuid.New(), Name: `Name`, CreatedAt: time.Now()}
	t.Run("Test SaveUsername", func(t *testing.T) {
		err := repo.SaveUsername(context.Background(), username)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Test GetUsernames", func(t *testing.T) {
		usernames, err := repo.GetUsernames(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if len(usernames) != 1 {
			t.Fatal("Expected 1 username, got ", len(usernames))
		}
		if usernames[0].ID != username.ID {
			t.Fatal("Expected username id ", username.ID, "got ", usernames[0].ID)
		}
	})
}
