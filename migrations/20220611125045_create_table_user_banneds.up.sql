CREATE TABLE user_banneds
(
    id           char(9),
    moderator_id char(6)   NOT NULL,
    user_id      char(8)   NOT NULL,
    comment_id   char(9)   NOT NULL,
    reason       text      NOT NULL,
    status       status    NOT NULL DEFAULT 'review',
    created_at   timestamp NOT NULL DEFAULT current_timestamp,
    updated_at   timestamp NOT NULL DEFAULT current_timestamp,
    primary key (id),
    constraint fk_user_banneds_moderators foreign key (moderator_id) references moderators (id) on delete cascade,
    constraint fk_user_banneds_users foreign key (user_id) references users (id) on delete cascade,
    constraint fk_user_banneds_comments foreign key (comment_id) references comments (id) on delete cascade,
    unique (moderator_id, user_id)
);
