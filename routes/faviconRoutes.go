package routes

import (
	"github.com/nedik/spp-lobby/controllers"

	"github.com/gin-gonic/gin"
)

type FaviconRouteController struct {
    faviconController controllers.FaviconController
}

func NewFaviconRouteController(faviconController controllers.FaviconController) FaviconRouteController {
    return FaviconRouteController{faviconController}
}

func (self *FaviconRouteController) InitFaviconRoutes(routerGroup *gin.RouterGroup) {
    routerGroup.GET("/favicon.ico", self.faviconController.GetFavicon)
}

