create table if not exists channel
(
    id        TEXT NOT NULL,
    username  TEXT NOT NULL,
    createdAt TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY (id),
    foreign key (user_id) references user(id)
);
create table if not exists user
(
    id            TEXT NOT NULL,
    username      TEXT NOT NULL,
    access_token  TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at    TEXT NOT NULL,
    created_at    TEXT NOT NULL,
    PRIMARY KEY (id)
);

create unique index if not exists CHANNEL_NAME_INDEX on channel (username);
create unique index if not exists USER_USERNAME_INDEX on user (username);