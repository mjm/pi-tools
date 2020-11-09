-- name: SetMessageForTrip :exec
INSERT INTO trip_messages (trip_id, message_id)
VALUES ($1, $2)
ON CONFLICT (trip_id) DO UPDATE
SET message_id = EXCLUDED.message_id;

-- name: GetMessageForTrip :one
SELECT message_id
FROM trip_messages
WHERE trip_id = $1
LIMIT 1;

-- name: GetTripForMessage :one
SELECT trip_id
FROM trip_messages
WHERE message_id = $1
LIMIT 1;
