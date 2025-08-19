package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"r2-challenge/cmd/envs"
)

type Database struct{ *gorm.DB }

func Setup(e envs.Envs) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		e.DBHost, e.DBPort, e.DBUser, e.DBPassword, e.DBName, e.DBSSLMode,
	)

	g, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}
	sqlDB, err := g.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(e.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(e.DBMaxIdleConns)
	if d, err := time.ParseDuration(e.DBConnMaxLifetime); err == nil {
		sqlDB.SetConnMaxLifetime(d)
	}
	return &Database{g}, nil
}
