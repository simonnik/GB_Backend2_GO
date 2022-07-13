-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "activities"
(
    "user_id" INT,
    "date"    TIMESTAMP,
    "name"    VARCHAR
) PARTITION BY RANGE ("date");
CREATE INDEX "activities_user_id_date" ON "activities" ("user_id", "date");
CREATE TABLE "activities_202011" PARTITION OF "activities" FOR VALUES FROM
    ('2020-11-01'::TIMESTAMP) TO ('2020-12-01'::TIMESTAMP);
CREATE TABLE "activities_202012" PARTITION OF "activities" FOR VALUES FROM
    ('2020-12-01'::TIMESTAMP) TO ('2021-01-01'::TIMESTAMP);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE activities;
-- +goose StatementEnd
