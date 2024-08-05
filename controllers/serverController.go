package controllers

import (
    "errors"
    "strconv"
    "net/http"

    "github.com/nedik/spp-lobby/types"

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

func GetSpecificServer(c *gin.Context) {
    ip := c.Param("ip")
    port, err := getPortFromParams(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    foundServer, err := findServer(ip, port)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusNotFound, foundServer)
}

func GetPlayersOfServer(c *gin.Context) {
    ip := c.Param("ip")
    port, err := getPortFromParams(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    foundServer, err := findServer(ip, port)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, foundServer.Players)
}

func getPortFromParams(c *gin.Context) (uint16, error) {
    port64, err := strconv.ParseUint(c.Param("port"), 10, 16)
    if err != nil {
        return 0, errors.New("Invalid port")
    }
    port := uint16(port64)

    return port, nil
}

func findServer(ip string, port uint16) (*types.Server, error) {
    for _, candidateServer := range servers {
        if candidateServer.IP == ip && candidateServer.Port == port {
            return &candidateServer, nil
        }
    }

    return nil, errors.New("server not found")
}

