USE cryptocompare;

create table currency
(
    `key`      varchar(8)                          not null  primary key,
    value      json                                null,
    created_at timestamp default CURRENT_TIMESTAMP null,
    updated_at timestamp default CURRENT_TIMESTAMP null
);
