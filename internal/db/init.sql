create table if not exists channel
(
    id        TEXT NOT NULL ,
    name      TEXT NOT NULL,
    createdAt TEXT NOT NULL
);
create table if not exists user
(
    id            TEXT NOT NULL,
    username      TEXT NOT NULL,
    access_token  TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at    TEXT NOT NULL,
    created_at    TEXT NOT NULL
);

create unique index if not exists CHANNEL_ID_INDEX on channel (id);
create unique index if not exists CHANNEL_NAME_INDEX on channel (name);
create unique index if not exists USER_ID_INDEX on user (id);
create unique index if not exists USER_USERNAME_INDEX on user (username);