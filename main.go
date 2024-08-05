package main

import (
	"log"

	"github.com/nedik/spp-lobby/initializers"
	"github.com/nedik/spp-lobby/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    config, err := initializers.LoadConfig(".")
    if err != nil {
        log.Fatal("Could not load environment variables", err)
    }

	router := gin.Default()

	routes.InitServerRoutes(router)
    routes.InitFaviconRoutes(router)

    log.Fatal(router.Run(":" + config.ServerPort))
}
