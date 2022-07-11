/*
инициализация БД
*/

USE test;

CREATE TABLE IF NOT EXISTS entities
(
    id   INT PRIMARY KEY,
    data VARCHAR(32)
);

insert into entities (id, data)
values (1, "data_one"),
       (2, "data_two");