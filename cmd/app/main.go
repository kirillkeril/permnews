package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"permnews/internal/app"
	"permnews/internal/pkg"
	"permnews/internal/storage/repository"
	"permnews/internal/transport/kafka"
	"syscall"
)

// @title			PermNews
// @description	Спецификация приложения
func main() {
	st := initDB()
	r := initGin()
	cons := initConsumer(st)
	appInstance := app.New(r, "0.0.0.0", "8080", cons)
	err := st.Migrate()
	if err != nil {
		log.Printf("migration err: %s", err.Error())
	}
	appInstance.InitServer()
	err = appInstance.Run(st)
	if err != nil {
		log.Fatal(err)
		return
	}

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		fmt.Printf("caught sig: %+v", sig)
		log.Println("Closing DB connection")
		st.Down()
		st.GetDb().Close()

		log.Println("Stopping http server")
	}()

}

func initConsumer(st repository.Storage) *kafka.Consumer {
	return kafka.NewConsumer([]string{"kafka:29092"}, "afisha", st)
	//return kafka.NewConsumer([]string{"localhost:9000"}, "afisha", st)
}

func initDB() *pkg.PostgresStorage {
	pg := pkg.NewPostgresStorage("postgres://postgres:PermNewsPSPU@postgres:5432/db?sslmode=disable", "file:///app/migrations")
	//pg := pkg.NewPostgresStorage("postgres://postgres:PermNewsPSPU@localhost:8082/db?sslmode=disable", "file:///app/migrations")
	return pg
}

func initGin() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type, access-control-allow-origin, access-control-allow-headers"},
	}))
	return r
}
