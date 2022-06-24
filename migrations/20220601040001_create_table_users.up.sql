CREATE TABLE users
(
    id         char(8),
    username   varchar(20) NOT NULL UNIQUE,
    email      varchar(50) NOT NULL UNIQUE,
    name       varchar(50) NOT NULL,
    password   varchar(60) NOT NULL,
    role       roles       NOT NULL,
    is_active  bool        NOT NULL DEFAULT true,
    created_at timestamp   NOT NULL DEFAULT current_timestamp,
    updated_at timestamp   NOT NULL DEFAULT current_timestamp,
    primary key (id)
);
