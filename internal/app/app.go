package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"permnews/internal/storage/repository"
	"permnews/internal/transport/http/router"
)

type Consumer interface {
	StartListening()
}

type App struct {
	router   *gin.Engine
	srv      *http.Server
	addr     string
	port     string
	consumer Consumer
}

func New(router *gin.Engine, addr, port string, consumer Consumer) *App {
	return &App{
		router:   router,
		addr:     addr,
		port:     port,
		consumer: consumer,
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
	a.consumer.StartListening()
	if err := a.srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
