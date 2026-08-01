package api

import (
	"GoChat/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	router *gin.Engine
}

func NewServer(config config.Config) (*Server, error) {
	server := &Server{
		config: config,
	}

	if err := server.setupServer(config); err != nil {
		return nil, err
	}

	return server, nil
}

func (server *Server) setupServer(config config.Config) error {
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "OK",
		})
	})

	server.router = router

	return nil
}
