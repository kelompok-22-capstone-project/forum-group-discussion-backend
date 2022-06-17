CREATE TABLE likes
(
    id         char(9),
    user_id    char(8)   NOT NULL,
    thread_id  char(9)   NOT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp,
    updated_at timestamp NOT NULL DEFAULT current_timestamp,
    primary key (id),
    constraint fk_likes_users foreign key (user_id) references users (id) on delete cascade,
    constraint fk_likes_threads foreign key (thread_id) references threads (id) on delete cascade,
    unique (user_id, thread_id)
);
