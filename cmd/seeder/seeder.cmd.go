package main

import (
	"bingo/cmd/seeder/seeders"
	"bingo/config"
	"bingo/lib"
	"bingo/model"
	"errors"
	"fmt"
	"time"

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
		{name: "userSeeder", seed: seeders.UserSeeder},
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
