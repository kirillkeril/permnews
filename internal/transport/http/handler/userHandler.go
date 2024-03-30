package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"permnews/internal/models/entity"
	jwtService "permnews/internal/pkg/jwt-service"
	"permnews/internal/services"
	"permnews/internal/transport/http/models/input"
	"permnews/internal/transport/http/models/output"
)

type UserService interface {
	Create(user *entity.User, password string) error
	GetById(id string) (*entity.User, error)
}

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterUser
//
//	Post Регистрация пользователя
//	@Summary		Регистрация пользователя
//	@Description	Регистрация пользователя
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	input.UserInput	true	"User"
//	@Success		200		{object}
//	@Router			/user [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var inp = input.User{}
	err := c.BindJSON(inp)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(400)
		return
	}
	user := entity.NewUser(inp.Email)
	err = h.userService.Create(user, inp.Password)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(400)
		return
	}
	c.Status(http.StatusCreated)
	return
}

// Login
//
//	@Summary		Вход пользователя
//	@Description	Возвращает пару токенов
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	jwtService.JwtPair
//	@Failure		409	{object}	error
//	@Failure		401	{object}	error
//	@Router			/user/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	inp := &input.User{}
	if err := ctx.Bind(&inp); err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, httpErr := h.userService.Login(inp.Email, inp.Password)
	if httpErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": httpErr.Error()})
		return
	}
	tokenPair, err := jwtService.CreateTokenPair(*user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
		return
	}
	pair := output.TokenPair{Access: tokenPair.Access, Refresh: tokenPair.Refresh}
	ctx.JSON(http.StatusOK, pair)
	return
}
