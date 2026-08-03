package api

import (
	"GoChat/config"
	db "GoChat/db/sqlc"
	"GoChat/middleware"
	"GoChat/servers/hub"
	"GoChat/token"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     config.Config
	router     *gin.Engine
	store      db.Store
	hub        *hub.Hub
	tokenMaker token.Maker
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		hub:        hub.NewHub(),
		tokenMaker: tokenMaker,
	}

	server.setupServer()

	return server, nil
}

func (server *Server) setupServer() {
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "OK",
		})
	})

	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.tokenMaker))

	authRoutes.POST("/rooms", server.handleCreateRoom)
	authRoutes.GET("/rooms", server.handleListRooms)
	authRoutes.POST("/rooms/:id/join", server.handleJoinRoom)
	authRoutes.GET("/ws", server.handleWebSocket)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) errorResponse(gc *gin.Context, statusCode int, message string) {
	gc.JSON(statusCode, gin.H{"error": message})
}
