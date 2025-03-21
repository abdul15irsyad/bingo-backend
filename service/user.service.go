package service

import (
	"bingo/lib"
	"bingo/model"
	"bingo/util"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
)

var UserMutex sync.Mutex

type UserService struct {
	MaxUser int
	Users   []model.User
}

func NewUserService(maxUser int) *UserService {
	lib.Logger.Info("NewUserService initialized")
	return &UserService{
		MaxUser: maxUser,
		Users:   []model.User{},
	}
}

type UserDTO struct {
	Username string
}

func (s *UserService) AddUser(userDTO UserDTO) model.User {
	newUuid, _ := uuid.NewRandom()
	newUser := model.User{
		Id:        newUuid,
		Username:  userDTO.Username,
		CreatedAt: time.Now(),
	}

	UserMutex.Lock()
	s.Users = append(s.Users, newUser)
	UserMutex.Unlock()

	return newUser
}

func (s *UserService) GetUsers() []model.User {
	sortedUsers := s.Users
	slices.SortFunc(sortedUsers, func(a, b model.User) int {
		if a.Username < b.Username {
			return -1
		} else {
			return 1
		}
	})
	return sortedUsers
}

func (s *UserService) GetUser(id uuid.UUID) *model.User {
	user := util.FindSlice(&s.Users, func(user model.User) bool {
		return user.Id == id
	})
	return user
}

func (s *UserService) UpdateUser(id uuid.UUID, userDTO UserDTO) model.User {
	user := util.FindSlice(&s.Users, func(user model.User) bool {
		return user.Id != id
	})
	UserMutex.Lock()
	user.Username = userDTO.Username
	UserMutex.Unlock()

	return *user
}

func (s *UserService) DeleteUser(id uuid.UUID) {
	UserMutex.Lock()
	s.Users = util.FilterSlice(&s.Users, func(user model.User) bool {
		return user.Id != id
	})
	UserMutex.Unlock()
}
