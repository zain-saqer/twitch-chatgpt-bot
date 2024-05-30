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
	_, err := repo.db.ExecContext(ctx, `create table if not exists channel (id  TEXT, name TEXT NOT NULL, createdAt  TEXT NOT NULL); create table if not exists username_whitelist (id  TEXT, name TEXT NOT NULL, createdAt  TEXT NOT NULL);`)
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

func (repo *SqliteRepository) GetUsernames(ctx context.Context) (usernames []*chat.Username, err error) {
	rows, err := repo.db.QueryContext(ctx, `select * from username_whitelist`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_err := rows.Close()
		if _err != nil {
			err = _err
		}
	}(rows)
	usernames = make([]*chat.Username, 0)
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
		usernames = append(usernames, &chat.Username{ID: id, Name: name, CreatedAt: createdAt})
	}
	return
}

func (repo *SqliteRepository) SaveUsername(ctx context.Context, username *chat.Username) error {
	stmt, err := repo.db.PrepareContext(ctx, `insert into username_whitelist (id, name, createdAt) values (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_err := stmt.Close()
		if _err != nil {
			err = _err
		}
	}(stmt)
	_, err = stmt.Exec(username.ID.String(), username.Name, username.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteRepository) DeleteUsername(ctx context.Context, id uuid.UUID) error {
	stmt, err := repo.db.PrepareContext(ctx, `delete from username_whitelist where id = ?`)
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

func (repo *SqliteRepository) GetUsername(ctx context.Context, id uuid.UUID) (username *chat.Username, err error) {
	stmt, err := repo.db.PrepareContext(ctx, `select * from username_whitelist where id = ?`)
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
	username = &chat.Username{ID: id, Name: name, CreatedAt: createdAt}
	return username, nil
}
