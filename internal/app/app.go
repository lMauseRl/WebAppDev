package app

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/WebAppDev/internal/api"
	"github.com/lud0m4n/WebAppDev/internal/app/config"
	"github.com/lud0m4n/WebAppDev/internal/app/dsn"
	"github.com/lud0m4n/WebAppDev/internal/app/repository"
)

// Application представляет основное приложение.
type Application struct {
	Config       *config.Config
	Repository   *repository.Repository
	RequestLimit int
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
	// Инициализируйте и настройте объект Application
	app := &Application{
		Config:     cfg,
		Repository: repo,
		// Установите другие параметры вашего приложения, если необходимо
	}

	return app, nil
}

// Run запускает приложение.
func (app *Application) Run() {

	handler := api.NewHandler(app.Repository)
	r := gin.Default()

	// Группа запросов для периода
	PeriodGroup := r.Group("/period")
	{
		PeriodGroup.GET("/", handler.GetPeriods)
		PeriodGroup.GET("/:id_period", handler.GetPeriodByID)
		PeriodGroup.DELETE("/:id_period/delete", handler.DeletePeriod)
		PeriodGroup.POST("/create", handler.CreatePeriod)
		PeriodGroup.PUT("/:id_period/update", handler.UpdatePeriod)
		PeriodGroup.POST("/:id_period/fossil", handler.AddPeriodToFossil)
		PeriodGroup.DELETE("/:id_period/fossil/delete", handler.RemovePeriodFromFossil)
		PeriodGroup.POST("/:id_period/image", handler.AddPeriodImage)
	}

	// Группа запросов для ископаемого
	FossilGroup := r.Group("/fossil")
	{
		FossilGroup.GET("/", handler.GetFossil)
		FossilGroup.GET("/:id", handler.GetFossilByID)
		FossilGroup.DELETE("/:id/delete", handler.DeleteFossil)
		FossilGroup.PUT("/:id/update", handler.UpdateFossil)
		FossilGroup.PUT("/:id/status/user", handler.UpdateFossilStatusForUser)           // Новый маршрут для обновления статуса ископаемого пользователем
		FossilGroup.PUT("/:id/status/moderator", handler.UpdateFossilStatusForModerator) // Новый маршрут для обновления статуса ископаемого модератором
	}
	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	r.Run(addr)
	log.Println("Server down")
}
