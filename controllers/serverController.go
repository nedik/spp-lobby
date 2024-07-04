package controllers

import (
    "github.com/nedik/spp-lobby/types"

    "net/http"

    "github.com/gin-gonic/gin"
)

var servers = []types.Server{}

func ListAllServers(c *gin.Context) {
    c.JSON(http.StatusOK, servers)
}

