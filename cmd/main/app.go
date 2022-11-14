package main

import (
	"HTTP-REST-API/internal/domain/repository/postgres"
	"HTTP-REST-API/internal/handlers"
	"HTTP-REST-API/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// TODO разобраться как поднимать базу нормально и еще докер композе
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "12345"
		dbName   = "test_db"
	)

	db, err := sqlx.Connect("postgres", "postgres://postgres:12345@localhost:5432/testdb?sslmode=disable")
	defer db.Close()

	if err != nil {
		log.Fatalln(err)
	}

	repository := postgres.InitPostgresDb(db)
	initService, err := service.InitService(repository)

	if err != nil {
		log.Fatalln(err)
	}

	handler, err := handlers.NewHandler(initService)

	if err != nil {
		log.Fatalln(err)
	}

	router := handler.InitRoutes()
	router.Run("localhost:8080")
}
