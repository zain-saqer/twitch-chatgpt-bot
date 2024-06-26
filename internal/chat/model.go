package chat

import (
	"context"
	"time"
)

type Message struct {
	Username    string
	ChannelName string
	Message     string
	MessageType uint8
	Time        time.Time
}

type Channel struct {
	ID        string
	Name      string
	UserId    string
	CreatedAt time.Time
}

type User struct {
	ID           string
	Username     string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

const (
	Unset           uint8 = iota
	Whisper         uint8 = iota
	PrivMsg         uint8 = iota
	ClearChat       uint8 = iota
	RoomState       uint8 = iota
	UserNotice      uint8 = iota
	UserState       uint8 = iota
	Notice          uint8 = iota
	Join            uint8 = iota
	Part            uint8 = iota
	Reconnect       uint8 = iota
	Names           uint8 = iota
	Ping            uint8 = iota
	Pong            uint8 = iota
	ClearMsg        uint8 = iota
	GlobalUserState uint8 = iota
)

type Repository interface {
	GetChannelsByUser(ctx context.Context, userId string) ([]*Channel, error)
	SaveChannel(ctx context.Context, channel *Channel) error
	GetChannel(ctx context.Context, id string) (*Channel, error)
	DeleteChannel(ctx context.Context, id string) error
	GetUsers(ctx context.Context) ([]*User, error)
	SaveUser(ctx context.Context, username *User) error
	DeleteUser(ctx context.Context, id string) error
	GetUser(ctx context.Context, id string) (user *User, err error)
	UpdateUser(ctx context.Context, user *User) error
}
