package main

import (
	"github.com/nedik/spp-lobby/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.InitServerRoutes(router)
    routes.InitFaviconRoutes(router)

	router.Run()
}
