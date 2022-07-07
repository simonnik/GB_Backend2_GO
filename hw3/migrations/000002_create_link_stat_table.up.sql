BEGIN;

CREATE TABLE IF NOT EXISTS links_stat
(
    id         serial PRIMARY KEY,
    link_id    int     NOT NULL constraint links_stat_links_id_fk references links,
    ip         varchar NOT NULL,
    created_at timestamp default current_timestamp
);

COMMIT ;