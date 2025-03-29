package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	Queue struct {
		Id              uuid.UUID `json:"id"`
		GameTotalPlayer int       `json:"game_type"`
		Client          *Client   `json:"client"`
		CreatedAt       time.Time `json:"created_at"`
	}

	Game struct {
		Id           uuid.UUID  `json:"id" gorm:"column:id;type:varchar(40);primaryKey"`
		Code         string     `json:"code" gorm:"column:code;type:varchar(40);not null"`
		TotalPlayer  int        `json:"total_player" gorm:"column:total_player;type:int;not null"`
		CreatedAt    time.Time  `json:"-" gorm:"column:created_at;type:timestamptz"`
		UpdatedAt    time.Time  `json:"-" gorm:"column:updated_at;type:timestamptz"`
		PlayerTurnId *uuid.UUID `json:"player_turn_id"`
		StartAt      *time.Time `json:"start_at"`
		// WinnerId  *uuid.UUID `json:"winner_id" gorm:"column:winner_id;type:varchar(40)"`

		PlayerWinnerId *uuid.UUID    `json:"player_winner_id"`
		Players        []*Player     `json:"players"`
		GameTiles      []*GameTile   `json:"game_tiles"`
		MarkedTiles    []*MarkedTile `json:"marked_tiles"`
	}

	Player struct {
		Id     uuid.UUID `json:"id"`
		User   *User     `json:"user"`
		Status bool      `json:"status"`
	}

	GameTile struct {
		Id    uuid.UUID `json:"id" gorm:"column:id;type:varchar(40);primaryKey"`
		Tiles []*Tile
		// GameId uuid.UUID `json:"game_id" gorm:"column:game_id;type:varchar(40);not null"`
		// UserId uuid.UUID `json:"user_id" gorm:"column:user;type:varchar(40);not null"`

		PlayerId uuid.UUID `json:"player_id"`
	}

	MarkedTile struct {
		Id     uuid.UUID `json:"id" gorm:"column:id;type:varchar(40);primaryKey"`
		Order  int       `json:"order" gorm:"column:order;type:int;not null"`
		Number int       `json:"number"`
		// CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamptz"`
		// UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamptz"`
	}

	Tile struct {
		Id       uuid.UUID `json:"id" gorm:"column:id;type:varchar(40);primaryKey"`
		X        int       `json:"x" gorm:"column:x;type:int"`
		Y        int       `json:"y" gorm:"column:y;type:int"`
		IsMarked bool      `json:"is_marked" gorm:"column:is_marked;type:bool;default:false"`
		Number   int       `json:"number"`
		// GameTileId uuid.UUID `json:"game_tile_id" gorm:"column:game_tile_id;type:varchar(40);not null"`
		// CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamptz"`
		// UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamptz"`

		// GameTile GameTile `json:"game_tile"`
	}
)
