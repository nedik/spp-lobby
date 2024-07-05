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

func RegisterServer(c *gin.Context) {
    var registerServerInput types.RegisterServerInput

    if err := c.BindJSON(&registerServerInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    if !types.ValidateRegisterServerInput(registerServerInput) {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    incomingServerIP := c.ClientIP()
    incomingServerPort := registerServerInput.Port
    incomingServer := types.ConvertRegisterServerInputToServer(registerServerInput)
    incomingServer.IP = incomingServerIP

    // Find and update duplicate if exists
    for i, server := range servers {
        if server.IP == incomingServerIP && server.Port == incomingServerPort {
            servers[i] = incomingServer
            c.JSON(http.StatusCreated, gin.H{})
            return
        }
    }

    // If doesn't exist, then add a new one
    servers = append(servers, incomingServer)
    c.JSON(http.StatusCreated, gin.H{})
}

