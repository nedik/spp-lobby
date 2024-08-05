package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/nedik/spp-lobby/types"

	"github.com/gin-gonic/gin"
	"github.com/igrmk/treemap/v2"
)

type TServersByUpdatedTime = *treemap.TreeMap[int64, []types.Server]

var serversByUpdatedTime = treemap.New[int64, []types.Server]()
var serversByIPPort = make(map[string]types.Server)

func treeToList(serversTree TServersByUpdatedTime) []types.Server {
    serversList := []types.Server{}
    for it := serversTree.Iterator(); it.Valid(); it.Next() {
        for _, currentServer := range it.Value() {
            serversList = append(serversList, currentServer)
        }
    }
    return serversList
}

func ListAllServers(c *gin.Context) {
    c.JSON(http.StatusOK, treeToList(serversByUpdatedTime))
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
    incomingServerIPPort := convertToIPPort(incomingServerIP, incomingServerPort)
    incomingServer.UpdatedAt = time.Now().Unix()

    // Find and update duplicate if exists
    foundServer, serverFound := serversByIPPort[incomingServerIPPort]
    if serverFound {
        updateInServersTree(serversByUpdatedTime, foundServer, incomingServer.UpdatedAt)
        serverByIPPort := serversByIPPort[incomingServerIPPort]
        serverByIPPort.UpdatedAt = incomingServer.UpdatedAt
        serversByIPPort[incomingServerIPPort] = serverByIPPort

        c.JSON(http.StatusCreated, gin.H{})
        return
    }

    // If doesn't exist, then add a new one
    appendToServersTree(serversByUpdatedTime, incomingServer)
    serversByIPPort[incomingServerIPPort] = incomingServer
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
    serverIPPort := convertToIPPort(ip, port)
    candidateServer, serverFound := serversByIPPort[serverIPPort]
    if serverFound {
        return &candidateServer, nil
    }

    return nil, errors.New("server not found")
}

func convertToIPPort(ip string, port uint16) string {
    return ip + ":" + strconv.FormatUint(uint64(port), 10)
}

func appendToServersTree(serversTree TServersByUpdatedTime, newServer types.Server) {
    serverListNow, found := serversTree.Get(newServer.UpdatedAt)
    if found {
        serverListNowUpdated := append(serverListNow, newServer)
        serversTree.Del(newServer.UpdatedAt)
        serversTree.Set(newServer.UpdatedAt, serverListNowUpdated)
    } else {
        serversTree.Set(newServer.UpdatedAt, []types.Server{newServer})
    }
}

func updateInServersTree(serversTree TServersByUpdatedTime, server types.Server, newUpdateTime int64) {
    lastUpdatedAt := server.UpdatedAt
    serversListAtTime, found := serversByUpdatedTime.Get(lastUpdatedAt)
    if found {
        // Find the server in the list
        var serverIndex *int = nil
        for index, currentServer := range serversListAtTime {
            if currentServer.IP == server.IP && currentServer.Port == server.Port {
                serverIndex = new(int)
                *serverIndex = index
                break
            }
        }

        if serverIndex != nil {
            // Remove server from the list
            lastIndex := len(serversListAtTime) - 1
            serversListAtTime[*serverIndex] = serversListAtTime[lastIndex]
            serversListAtTime = serversListAtTime[:lastIndex]

            // Remove current list and put back  on the tree without the server
            serversByUpdatedTime.Del(lastUpdatedAt)
            serversByUpdatedTime.Set(lastUpdatedAt, serversListAtTime)
        }
    }

    server.UpdatedAt = newUpdateTime
    appendToServersTree(serversTree, server)
}

