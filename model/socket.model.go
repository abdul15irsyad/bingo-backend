package model

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

	Queue struct {
		Id              uuid.UUID `json:"id"`
		GameTotalPlayer int       `json:"game_type"`
		Client          *Client   `json:"client"`
		CreatedAt       time.Time `json:"created_at"`
	}

	Message struct {
		Type      string    `json:"type"`
		User      *User     `json:"user"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	BroadcastRoom struct {
		Room    *Room
		Message Message
	}

	AddClientToRoom struct {
		Room   *Room
		Client *Client
	}

	RemoveClientFromRoom = AddClientToRoom
)
