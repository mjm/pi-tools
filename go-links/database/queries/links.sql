-- name: ListRecentLinks :many
SELECT *
FROM links
ORDER BY created_at DESC
LIMIT 30;

-- name: GetLink :one
SELECT *
FROM links
WHERE id = $1;

-- name: GetLinkByShortURL :one
SELECT *
FROM links
WHERE short_url = $1;

-- name: CreateLink :one
INSERT INTO links
    (id, short_url, destination_url, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateLink :one
UPDATE links
SET short_url       = $2,
    destination_url = $3,
    description     = $4
WHERE id = $1
RETURNING *;
