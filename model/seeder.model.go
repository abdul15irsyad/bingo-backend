package model

import (
	"time"

	"github.com/google/uuid"
)

type Seeder struct {
	Id        uuid.UUID `gorm:"column:id;varchar(40);primaryKey"`
	Name      string    `gorm:"column:name;varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:created_at;timestamptz;not null"`
}
