package http

import (
	"authorization_service/internal/storage/repository"
	"authorization_service/internal/transport/http/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type App struct {
	router *gin.Engine
	srv    *http.Server
	addr   string
	port   string
}

func New(router *gin.Engine, addr, port string) *App {
	return &App{
		router: router,
		addr:   addr,
		port:   port,
	}
}

func (a *App) InitServer() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", a.addr, a.port),
		Handler: a.router,
	}
	a.srv = srv
}

func (a *App) Run(storage repository.Storage) error {
	a.router = router.InitRouter(a.router, storage)
	a.InitServer()
	if err := a.srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
