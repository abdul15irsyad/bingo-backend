package service

import (
	"bingo/model"
	"bingo/util"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id   uuid.UUID  `json:"id"`
	User model.User `json:"user"`
	Send chan []byte
	Conn *websocket.Conn
}

type Room struct {
	Id      uuid.UUID `json:"id"`
	Code    string    `json:"code"`
	Type    string    `json:"type"`
	Clients []Client
	Mutex   sync.RWMutex
}

type RoomMessage struct {
	Type      string    `json:"type"`
	Client    *Client   `json:"client"`
	Content   []byte    `json:"content"`
	RoomId    uuid.UUID `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
}

type HubService struct {
	Clients       []Client
	Rooms         []Room
	Register      chan *Client
	Unregister    chan *Client
	BroadcastRoom chan RoomMessage
	Mutex         sync.RWMutex
}

func NewHubService() *HubService {
	return &HubService{
		Clients:       []Client{},
		Rooms:         []Room{},
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		BroadcastRoom: make(chan RoomMessage),
	}
}

func (s *HubService) Run() {
	for {
		select {
		// case client := <-s.Register:
		// 	s.registerClient(client)
		// case client := <-s.Unregister:
		// 	s.unregisterClient(client)
		case message := <-s.BroadcastRoom:
			s.BroadcastToRoom(message)
		}
	}
}

func (h *HubService) BroadcastToRoom(message RoomMessage) {
	h.Mutex.RLock()
	defer h.Mutex.RUnlock()

	room := util.FindSlice(&h.Rooms, func(room *Room) bool {
		return room.Id == message.RoomId
	})
	if room == nil {
		fmt.Println("room not found")
		return
	}

	room.Mutex.RLock()
	defer room.Mutex.RUnlock()

	for _, client := range room.Clients {
		select {
		case client.Send <- message.Content:
		default:
			close(client.Send)
			room.Clients = util.FilterSlice(&room.Clients, func(roomClient *Client) bool {
				return roomClient.Id == client.Id
			})
		}
	}
}
