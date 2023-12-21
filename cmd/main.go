package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"exam/api"
	"exam/config"
	"exam/storage/postgres"
)

func main() {

	var cfg = config.Load()
	fmt.Println(cfg)
	pgStorage, err := postgres.NewConnectionPostgres(&cfg)
	if err != nil {
		panic(err)
	}

	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	api.SetUpApi(r, &cfg, pgStorage)

	log.Println("Listening:", cfg.ServiceHost+cfg.ServiceHTTPPort, "...")
	if err := r.Run(cfg.ServiceHost + cfg.ServiceHTTPPort); err != nil {
		panic("Listent and service panic:" + err.Error())
	}
}
