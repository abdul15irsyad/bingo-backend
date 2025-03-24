package service

import (
	"bingo/lib"
	"bingo/model"
	"bingo/util"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	lib.Logger.Info("NewUserService initialized")
	return &UserService{
		db: db,
	}
}

type CreateUserDTO struct {
	Name     string
	Username string
	Email    string
	Password string
}

func (s *UserService) CreateUser(dto CreateUserDTO) (model.User, error) {
	newUuid, _ := uuid.NewRandom()
	hashPassword, err := util.HashPassword(dto.Password)
	if err != nil {
		return model.User{}, err
	}

	newUser := model.User{
		Id:        newUuid,
		Name:      dto.Name,
		Username:  dto.Username,
		Email:     dto.Email,
		Password:  hashPassword,
		CreatedAt: time.Now(),
	}

	if err := s.db.Model(&model.User{}).Create(newUser).Error; err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

type GetPaginatedUsersDto struct {
	Page   int
	Limit  int
	Search *string
}

func (s *UserService) GetPaginatedUsers(dto GetPaginatedUsersDto) ([]model.User, int64, error) {
	users := []model.User{}
	offset := (dto.Page - 1) * dto.Limit
	query := s.db.Model(&model.User{})
	if dto.Search != nil && *dto.Search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+*dto.Search+"%", "%"+*dto.Search+"%")
	}
	if err := query.Limit(dto.Limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *UserService) GetUser(id uuid.UUID) (model.User, error) {
	var user model.User
	if err := s.db.Model(&model.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserByUsernameOrEmail(usernameOrEmail string) (model.User, error) {
	var user model.User
	if err := s.db.Model(&model.User{}).Select("id", "password").Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

type UpdateUserDTO struct {
	Name     string
	Username string
	Email    string
}

func (s *UserService) UpdateUser(id uuid.UUID, dto UpdateUserDTO) (model.User, error) {
	var user model.User
	if err := s.db.Model(&user).Where("id = ?", id).Updates(model.User{
		Name:      dto.Name,
		Username:  dto.Username,
		Email:     dto.Email,
		UpdatedAt: time.Now(),
	}).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	if err := s.db.Delete(&model.User{}).Where("id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
