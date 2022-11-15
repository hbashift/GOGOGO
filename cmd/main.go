package main

import (
	"HTTP-REST-API/internal/domain/repository/postgres"
	"HTTP-REST-API/internal/handlers"
	"HTTP-REST-API/internal/service"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

func main() {
	err := initConfig()

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

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
	dns := fmt.Sprintf(":" + viper.GetString("port"))
	router.Run(dns)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
