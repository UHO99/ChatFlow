package api

import (
	db "GoChat/db/sqlc"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type roomResponse struct {
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

	ctx.JSON(http.StatusCreated, roomResponse{Name: room.Name})
}

func (s *Server) handleListRooms(ctx *gin.Context) {
	rooms, err := s.store.GetListRoom(ctx)
	if err != nil {
		s.errorResponse(ctx, http.StatusInternalServerError, "cannot list rooms")
		return
	}

	res := make([]roomResponse, 0, len(rooms))
	for _, room := range rooms {
		res = append(res, roomResponse{Name: room.Name})
	}

	ctx.JSON(http.StatusOK, res)
}
