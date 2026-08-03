-- name: CreateRoomAddMember :exec
INSERT INTO room_members (room_id, user_id)
VALUES ($1, $2);

-- name: JoinRoomByUsername :execrows
INSERT INTO room_members (room_id, user_id)
SELECT  $1, id
FROM    users
WHERE   username = $2
ON CONFLICT (room_id, user_id) DO NOTHING;

-- name: GetRoomMemberByUsername :one
SELECT  room_members.room_id,
        room_members.user_id,
        room_members.joined_at
FROM    room_members
JOIN    users ON users.id = room_members.user_id
WHERE   room_members.room_id = $1
AND     users.username = $2;