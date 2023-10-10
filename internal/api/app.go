package app

import (
	"github.com/joho/godotenv"
	"github.com/lMauseRl/WebAppDev/internal/app/dsn"
	"github.com/lMauseRl/WebAppDev/internal/app/repository"
)

type Application struct {
	repository *repository.Repository
}

func New() (Application, error) {
	_ = godotenv.Load()
	repo, err := repository.New(dsn.SetConnectionString())
	if err != nil {
		return Application{}, err
	}

	return Application{repository: repo}, nil
}
