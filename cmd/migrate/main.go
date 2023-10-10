package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lMauseRl/WebAppDev/internal/app/ds"
	"github.com/lMauseRl/WebAppDev/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.SetConnectionString()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&ds.Periods{})
	if err != nil {
		panic("cant migrate db")
	}
}
