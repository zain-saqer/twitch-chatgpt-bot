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
	user := &chat.User{ID: uuid.New().String(), Username: `Username`, CreatedAt: time.Now()}
	channel := &chat.Channel{ID: uuid.New().String(), Name: `Username`, CreatedAt: time.Now(), UserId: user.ID}
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
		channels, err := repo.GetChannelsByUser(context.Background(), user.ID)
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
	t.Run("Test GetChannel and DeleteChannel", func(t *testing.T) {
		channel2, err := repo.GetChannel(context.Background(), channel.ID)
		if err != nil {
			t.Fatal(err)
		}
		if channel2 == nil {
			t.Fatal("Expected channel got nil")
		}
		if channel2.ID != channel.ID {
			t.Fatal("Expected channel id ", channel.ID, "got ", channel2.ID)
		}
		err = repo.DeleteChannel(context.Background(), channel.ID)
		if err != nil {
			t.Fatal(err)
		}
		channel2, err = repo.GetChannel(context.Background(), channel.ID)
		if err != nil {
			t.Fatal(err)
		}
		if channel2 != nil {
			t.Fatal("Expected channel got nil")
		}
		channels, err := repo.GetChannelsByUser(context.Background(), user.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(channels) != 0 {
			t.Fatal("Expected 0 channel, got ", len(channels))
		}
	})
	t.Run("Test SaveUser", func(t *testing.T) {
		err := repo.SaveUser(context.Background(), user)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Test GetUsers", func(t *testing.T) {
		usernames, err := repo.GetUsers(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if len(usernames) != 1 {
			t.Fatal("Expected 1 user, got ", len(usernames))
		}
		if usernames[0].ID != user.ID {
			t.Fatal("Expected user id ", user.ID, "got ", usernames[0].ID)
		}
	})
}
