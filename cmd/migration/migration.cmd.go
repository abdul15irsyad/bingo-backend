package main

import (
	"bingo/config"
	"bingo/lib"
	"bingo/model"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	dbManager := lib.NewGormManager()
	db, err := dbManager.InitPostgresDB("migration", config.DBConfig)
	if err != nil {
		panic(err)
	}

	// start migrate
	if !db.Migrator().HasTable(&model.User{}) {
		db.Migrator().CreateTable(&model.User{})
	}
	if !db.Migrator().HasTable(&model.Seeder{}) {
		db.Migrator().CreateTable(&model.Seeder{})
	}
}
