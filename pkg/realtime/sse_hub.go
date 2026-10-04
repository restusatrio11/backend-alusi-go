package realtime

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// StatusEvent represents a live status change or snapshot event for an application
type StatusEvent struct {
	Type           string    `json:"type"` // "snapshot", "status_change", "probe_result"
	AppID          int       `json:"app_id,omitempty"`
	AppSlug        string    `json:"app_slug,omitempty"`
	Nama           string    `json:"nama,omitempty"`
	StatusLayanan  string    `json:"status_layanan,omitempty"`
	PreviousStatus string    `json:"previous_status,omitempty"`
	ResponseTimeMs int       `json:"response_time_ms,omitempty"`
	StatusCode     int       `json:"status_code,omitempty"`
	Message        string    `json:"message,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}

// Client represents a connected SSE client
type Client struct {
	ID       string
	SendChan chan []byte
}

// SSEHub manages all active SSE client streams and broadcasts events
type SSEHub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
	done       chan struct{}
}

// NewSSEHub creates a new SSEHub instance
func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client, 32),
		unregister: make(chan *Client, 32),
		broadcast:  make(chan []byte, 128),
		done:       make(chan struct{}),
	}
}

// Start launches the hub event loop in a goroutine
func (h *SSEHub) Start() {
	go func() {
		for {
			select {
			case <-h.done:
				return
			case client := <-h.register:
				h.mu.Lock()
				h.clients[client.ID] = client
				h.mu.Unlock()
				log.Debug().Str("client_id", client.ID).Int("total_clients", len(h.clients)).Msg("SSE client connected")

			case client := <-h.unregister:
				h.mu.Lock()
				if _, ok := h.clients[client.ID]; ok {
					delete(h.clients, client.ID)
					close(client.SendChan)
					log.Debug().Str("client_id", client.ID).Int("total_clients", len(h.clients)).Msg("SSE client disconnected")
				}
				h.mu.Unlock()

			case message := <-h.broadcast:
				h.mu.RLock()
				for _, client := range h.clients {
					select {
					case client.SendChan <- message:
					default:
						// Non-blocking write: if client buffer is full, drop to prevent blocking others
						log.Warn().Str("client_id", client.ID).Msg("SSE client buffer full, dropping message")
					}
				}
				h.mu.RUnlock()
			}
		}
	}()
}

// Stop terminates the hub event loop
func (h *SSEHub) Stop() {
	close(h.done)
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, client := range h.clients {
		close(client.SendChan)
	}
	h.clients = make(map[string]*Client)
}

// RegisterClient registers a new SSE client
func (h *SSEHub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient removes an SSE client
func (h *SSEHub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastStatusEvent sends a formatted SSE event to all connected clients
func (h *SSEHub) BroadcastStatusEvent(event StatusEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize SSE status event")
		return
	}

	formattedMsg := []byte(fmt.Sprintf("event: %s\ndata: %s\n\n", event.Type, string(data)))
	select {
	case h.broadcast <- formattedMsg:
	default:
		log.Warn().Msg("SSE broadcast channel full, skipping message")
	}
}

// ActiveClientsCount returns number of connected clients
func (h *SSEHub) ActiveClientsCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
