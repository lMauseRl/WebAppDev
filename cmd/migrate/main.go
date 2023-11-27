package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lud0m4n/WebAppDev/internal/app/ds"
	"github.com/lud0m4n/WebAppDev/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Явно мигрировать только нужные таблицы
	err = db.AutoMigrate(&ds.Period{}, &ds.Fossil{}, &ds.User{}, &ds.FossilPeriod{})
	if err != nil {
		panic("cant migrate db")
	}
}
