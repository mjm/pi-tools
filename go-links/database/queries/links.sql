-- name: CreateLink :one
INSERT INTO links
    (id, short_url, destination_url, description)
VALUES ($1, $2, $3, $4)
RETURNING *;
