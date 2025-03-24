package lib

import (
	"bingo/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormManager struct {
	DBs map[string]*gorm.DB
}

func NewGormManager() *GormManager {
	return &GormManager{
		map[string]*gorm.DB{},
	}
}

func (m *GormManager) InitPostgresDB(name string, databaseConfig config.DatabaseConfig) (*gorm.DB, error) {
	gormLogger := logger.Default.LogMode(logger.Info)
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", databaseConfig.User, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port, databaseConfig.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetConnMaxLifetime(time.Duration(databaseConfig.ConnMaxLifetime) * time.Second)

	m.DBs[name] = db

	return db, nil
}

func (m *GormManager) GetDatabase(name string) *gorm.DB {
	if _, ok := m.DBs[name]; !ok {
		return nil
	}

	sqlDB, err := m.DBs[name].DB()
	if err != nil {
		return nil
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil
	}

	return m.DBs[name]
}

func (m *GormManager) CloseDatabase(name string) error {
	if _, ok := m.DBs[name]; !ok {
		return nil
	}

	sqlDB, err := m.DBs[name].DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		return err
	}

	delete(m.DBs, name)
	return nil
}
