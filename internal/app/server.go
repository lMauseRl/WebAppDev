package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/WebAppDev/docs"
	"github.com/lud0m4n/WebAppDev/internal/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// Run запускает приложение.
func (app *Application) Run() {
	r := gin.Default()
	// Это нужно для автоматического создания папки "docs" в вашем проекте
	docs.SwaggerInfo.Title = "BagTracker RestAPI"
	docs.SwaggerInfo.Description = "API server for BagTracker application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Группа запросов для периода
	PeriodGroup := r.Group("/period")
	{
		PeriodGroup.GET("/", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetPeriods)
		PeriodGroup.GET("/:id_period", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetPeriodByID)
		PeriodGroup.DELETE("/:id_period/delete", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.DeletePeriod)
		PeriodGroup.POST("/create", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.CreatePeriod)
		PeriodGroup.PUT("/:id_period/update", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.UpdatePeriod)
		PeriodGroup.POST("/:id_period/fossil", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddPeriodToFossil)
		PeriodGroup.DELETE("/:id_period/fossil/delete", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.RemovePeriodFromFossil)
		PeriodGroup.POST("/:id_period/image", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddPeriodImage)
	}
	// Группа запросов для ископаемого
	FossilGroup := r.Group("/fossil").Use(middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository))
	{
		FossilGroup.GET("/", app.Handler.GetFossil)
		FossilGroup.GET("/:id", app.Handler.GetFossilByID)
		FossilGroup.DELETE("/:id/delete", app.Handler.DeleteFossil)
		FossilGroup.PUT("/:id/update", app.Handler.UpdateFossil)
		FossilGroup.PUT("/:id/status/user", app.Handler.UpdateFossilStatusForUser)           // Новый маршрут для обновления статуса ископаемого пользователем
		FossilGroup.PUT("/:id/status/moderator", app.Handler.UpdateFossilStatusForModerator) // Новый маршрут для обновления статуса ископаемого модератором
	}

	UserGroup := r.Group("/user")
	{
		UserGroup.GET("/", app.Handler.GetUserByID)
		UserGroup.POST("/register", app.Handler.Register)
		UserGroup.POST("/login", app.Handler.Login)
		UserGroup.POST("/logout", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.Logout)
	}
	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	r.Run(addr)
	log.Println("Server down")
}
