package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FrostyCreator/NewsCollector/config"
	"github.com/FrostyCreator/NewsCollector/controller"
	"github.com/FrostyCreator/NewsCollector/server"
	"github.com/FrostyCreator/NewsCollector/store/db"
)

func main(){
	if err := run(); err != nil{
		log.Fatal(err)
	}
}


func run() error {
	ctx := context.Background()

	// config
	cfg := config.GetConfig()

	// connect to database
	pgDB, err := db.Dial(*cfg)
	if err != nil {
		return err
	}
	defer pgDB.Close()



	newsRepo := db.NewNewsRepo(pgDB)
	newsController := controller.NewNewsController(ctx, newsRepo)
	router := server.NewRouter(newsController)

	// Обновление новостей в бд
	newsController.AddNews()

	// create new server instance and run http server
	addr := ":8080"
	_, err = server.Init(ctx, cfg, newsRepo, *router,  addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Running http server on %s\n", addr)

	return nil
}