package hub

import (
	db "GoChat/db/sqlc"
	"GoChat/servers/messages"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

const (
	sendBufferSize = 256
	writeTimeout   = 5 * time.Second
)

type Client struct {
	conn     *websocket.Conn
	room     *Room
	userID   int64
	nickname string
	send     chan []byte
	store    db.Store
}

func NewClient(conn *websocket.Conn, room *Room, userID int64, nickname string, store db.Store) *Client {
	return &Client{
		conn:     conn,
		room:     room,
		userID:   userID,
		nickname: nickname,
		send:     make(chan []byte, sendBufferSize),
		store:    store,
	}
}

func (c *Client) Run(ctx *gin.Context) {
	c.room.addClient(c)
	c.room.broadcast(c.event(messages.TypeJoin, ""))

	done := make(chan struct{})
	go c.writePump(ctx, done)
	c.readPump(ctx)

	c.room.removeClient(c)
	<-done

	c.room.broadcast(c.event(messages.TypeLeave, ""))
}

func (c *Client) readPump(ctx context.Context) {
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			switch {
			case errors.Is(err, io.EOF), errors.Is(err, net.ErrClosed):
				log.Printf("client %s disconnected", c.nickname)
			case websocket.CloseStatus(err) == websocket.StatusNormalClosure:
				log.Printf("client %s closed the connection", c.nickname)
			default:
				log.Printf("error reading from client %s: %v", c.nickname, err)
			}

			return
		}

		var in messages.IncomingMessage
		if err := json.Unmarshal(data, &in); err != nil {
			log.Printf("error unmarshaling message from client %s: %v", c.nickname, err)
			continue
		}

		if err := c.store.CreateMessage(ctx, db.CreateMessageParams{
			RoomID:  c.room.id,
			UserID:  c.userID,
			Content: in.Content,
		}); err != nil {
			log.Printf("error creating message for client %s: %v", c.nickname, err)
			continue
		}

		c.room.broadcast(c.event(messages.TypeMessage, in.Content))
	}
}

func (c *Client) writePump(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	for payload := range c.send {
		if err := c.write(ctx, payload); err != nil {
			log.Printf("cannot write to client %s: %v", c.nickname, err)
			return
		}
	}
}

func (c *Client) write(ctx context.Context, payload []byte) error {
	writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	return c.conn.Write(writeCtx, websocket.MessageText, payload)
}

func (c *Client) event(eventType messages.Type, content string) []byte {
	payload, err := json.Marshal(&messages.OutgoingMessage{
		Type:    eventType,
		RoomID:  c.room.id,
		Content: content,
		Sender:  c.nickname,
		SendAt:  time.Now(),
	})
	if err != nil {
		log.Printf("error marshaling event for client %s: %v", c.nickname, err)
		return nil
	}

	return payload
}
