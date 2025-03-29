package service

import (
	"bingo/lib"
	"bingo/model"
	"bingo/util"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameService struct {
	Games         []*model.Game
	Mutex         sync.RWMutex
	socketService *SocketService
}

func NewGameService(socketService *SocketService) *GameService {
	lib.Logger.Info("NewGameService initialized")
	return &GameService{[]*model.Game{}, sync.RWMutex{}, socketService}
}

type CreateGameDTO struct {
	TotalPlayer int
	Users       []*model.User
}

// Game
func (s *GameService) CreateGame(dto CreateGameDTO) (model.Game, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	newUuid, _ := uuid.NewRandom()
	var newCode string
	for {
		randomString, err := util.RandomString(6, nil)
		if err != nil {
			return model.Game{}, err
		}
		if util.FindSlice(&s.Games, func(g **model.Game) bool {
			return (*g).Code == randomString
		}) == nil {
			newCode = randomString
			break
		}
	}

	players := util.MapSlice(&dto.Users, func(user *model.User) *model.Player {
		newUuid, _ := uuid.NewRandom()
		return &model.Player{
			Id:     newUuid,
			User:   user,
			Status: false,
		}
	})

	newGame := model.Game{
		Id:          newUuid,
		Code:        newCode,
		TotalPlayer: dto.TotalPlayer,
		Players:     players,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		GameTiles: util.MapSlice(&players, func(p *model.Player) *model.GameTile {
			newUuid, _ := uuid.NewRandom()
			return &model.GameTile{
				Id:       newUuid,
				Tiles:    createTiles(),
				PlayerId: p.Id,
			}
		}),
		MarkedTiles: []*model.MarkedTile{},
	}

	return newGame, nil
}

func (s *GameService) GetGame(id uuid.UUID) *model.Game {
	return *util.FindSlice(&s.Games, func(g **model.Game) bool {
		return (*g).Id == id
	})
}

// Player
func (s *GameService) GetPlayerFromUserId(game *model.Game, userId uuid.UUID) *model.Player {
	return *util.FindSlice(&game.Players, func(p **model.Player) bool {
		return (*p).User.Id == userId
	})
}

func (s *GameService) QueuePlayer(client *model.Client, totalPlayer int) error {
	s.socketService.Mutex.Lock()
	queuePlayers := util.FilterSlice(&s.socketService.Queues, func(queue **model.Queue) bool {
		return (*queue).GameTotalPlayer == totalPlayer && (*queue).Client.User.Id != client.User.Id
	})
	if len(queuePlayers)+1 < totalPlayer {
		// add to queue
		newUuid, _ := uuid.NewRandom()
		s.socketService.Queues = append(s.socketService.Queues, &model.Queue{
			Id:              newUuid,
			GameTotalPlayer: totalPlayer,
			Client:          client,
			CreatedAt:       time.Now(),
		})
	} else {
		// start game
		newUuid, _ := uuid.NewRandom()
		queuePlayers = append(queuePlayers, &model.Queue{
			Id:              newUuid,
			GameTotalPlayer: totalPlayer,
			Client:          client,
			CreatedAt:       time.Now(),
		})
		// remove from queue
		s.socketService.Queues = util.FilterSlice(&s.socketService.Queues, func(queue **model.Queue) bool {
			return (*queue).GameTotalPlayer != totalPlayer && (*queue).Client.User.Id != client.User.Id
		})

		game, err := s.CreateGame(CreateGameDTO{
			TotalPlayer: totalPlayer,
			Users: util.MapSlice(&queuePlayers, func(p *model.Queue) *model.User {
				return (*p).Client.User
			}),
		})
		if err != nil {
			return err
		}
		s.Games = append(s.Games, &game)

		room := s.socketService.CreateRoom(&game)
		room.Clients = util.MapSlice(&queuePlayers, func(queue *model.Queue) *model.Client {
			return queue.Client
		})

		err = s.socketService.BroadcastToRoom(room, model.Payload{
			Type:      model.GameMatchType,
			Content:   game.Id,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
	}
	s.socketService.Mutex.Unlock()

	return nil
}

func (s *GameService) CheckIsAllReady(game *model.Game) (isAllReady bool) {
	return util.FindSlice(&game.Players, func(p **model.Player) bool {
		return !(*p).Status
	}) == nil
}

func createTiles() []*model.Tile {
	var tiles []*model.Tile
	size := 5
	for y := range size {
		for x := range size {
			newUuid, _ := uuid.NewRandom()
			tiles = append(tiles, &model.Tile{
				Id:       newUuid,
				X:        x,
				Y:        y,
				IsMarked: false,
				Number:   size*y + (x + 1),
			})
		}
	}

	return tiles
}

func (s *GameService) MarkTile(game *model.Game, number int) {
	for _, gameTile := range game.GameTiles {
	playerTiles:
		for _, tile := range gameTile.Tiles {
			if tile.Number == number {
				tile.IsMarked = true
				break playerTiles
			}
		}
	}

	// add to marked tile
	newUuid, _ := uuid.NewRandom()
	game.MarkedTiles = append(game.MarkedTiles, &model.MarkedTile{
		Id:     newUuid,
		Order:  len(game.MarkedTiles) + 1,
		Number: number,
	})
}
