-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "users"
(
    "user_id" INT,
    "name"    VARCHAR,
    "age"     INT,
    "spouse"  INT
);
CREATE UNIQUE INDEX "users_user_id" ON "users" ("user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd
