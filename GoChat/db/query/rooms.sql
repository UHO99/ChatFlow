-- name: GetRoomByID :one
SELECT  *
FROM    rooms
WHERE   id = $1;

-- name: GetListRoom :many
SELECT  *
FROM    rooms
ORDER BY id;

-- name: CreateRoom :one
INSERT INTO rooms (name, description)
VALUES ($1, $2)
RETURNING *;