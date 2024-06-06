package db

import (
	"context"
	"database/sql"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/zain-saqer/twitch-chatgpt/internal/chat"
	"time"
)

//go:embed init.sql
var initSql string

type SqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SqliteRepository {
	return &SqliteRepository{
		db: db,
	}
}

func (repo *SqliteRepository) PrepareDatabase(ctx context.Context) error {
	_, err := repo.db.ExecContext(ctx, initSql)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) GetChannelsByUser(ctx context.Context, userId string) (channels []*chat.Channel, err error) {
	rows, err := repo.db.QueryContext(ctx, `select id, username, createdAt from channel where user_id = ?`, userId)
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
		var id string
		var name string
		var createdAtStr string
		err = rows.Scan(&id, &name, &createdAtStr)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}
		channels = append(channels, &chat.Channel{ID: id, Name: name, UserId: userId, CreatedAt: createdAt})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (repo *SqliteRepository) SaveChannel(ctx context.Context, channel *chat.Channel) error {
	stmt, err := repo.db.PrepareContext(ctx, `insert into channel (id, username, user_id, createdAt) values (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	_, err = stmt.Exec(channel.ID, channel.Name, channel.UserId, channel.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) DeleteChannel(ctx context.Context, id string) error {
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

func (repo *SqliteRepository) GetChannel(ctx context.Context, id string) (channel *chat.Channel, err error) {
	stmt, err := repo.db.PrepareContext(ctx, `select username, user_id, createdAt from channel where id = ?`)
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
	var name string
	var userId string
	var createdAtStr string
	err = row.Scan(&name, &userId, &createdAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}
	channel = &chat.Channel{ID: id, Name: name, UserId: userId, CreatedAt: createdAt}
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
		var id string
		var username string
		var accessToken string
		var refreshToken string
		var expiresAtStr string
		var createdAtStr string
		err = rows.Scan(&id, &username, &accessToken, &refreshToken, &expiresAtStr, &createdAtStr)
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
	_, err = stmt.Exec(user.ID, user.Username, user.AccessToken, user.RefreshToken, user.ExpiresAt.Format(time.RFC3339), user.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) DeleteUser(ctx context.Context, id string) error {
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

func (repo *SqliteRepository) GetUser(ctx context.Context, id string) (user *chat.User, err error) {
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

func (repo *SqliteRepository) UpdateUser(ctx context.Context, user *chat.User) error {
	stmt, err := repo.db.PrepareContext(ctx, `update user set username=?, access_token=?, refresh_token=?, expires_at=? where id = ?`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	_, err = stmt.Exec(user.Username, user.AccessToken, user.RefreshToken, user.ExpiresAt.Format(time.RFC3339), user.ID)
	if err != nil {
		return err
	}
	return nil
}
