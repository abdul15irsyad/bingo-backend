package service

import (
	"bingo/data"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AddUserDTO struct {
	Username string
}

func AddUser(addUserDTO AddUserDTO) data.User {
	newUuid, _ := uuid.NewRandom()
	newUser := data.User{
		Id:        newUuid,
		Username:  addUserDTO.Username,
		CreatedAt: time.Now(),
	}

	var mutex sync.Mutex
	mutex.Lock()
	data.Users = append(data.Users, newUser)
	mutex.Unlock()

	return newUser
}
