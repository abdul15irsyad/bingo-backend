package model

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	GameMatchType   = "game match"
	PlayerReadyType = "player ready"
	GameStartType   = "game start"
	MessageType     = "message"
)

type (
	Client struct {
		Id   uuid.UUID `json:"id"`
		Conn *websocket.Conn
		User *User `json:"user"`
		// Send chan []byte
	}

	Room struct {
		Id      uuid.UUID `json:"id"`
		Clients []Client
		Game    *Game
		Mutex   sync.RWMutex
	}

	Payload struct {
		Type      string    `json:"type"`
		User      *User     `json:"user"`
		Content   any       `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	MoveContent struct {
		Number int   `json:"number"`
		User   *User `json:""`
	}

	BroadcastRoom struct {
		Room    *Room
		Message Payload
	}

	AddClientToRoom struct {
		Room   *Room
		Client *Client
	}

	RemoveClientFromRoom = AddClientToRoom
)
