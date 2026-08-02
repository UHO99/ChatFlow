package api

import (
	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleWebSocket(c *gin.Context) {
	// ctx := c.Request.Context()

	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to accept WebSocket connection"})
		return
	}
	defer conn.CloseNow()

	// client := s.hub.NewClient(conn, ctx)

	// client.Run(ctx)
}
