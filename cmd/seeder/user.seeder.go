package main

import (
	"bingo/model"
	"bingo/util"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/bxcodec/faker"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func userSeeder(db *gorm.DB) {
	users := []model.User{
		newUser("Luffy"),
		newUser("Zoro"),
		newUser("Sanji"),
	}

	for range 20 - len(users) {
		randomUuid, _ := uuid.NewRandom()
		hashedPassword, _ := util.HashPassword("Qwerty123")
		name, _ := faker.GetPerson().Name(reflect.Value{})
		nameSlug := util.Slugify(name.(string))
		randomDate := util.RandomDate(time.Now().AddDate(0, 0, -1), time.Now())
		email := strings.ReplaceAll(nameSlug, "-", "") + "@email.com"
		user := model.User{
			Id:              randomUuid,
			Name:            name.(string),
			Username:        &nameSlug,
			Email:           &email,
			EmailVerifiedAt: nil,
			Password:        &hashedPassword,
			CreatedAt:       randomDate,
			UpdatedAt:       randomDate,
		}
		users = append(users, user)
	}

	result := db.Model(&model.User{}).Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Printf("%d users inserted\n", result.RowsAffected)
}

func newUser(name string) model.User {
	newUuid, _ := uuid.NewRandom()
	hashedPassword, _ := util.HashPassword("Qwerty123")
	randomDate := util.RandomDate(time.Now().AddDate(0, 0, -1), time.Now())
	username := util.Slugify(name)
	email := fmt.Sprintf("%s@email.com", username)
	return model.User{
		Id:              newUuid,
		Name:            name,
		Username:        &username,
		Email:           &email,
		EmailVerifiedAt: &randomDate,
		Password:        &hashedPassword,
		CreatedAt:       randomDate,
		UpdatedAt:       randomDate,
	}

}
