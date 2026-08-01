-- name: GetRoomByID :one
SELECT  * 
FROM    rooms 
WHERE   id = $1
RETURNING *;

-- name: CreateRoom :exec
INSERT INTO rooms (name, description)
VALUES ($1, $2);