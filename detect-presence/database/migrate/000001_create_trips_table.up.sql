CREATE TABLE IF NOT EXISTS trips
(
    id          text primary key not null,
    left_at     integer          not null,
    returned_at integer
);
