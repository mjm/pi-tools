-- name: TagTrip :exec
INSERT INTO trip_taggings (trip_id, tag)
VALUES ($1, $2);

-- name: UntagTrip :exec
DELETE FROM trip_taggings
WHERE trip_id = $1
AND tag = $2;
