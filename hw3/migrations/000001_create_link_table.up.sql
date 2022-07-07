BEGIN;

CREATE TABLE IF NOT EXISTS links
(
    id serial PRIMARY KEY,
    link text NOT NULL,
    token VARCHAR (25)
);

COMMIT;