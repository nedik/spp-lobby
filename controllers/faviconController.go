package controllers

import (
    "github.com/gin-gonic/gin"
)

func GetFavicon(c *gin.Context) {
    c.File("./assets/favicon.ico")
}

