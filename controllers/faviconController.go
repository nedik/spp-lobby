package controllers

import (
    "github.com/gin-gonic/gin"
)

type FaviconController struct {
}

func NewFaviconController() FaviconController {
    return FaviconController{}
}

func (self *FaviconController) GetFavicon(c *gin.Context) {
    c.File("./assets/favicon.ico")
}

