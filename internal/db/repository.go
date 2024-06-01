package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"time"
)

type SqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SqliteRepository {
	return &SqliteRepository{
		db: db,
	}
}

func (repo *SqliteRepository) PrepareDatabase(ctx context.Context) error {
	_, err := repo.db.ExecContext(ctx, `create table if not exists channel (id  TEXT, name TEXT NOT NULL, createdAt  TEXT NOT NULL); create table if not exists user (id TEXT, username TEXT NOT NULL, access_token TEXT NOT NULL, refresh_token TEXT NOT NULL, expires_at TEXT NOT NULL, created_at  TEXT NOT NULL);`)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) GetChannels(ctx context.Context) (channels []*chat.Channel, err error) {
	rows, err := repo.db.QueryContext(ctx, `select * from channel`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_err := rows.Close()
		if _err != nil {
			err = _err
		}
	}(rows)
	channels = make([]*chat.Channel, 0)
	for rows.Next() {
		var idStr string
		var name string
		var createdAtStr string
		err = rows.Scan(&idStr, &name, &createdAtStr)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		channels = append(channels, &chat.Channel{ID: id, Name: name, CreatedAt: createdAt})
		return channels, nil
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (repo *SqliteRepository) SaveChannel(ctx context.Context, channel *chat.Channel) error {
	stmt, err := repo.db.PrepareContext(ctx, `insert into channel (id, name, createdAt) values (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	_, err = stmt.Exec(channel.ID.String(), channel.Name, channel.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) DeleteChannel(ctx context.Context, id uuid.UUID) error {
	stmt, err := repo.db.PrepareContext(ctx, `delete from channel where id = ?`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) GetChannel(ctx context.Context, id uuid.UUID) (channel *chat.Channel, err error) {
	stmt, err := repo.db.PrepareContext(ctx, `select * from channel where id = ?`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	row := stmt.QueryRow(id)
	if err != nil {
		return nil, err
	}
	var idStr string
	var name string
	var createdAtStr string
	err = row.Scan(&idStr, &name, &createdAtStr)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}
	id, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	channel = &chat.Channel{ID: id, Name: name, CreatedAt: createdAt}
	return channel, nil
}

func (repo *SqliteRepository) GetUsers(ctx context.Context) (users []*chat.User, err error) {
	rows, err := repo.db.QueryContext(ctx, `select id, username, access_token, refresh_token, expires_at, created_at from user`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_err := rows.Close()
		if _err != nil {
			err = _err
		}
	}(rows)
	users = make([]*chat.User, 0)
	for rows.Next() {
		var idStr string
		var username string
		var accessToken string
		var refreshToken string
		var expiresAtStr string
		var createdAtStr string
		err = rows.Scan(&idStr, &username, &accessToken, &refreshToken, &expiresAtStr, &createdAtStr)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}
		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			return nil, err
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		users = append(users, &chat.User{ID: id, Username: username, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresAt: expiresAt, CreatedAt: createdAt})
	}
	return
}

func (repo *SqliteRepository) SaveUser(ctx context.Context, user *chat.User) error {
	stmt, err := repo.db.PrepareContext(ctx, `insert into user (id, username, access_token, refresh_token, expires_at, created_at) values (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	_, err = stmt.Exec(user.ID.String(), user.Username, user.AccessToken, user.RefreshToken, user.ExpiresAt.Format(time.RFC3339), user.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	stmt, err := repo.db.PrepareContext(ctx, `delete from user where id = ?`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) GetUser(ctx context.Context, id uuid.UUID) (user *chat.User, err error) {
	stmt, err := repo.db.PrepareContext(ctx, `select username, access_token, refresh_token, expires_at, created_at from user where id = ?`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	row := stmt.QueryRow(id)
	if err != nil {
		return nil, err
	}
	var username string
	var accessToken string
	var refreshToken string
	var expiresAtStr string
	var createdAtStr string
	err = row.Scan(&username, &accessToken, &refreshToken, &expiresAtStr, &createdAtStr)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}
	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		return nil, err
	}
	user = &chat.User{ID: id, Username: username, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresAt: expiresAt, CreatedAt: createdAt}
	return user, nil
}
