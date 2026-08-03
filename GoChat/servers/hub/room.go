package hub

import "sync"

type Room struct {
	mu      sync.RWMutex
	id      int64
	name    string
	clients map[*Client]struct{}
}

func NewRoom(id int64, name string) *Room {
	return &Room{
		id:      id,
		name:    name,
		clients: make(map[*Client]struct{}),
	}
}

func (r *Room) addClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

func (r *Room) removeClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
		close(c.send)
	}
}

func (r *Room) broadcast(msg []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for c := range r.clients {
		select {
		case c.send <- msg:
		default:
		}
	}
}

func (r *Room) isEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients) == 0
}
