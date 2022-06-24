CREATE TABLE categories
(
    id          char(5),
    name        varchar(50) NOT NULL,
    description text        NOT NULL,
    created_at  timestamp   NOT NULL DEFAULT current_timestamp,
    updated_at  timestamp   NOT NULL DEFAULT current_timestamp,
    primary key (id)
);
