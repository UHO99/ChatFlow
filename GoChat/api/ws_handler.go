package api

import (
	db "GoChat/db/sqlc"
	"GoChat/middleware"
	"GoChat/servers/hub"
	"errors"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (s *Server) handleWebSocket(ctx *gin.Context) {
	roomIDStr := ctx.Request.URL.Query().Get("room_id")
	if roomIDStr == "" {
		s.errorResponse(ctx, http.StatusBadRequest, "not exists room")
		return
	}

	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		s.errorResponse(ctx, http.StatusBadRequest, "invalid room Id")
		return
	}

	authPayload, ok := middleware.GetAuthPayload(ctx)
	if !ok {
		s.errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	room, err := s.store.GetRoomByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.errorResponse(ctx, http.StatusNotFound, "room not found")
			return
		}

		s.errorResponse(ctx, http.StatusInternalServerError, "failed to get room")
		return
	}

	member, err := s.store.GetRoomMemberByUsername(ctx, db.GetRoomMemberByUsernameParams{
		RoomID:   room.ID,
		Username: authPayload.Username(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.errorResponse(ctx, http.StatusForbidden, "not a user of this room")
			return
		}
		s.errorResponse(ctx, http.StatusInternalServerError, "failed server request")
		return
	}

	conn, err := websocket.Accept(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	hubRoom, _ := s.hub.GetOrCreateRoom(room.ID, room.Name)

	client := hub.NewClient(conn, hubRoom, member.UserID, authPayload.Username(), s.store)

	client.Run(ctx)
}
