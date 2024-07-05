package routes

import (
	"github.com/nedik/spp-lobby/controllers"

	"github.com/gin-gonic/gin"
)

func InitServerRoutes(router *gin.Engine) {
    router.GET("/servers", controllers.ListAllServers)
    router.POST("/servers", controllers.RegisterServer)
}

