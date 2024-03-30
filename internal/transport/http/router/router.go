package router

import (
	nocache "github.com/alexander-melentyev/gin-nocache"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"permnews/internal/pkg/password_service"
	"permnews/internal/services"
	"permnews/internal/storage/repository"
	"permnews/internal/transport/http/handler"
)

func InitRouter(r *gin.Engine, storage repository.Storage) *gin.Engine {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(nocache.NoCache())

	api := r.Group("/api")
	userRepository := repository.NewUserRepository(storage)
	eventRepository := repository.NewEventsRepository(storage)
	eventService := services.NewEventsService(eventRepository)
	passwordService := password_service.NewPasswordService("test") // TODO: change secret
	userService := services.NewUserService(userRepository, passwordService)
	userGroup := api.Group("/user")
	userHandler := handler.NewUserHandler(*userService)

	userGroup.POST("/register", userHandler.RegisterUser)
	userGroup.POST("/login", userHandler.Login)

	eventGroup := api.Group("/events")
	eventHandler := handler.NewEventsHandler(eventService)
	eventGroup.GET("/", eventHandler.GetEvents)

	return r
}
