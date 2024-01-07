package app

import (
	"context"

	"github.com/lud0m4n/WebAppDev/internal/config"
	"github.com/lud0m4n/WebAppDev/internal/dsn"
	"github.com/lud0m4n/WebAppDev/internal/http/fossil"
	"github.com/lud0m4n/WebAppDev/internal/http/repository"
	"github.com/lud0m4n/WebAppDev/internal/http/usecase"
)

// Application представляет основное приложение.
type Application struct {
	Config     *config.Config
	Repository *repository.Repository
	UseCase    *usecase.UseCase
	Handler    *fossil.Handler
}

// New создает новый объект Application и настраивает его.
func New(ctx context.Context) (*Application, error) {
	// Инициализируйте конфигурацию
	cfg, err := config.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Инициализируйте подключение к базе данных (DB)
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}
	uc := usecase.NewUseCase(repo)
	h := fossil.NewHandler(uc)
	// Инициализируйте и настройте объект Application
	app := &Application{
		Config:     cfg,
		Repository: repo,
		UseCase:    uc,
		Handler:    h,
	}

	return app, nil
}
