CREATE TABLE trip_taggings
(
    trip_id string not null,
    tag string not null
);

CREATE UNIQUE INDEX trip_taggings_unique_idx ON trip_taggings (trip_id, tag);
