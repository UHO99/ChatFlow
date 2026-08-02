-- name: CreateMessage :exec
INSERT INTO messages (room_id, user_id, content)
VALUES ($1, $2, $3);