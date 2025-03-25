package service

import (
	"bingo/model"
	"bingo/util"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id   uuid.UUID `json:"id"`
	Conn *websocket.Conn
	User *model.User `json:"user"`
	Send chan []byte
}

type Room struct {
	Id      uuid.UUID `json:"id"`
	Clients []Client
	Game    *model.Game
	Mutex   sync.RWMutex
}

type Message struct {
	Type      string    `json:"type"`
	Client    *Client   `json:"client"`
	Content   []byte    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type BroadcastRoom struct {
	Room    *Room
	Message Message
}

type AddClientToRoom struct {
	Room   *Room
	Client *Client
}

type RemoveClientFromRoom = AddClientToRoom

type SocketService struct {
	Clients                  []Client
	Rooms                    []Room
	MaxRoom                  int
	AddClientToRoomChan      chan *AddClientToRoom
	RemoveClientFromRoomChan chan *RemoveClientFromRoom
	BroadcastRoomChan        chan *BroadcastRoom
	Mutex                    sync.RWMutex
}

func NewSocketService(maxRoom int) *SocketService {
	return &SocketService{
		Clients:                  []Client{},
		Rooms:                    []Room{},
		MaxRoom:                  maxRoom,
		AddClientToRoomChan:      make(chan *AddClientToRoom),
		RemoveClientFromRoomChan: make(chan *RemoveClientFromRoom),
		BroadcastRoomChan:        make(chan *BroadcastRoom),
	}
}

func (s *SocketService) Run() {
	for {
		select {
		case register := <-s.AddClientToRoomChan:
			s.AddClientToRoom(register.Room, register.Client)
		case unregister := <-s.RemoveClientFromRoomChan:
			s.RemoveClientFromRoom(unregister.Room, unregister.Client)
		case broadcastRoom := <-s.BroadcastRoomChan:
			s.BroadcastToRoom(broadcastRoom.Room, broadcastRoom.Message)
		}
	}
}

func (s *SocketService) CreateClient(conn *websocket.Conn, user *model.User) *Client {
	id, _ := uuid.NewRandom()
	return &Client{
		Id:   id,
		User: user,
		Send: make(chan []byte),
		Conn: conn,
	}
}

func (s *SocketService) RemoveClient(client *Client) {
	s.Clients = util.FilterSlice(&s.Clients, func(c *Client) bool {
		return c.Id != client.Id
	})
	client.Conn.Close()
}

func (s *SocketService) CreateRoom(game *model.Game) *Room {
	id, _ := uuid.NewRandom()
	return &Room{
		Id:      id,
		Clients: []Client{},
		Game:    game,
	}
}

func (s *SocketService) AddClientToRoom(room *Room, client *Client) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room.Mutex.Lock()
	room.Clients = append(room.Clients, *client)
	room.Mutex.Unlock()
}

func (s *SocketService) RemoveClientFromRoom(room *Room, client *Client) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room.Mutex.Lock()
	room.Clients = util.FilterSlice(&room.Clients, func(roomClient *Client) bool {
		return roomClient.Id != client.Id
	})
	close(client.Send)
	room.Mutex.Unlock()

	if len(room.Clients) == 0 {
		s.Rooms = util.FilterSlice(&s.Rooms, func(r *Room) bool {
			return r.Id != room.Id
		})
	}
}

func (s *SocketService) BroadcastToRoom(room *Room, message Message) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room.Mutex.RLock()
	defer room.Mutex.RUnlock()

	for _, client := range room.Clients {
		select {
		case client.Send <- message.Content:
		default:
			room.Clients = util.FilterSlice(&room.Clients, func(roomClient *Client) bool {
				return roomClient.Id != client.Id
			})
			close(client.Send)
		}
	}
}
