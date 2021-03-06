CREATE TABLE IF NOT EXISTS links
(
    id              text primary key not null,
    short_url       text             not null,
    destination_url text             not null,
    description     text             not null,
    created_at      timestamp        not null default current_timestamp,
    updated_at      timestamp        not null default current_timestamp
);

CREATE UNIQUE INDEX links_short_url_unique_idx ON links (short_url);
