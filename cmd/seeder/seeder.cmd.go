package main

import (
	"bingo/config"
	"bingo/lib"
	"bingo/model"
	"bingo/util"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/bxcodec/faker"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	dbManager := lib.NewGormManager()
	db, err := dbManager.InitPostgresDB("seeder", config.DBConfig)
	if err != nil {
		panic(err)
	}

	for _, seeder := range []struct {
		name string
		seed func(*gorm.DB)
	}{
		{name: "userSeeder", seed: userSeeder},
	} {
		// check name in seeders table
		result := db.Model(&model.Seeder{}).Where("name = ?", seeder.name).First(&seeder)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			panic(result.Error)
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Printf("%s already executed\n", seeder.name)
			continue
		}

		fmt.Printf("executing %s...\n", seeder.name)
		seeder.seed(db)

		// add name to seeders table
		Id, _ := uuid.NewRandom()
		db.Create(&model.Seeder{
			Id:        Id,
			Name:      seeder.name,
			CreatedAt: time.Now(),
		})
		fmt.Printf("%s executed\n", seeder.name)
	}
}

func userSeeder(db *gorm.DB) {
	users := []model.User{}
	newUuid, _ := uuid.Parse("0e3cec4a-206c-4dfb-96f6-3f6b85db9543")
	hashedPassword, _ := util.HashPassword("Qwerty123")
	randomDate := util.RandomDate(time.Now().AddDate(0, 0, -1), time.Now())
	username := "irsyadabdul"
	email := "abdulirsyad@email.com"
	user := model.User{
		Id:              newUuid,
		Name:            "Irsyad Abdul",
		Username:        &username,
		Email:           &email,
		EmailVerifiedAt: &randomDate,
		Password:        &hashedPassword,
		CreatedAt:       randomDate,
		UpdatedAt:       randomDate,
	}
	users = append(users, user)

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
