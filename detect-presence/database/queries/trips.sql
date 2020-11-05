-- name: ListTrips :many
SELECT id, left_at, returned_at, array_remove(array_agg(DISTINCT tag ORDER BY tag), NULL)::text[] as tags
FROM trips
         LEFT JOIN trip_taggings tt ON trips.id = tt.trip_id
WHERE ignored_at IS NULL
GROUP BY id, left_at
ORDER BY left_at DESC
LIMIT 30;

-- name: GetTrip :one
SELECT id, left_at, returned_at, array_remove(array_agg(DISTINCT tag ORDER BY tag), NULL)::text[] as tags
FROM trips
         LEFT JOIN trip_taggings tt ON trips.id = tt.trip_id
WHERE id = $1
GROUP BY id
LIMIT 30;

-- name: GetLastCompletedTrip :one
SELECT *
FROM trips
WHERE ignored_at IS NULL
AND returned_at IS NOT NULL
ORDER BY left_at DESC
LIMIT 1;

-- name: GetCurrentTrip :one
SELECT *
FROM trips
WHERE ignored_at IS NULL
  AND returned_at IS NULL
ORDER BY left_at DESC
LIMIT 1;

-- name: BeginTrip :one
INSERT INTO trips (id, left_at)
VALUES ($1, $2)
RETURNING *;

-- name: EndTrip :exec
UPDATE trips
SET returned_at = $2
WHERE id = $1
AND returned_at IS NULL;

-- name: IgnoreTrip :execrows
UPDATE trips
SET ignored_at = CURRENT_TIMESTAMP
WHERE id = $1;
