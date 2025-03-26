package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id              uuid.UUID  `json:"id" gorm:"column:id;type:varchar(40);primaryKey"`
	Name            string     `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Username        string     `json:"username" gorm:"column:username;unique;type:varchar(255);not null"`
	Email           string     `json:"email" gorm:"column:email;unique;type:varchar(255);not null"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"column:email_verified_at;type:timestamptz"`
	Password        string     `json:"-" gorm:"select:false;column:password;type:varchar(255);not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at;type:timestamptz"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamptz"`
}
