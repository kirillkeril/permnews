package router

import (
	"authorization_service/internal/pkg/password_service"
	"authorization_service/internal/services"
	"authorization_service/internal/storage/repository"
	"authorization_service/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(r *gin.Engine, storage repository.Storage) *gin.Engine {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	userRepository := repository.NewUserRepository(storage)
	passwordService := password_service.NewPasswordService("test") // TODO: change secret
	userService := services.NewUserService(userRepository, passwordService)
	userGroup := api.Group("/user")
	userHandler := handler.NewUserHandler(*userService)

	userGroup.POST("/register", userHandler.RegisterUser)
	userGroup.POST("/login", userHandler.Login)

	return r
}
