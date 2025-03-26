package service

import (
	"bingo/model"
	"bingo/util"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SocketService struct {
	Clients []model.Client
	Queues  []model.Queue
	Rooms   []model.Room
	MaxRoom int
	// AddClientToRoomChan      chan *model.AddClientToRoom
	// RemoveClientFromRoomChan chan *model.RemoveClientFromRoom
	// BroadcastRoomChan        chan *model.BroadcastRoom
	Mutex sync.RWMutex
}

func NewSocketService(maxRoom int) *SocketService {
	return &SocketService{
		Clients: []model.Client{},
		Queues:  []model.Queue{},
		Rooms:   []model.Room{},
		MaxRoom: maxRoom,
		// AddClientToRoomChan:      make(chan *model.AddClientToRoom),
		// RemoveClientFromRoomChan: make(chan *model.RemoveClientFromRoom),
		// BroadcastRoomChan:        make(chan *model.BroadcastRoom),
	}
}

// func (s *SocketService) Run() {
// 	for {
// 		select {
// 		case register := <-s.AddClientToRoomChan:
// 			s.AddClientToRoom(register.Room, register.Client)
// 		case unregister := <-s.RemoveClientFromRoomChan:
// 			s.RemoveClientFromRoom(unregister.Room, unregister.Client)
// 		case broadcastRoom := <-s.BroadcastRoomChan:
// 			s.BroadcastToRoom(broadcastRoom.Room, broadcastRoom.Message)
// 		}
// 	}
// }

func (s *SocketService) CreateClient(conn *websocket.Conn, user *model.User) *model.Client {
	id, _ := uuid.NewRandom()
	return &model.Client{
		Id:   id,
		User: user,
		Conn: conn,
	}
}

func (s *SocketService) RemoveClient(client *model.Client) {
	s.Clients = util.FilterSlice(&s.Clients, func(c *model.Client) bool {
		return c.Id != client.Id
	})
}

func (s *SocketService) SendMessage(client *model.Client, message *model.Message) error {
	err := client.Conn.WriteJSON(*message)
	if err != nil {
		return err
	}

	return nil
}

func (s *SocketService) CreateRoom(game *model.Game) *model.Room {
	id, _ := uuid.NewRandom()
	return &model.Room{
		Id:      id,
		Clients: []model.Client{},
		Game:    game,
	}
}

func (s *SocketService) AddClientToRoom(room *model.Room, client *model.Client) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room.Mutex.Lock()
	room.Clients = append(room.Clients, *client)
	room.Mutex.Unlock()
}

func (s *SocketService) RemoveClientFromRoom(room *model.Room, client *model.Client) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room.Mutex.Lock()
	room.Clients = util.FilterSlice(&room.Clients, func(roomClient *model.Client) bool {
		return roomClient.Id != client.Id
	})
	room.Mutex.Unlock()

	if len(room.Clients) == 0 {
		s.Rooms = util.FilterSlice(&s.Rooms, func(r *model.Room) bool {
			return r.Id != room.Id
		})
	}
}

func (s *SocketService) BroadcastToRoom(room *model.Room, message model.Message) error {
	room.Mutex.RLock()
	defer room.Mutex.RUnlock()

	for _, client := range room.Clients {
		err := s.SendMessage(&client, &message)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
