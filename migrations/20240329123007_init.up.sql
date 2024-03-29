create table if not exists users
(
    id         uuid primary key unique not null,
    email      varchar(255) unique     not null,
    password   varchar(255)            not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

create table if not exists roles
(
    id         uuid primary key unique not null,
    name       varchar(255) unique     not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

create table if not exists users_roles
(
    user_id uuid not null references users (id),
    role_id uuid not null references roles (id),
    primary key (user_id, role_id)
)