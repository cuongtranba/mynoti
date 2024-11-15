-- name: CreateComicTracking :exec
INSERT INTO comic_tracking (url, name, description, html)
VALUES ($1, $2, $3, $4)
RETURNING id, url, name, description, html, last_checked;

-- name: GetAllComicTrackings :many
SELECT id, url, name, description, html, last_checked
FROM comic_tracking;

-- name: GetComicTrackingByID :one
SELECT id, url, name, description, html, last_checked
FROM comic_tracking
WHERE id = $1;

-- name: UpdateComicTracking :exec
UPDATE comic_tracking
SET url = $2, 
    name = $3, 
    description = $4, 
    html = $5, 
    last_checked = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, url, name, description, html, last_checked;

-- name: DeleteComicTracking :exec
DELETE FROM comic_tracking
WHERE id = $1
RETURNING id;

-- name: UpdateLastChecked :exec
UPDATE comic_tracking
SET last_checked = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, last_checked;

-- name: GetComicTrackingsList :many
SELECT id, url, name, description, html, last_checked
FROM comic_tracking
ORDER BY last_checked DESC
LIMIT $1 OFFSET $2;