-- name: CreateRoomAddMember :exec
INSERT INTO room_members (room_id, user_id)
VALUES ($1, $2);