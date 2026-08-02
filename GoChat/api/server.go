package api

import (
	"GoChat/config"
	db "GoChat/db/sqlc"
	"GoChat/servers/hub"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	router *gin.Engine
	store  db.Store
	hub    *hub.Hub
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
		hub:    hub.NewHub(),
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

	router.GET("/ws", server.handleWebSocket)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
