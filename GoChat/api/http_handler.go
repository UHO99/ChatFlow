package api

import (
	db "GoChat/db/sqlc"
	"GoChat/middleware"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type roomResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type createRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) handleCreateRoom(ctx *gin.Context) {
	var req createRoomRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		s.errorResponse(ctx, http.StatusBadRequest, "Request Parameter Error")
		return
	}

	room, err := s.store.CreateRoom(ctx, db.CreateRoomParams{
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			s.errorResponse(ctx, http.StatusConflict, "room aleary exists")
			return
		}

		s.errorResponse(ctx, http.StatusInternalServerError, "cannot create room")
		return
	}

	authPayload, ok := middleware.GetAuthPayload(ctx)
	if !ok {
		s.errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	if _, err := s.store.JoinRoomByUsername(ctx, db.JoinRoomByUsernameParams{
		RoomID:   room.ID,
		Username: authPayload.Username(),
	}); err != nil {
		s.errorResponse(ctx, http.StatusInternalServerError, "cannot add room creator as member")
		return
	}

	ctx.JSON(http.StatusCreated, roomResponse{ID: room.ID, Name: room.Name})
}

func (s *Server) handleJoinRoom(ctx *gin.Context) {
	roomID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		s.errorResponse(ctx, http.StatusBadRequest, "invalid room id")
		return
	}

	if _, err := s.store.GetRoomByID(ctx, roomID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.errorResponse(ctx, http.StatusNotFound, "room not found")
			return
		}
		s.errorResponse(ctx, http.StatusInternalServerError, "failed to get room")
		return
	}

	authPayload, ok := middleware.GetAuthPayload(ctx)
	if !ok {
		s.errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	if _, err := s.store.JoinRoomByUsername(ctx, db.JoinRoomByUsernameParams{
		RoomID:   roomID,
		Username: authPayload.Username(),
	}); err != nil {
		s.errorResponse(ctx, http.StatusInternalServerError, "cannot join room")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "joined room"})
}

func (s *Server) handleListRooms(ctx *gin.Context) {
	rooms, err := s.store.GetListRoom(ctx)
	if err != nil {
		s.errorResponse(ctx, http.StatusInternalServerError, "cannot list rooms")
		return
	}

	res := make([]roomResponse, 0, len(rooms))
	for _, room := range rooms {
		res = append(res, roomResponse{ID: room.ID, Name: room.Name})
	}

	ctx.JSON(http.StatusOK, res)
}
