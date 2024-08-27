package routes

import (
	"github.com/nedik/spp-lobby/controllers"

	"github.com/gin-gonic/gin"
)

type ServerRouteController struct {
    serverController controllers.ServerController
}

func NewServerRouteController(serverController controllers.ServerController) ServerRouteController {
    return ServerRouteController{serverController}
}

func (self *ServerRouteController) InitServerRoutes(routerGroup *gin.RouterGroup) {
    router := routerGroup.Group("servers")

    router.GET("", self.serverController.ListAllServers)
    router.GET("/:ip/:port", self.serverController.GetSpecificServer)
    router.GET("/:ip/:port/players", self.serverController.GetPlayersOfServer)
    router.POST("", self.serverController.RegisterServer)
}

