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
	Games []model.Game
	Mutex sync.RWMutex
}

func NewGameService() *GameService {
	lib.Logger.Info("NewGameService initialized")
	return &GameService{[]model.Game{}, sync.RWMutex{}}
}

type CreateGameDTO struct {
	Type  string
	Users []model.User
}

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
		if util.FindSlice(&s.Games, func(g *model.Game) bool {
			return g.Code == randomString
		}) == nil {
			newCode = randomString
			break
		}
	}

	newGame := model.Game{
		Id:        newUuid,
		Code:      newCode,
		Type:      dto.Type,
		Users:     dto.Users,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		GameTiles: util.MapSlice(dto.Users, func(u model.User) model.GameTile {
			newUuid, _ := uuid.NewRandom()
			return model.GameTile{
				Id:    newUuid,
				Tiles: CreateTiles(),
			}
		}),
		MarkedTiles: []model.MarkedTile{},
	}

	return newGame, nil
}

func CreateTiles() []model.Tile {
	var tiles []model.Tile
	size := 5
	for y := range size {
		for x := range size {
			newUuid, _ := uuid.NewRandom()
			tiles = append(tiles, model.Tile{
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

func MarkTile(number int, game *model.Game) {
	for i := range game.GameTiles {
		for j := range game.GameTiles[i].Tiles {
			if game.GameTiles[i].Tiles[j].Number == number {
				game.GameTiles[i].Tiles[j].IsMarked = true
				break
			}
		}
	}

	// add to marked tile
	newUuid, _ := uuid.NewRandom()
	game.MarkedTiles = append(game.MarkedTiles, model.MarkedTile{
		Id:     newUuid,
		Order:  len(game.MarkedTiles) + 1,
		Number: number,
	})
}
