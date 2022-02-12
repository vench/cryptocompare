CREATE DATABASE cryptocompare;
CREATE USER 'user_cryptocompare'@'localhost' IDENTIFIED BY 'password_cryptocompare';
GRANT ALL PRIVILEGES ON cryptocompare.* TO 'user_cryptocompare'@'%' IDENTIFIED BY 'mysql';
GRANT ALL PRIVILEGES ON cryptocompare.* TO 'user_cryptocompare'@'localhost' IDENTIFIED BY 'mysql';
USE cryptocompare;

-- auto-generated definition
create table currency
(
    `key`      varchar(8)                          not null  primary key,
    value      json                                null,
    created_at timestamp default CURRENT_TIMESTAMP null,
    updated_at timestamp default CURRENT_TIMESTAMP null
);
