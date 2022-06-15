CREATE TABLE user_follows
(
    id           char(9),
    user_id      char(8)   NOT NULL,
    following_id char(8)   NOT NULL,
    created_at   timestamp NOT NULL DEFAULT current_timestamp,
    updated_at   timestamp NOT NULL DEFAULT current_timestamp,
    primary key (id),
    constraint fk_user_follows_user_users foreign key (user_id) references users (id) on delete cascade,
    constraint fk_user_follows_following_users foreign key (following_id) references users (id) on delete cascade,
    unique (user_id, following_id)
);
