create table if not exists events
(
    id       uuid primary key unique not null,
    title    text                    not null,
    notice   text                    not null,
    source   varchar(255)            not null,
    image    varchar(255)            not null,
    link     varchar(255)            not null,
    category varchar(255)            not null
)