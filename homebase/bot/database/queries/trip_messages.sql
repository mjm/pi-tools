-- name: SetMessageForTrip :exec
INSERT INTO trip_messages (trip_id, message_id)
VALUES ($1, $2);

-- name: GetTripForMessage :one
SELECT trip_id
FROM trip_messages
WHERE message_id = $1
LIMIT 1;
