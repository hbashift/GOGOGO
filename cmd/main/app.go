package main

import (
	"HTTP-REST-API/internal/domain/repository/postgres"
	"HTTP-REST-API/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "12345"
		dbName   = "test_db"
	)

	db, err := sqlx.Connect("postgres", "postgres://postgres:12345@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	router := gin.Default()
	test := postgres.InitPostgresDb(db)
	if err != nil {
		log.Fatalln(err)
	}

	handler, _ := handlers.NewHandler(test)

	router.GET("/balance/:id", handler.GetAccountBalance)
	router.POST("/account/add", handler.AddToAccountBalance)
	router.POST("/reservation", handler.ReserveUsersAmount)
	router.POST("/admit", handler.AdmitPurchase)
	// TODO здесь все методы для ТЗ
	router.Run("localhost:8080")

	/*	t := time.Now()
		p := t.Format("2006-01-02")

		fmt.Println(p)
	*/
}
