-- name: ListTags :many
SELECT tag AS name, COUNT(trip_id) AS trip_count
FROM trip_taggings
GROUP BY tag
ORDER BY COUNT(trip_id) DESC
LIMIT $1;

-- name: TagTrip :exec
INSERT INTO trip_taggings (trip_id, tag)
VALUES ($1, $2);

-- name: UntagTrip :exec
DELETE
FROM trip_taggings
WHERE trip_id = $1
  AND tag = $2;
