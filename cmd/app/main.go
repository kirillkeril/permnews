package main

import (
	"authorization_service/internal/pkg"
	"authorization_service/internal/transport/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

// @title			PermNews
// @description	Спецификация приложения
func main() {
	st := initDB()
	r := initGin()
	app := http.New(r, "0.0.0.0", "8080")
	app.InitServer()
	err := app.Run(st)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func initDB() *pkg.PostgresStorage {
	pg := pkg.NewPostgresStorage("postgres://postgres:PermNewsPSPU@localhost:8082/db?sslmode=disable", "")
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
