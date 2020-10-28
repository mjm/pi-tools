CREATE TABLE IF NOT EXISTS trips
(
    id          uuid primary key not null,
    left_at     timestamp        not null,
    returned_at timestamp
);
