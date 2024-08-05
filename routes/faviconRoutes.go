package routes

import (
	"github.com/nedik/spp-lobby/controllers"

	"github.com/gin-gonic/gin"
)

func InitFaviconRoutes(router *gin.Engine) {
    router.GET("/favicon.ico", controllers.GetFavicon)
}

