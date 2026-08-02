package hub

import (
	"context"
	"sync"
)

type Hub struct {
	mu    sync.Mutex
	rooms map[int64]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[int64]*Room),
	}
}

func (h *Hub) GetRoom(ctx context.Context, roomID int64, name string) (*Room, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		return room, nil
	}

	room := NewRoom(roomID, name)
	h.rooms[roomID] = room

	return room, nil
}

func (h *Hub) RemoveRoom(ctx context.Context, roomID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok && room.IsEmpty() {
		delete(h.rooms, roomID)
	}
}
