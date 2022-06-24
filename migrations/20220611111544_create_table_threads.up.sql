CREATE TABLE threads
(
    id           char(9)      NOT NULL,
    title        varchar(255) NOT NULL,
    description  text         NOT NULL,
    total_viewer int          NOT NULL DEFAULT 0,
    creator_id   char(8)      NOT NULL,
    category_id  char(5)      NOT NULL,
    created_at   timestamp    NOT NULL DEFAULT current_timestamp,
    updated_at   timestamp    NOT NULL DEFAULT current_timestamp,
    primary key (id),
    constraint fk_threads_users
        foreign key (creator_id)
            references users (id) on delete cascade,
    constraint fk_threads_categories
        foreign key (category_id)
            references categories (id) on delete cascade
);
