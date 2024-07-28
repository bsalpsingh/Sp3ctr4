package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"github.com/sp3ctr4/database"
	"github.com/sp3ctr4/routes"
)

func main() {
	// load the env variables and instantiate the server app
	if err := env.Load("./.env"); err != nil {
		panic(err)
	}
	value := env.Get("PORT", "8080")

	database.Seed()

	ginEngine := gin.Default()
	routes.Init(ginEngine)

	ginEngine.Run(fmt.Sprintf("localhost:%v", value)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
