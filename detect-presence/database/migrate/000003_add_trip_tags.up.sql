CREATE TABLE IF NOT EXISTS trip_taggings
(
    trip_id uuid not null,
    tag     text not null,
    constraint fk_trip
        foreign key (trip_id) references trips (id)
);

CREATE UNIQUE INDEX trip_taggings_unique_idx ON trip_taggings (trip_id, tag);
