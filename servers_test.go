package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nedik/spp-lobby/routes"
	"github.com/nedik/spp-lobby/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type InvalidInput struct {
	SomeValue string `json:"some_value"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	routes.InitServerRoutes(router)
	return router
}

func TestGetServersEmptyList(t *testing.T) {
	router := setupRouter()
	getServersAndAssert(t, router, []types.Server{})
}

func TestRegisterNewServer(t *testing.T) {
	router := setupRouter()

	registeredServer := registerServerAndAssert(t, router, 23073, "Test Server")

	getServersAndAssert(t, router, []types.Server{registeredServer})
}

func TestRegisterNewServerTwice(t *testing.T) {
	router := setupRouter()

	registerServerAndAssert(t, router, 23073, "Test Server")
	registeredServer := registerServerAndAssert(t, router, 23073, "Test Server")

	getServersAndAssert(t, router, []types.Server{registeredServer})
}

func TestRegisterTwoNewServers(t *testing.T) {
	router := setupRouter()

	registeredServer1 := registerServerAndAssert(t, router, 23073, "Test Server")
	registeredServer2 := registerServerAndAssert(t, router, 23074, "Test Server 2")

	getServersAndAssert(t, router, []types.Server{registeredServer1, registeredServer2})
}

func TestRegisterTwoNewServersTwice(t *testing.T) {
	router := setupRouter()

	registeredServer1 := registerServerAndAssert(t, router, 23073, "Test Server")
	registeredServer2 := registerServerAndAssert(t, router, 23074, "Test Server 2")
	registerServerAndAssert(t, router, 23073, "Test Server")
	registerServerAndAssert(t, router, 23074, "Test Server 2")

	getServersAndAssert(t, router, []types.Server{registeredServer1, registeredServer2})
}

func TestMissingFieldsOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()

	newServerInputWithMissingFields := types.RegisterServerInput{
		Advanced:       true,
		AntiCheatOn:    true,
		BonusFrequency: 10,
		CurrentMap:     "ctf_Ash",
		Info:           "Test Server Info",
		MaxPlayers:     32,
		NumBots:        1,
		OS:             "Linux",
		Port:           23073,
		Realistic:      true,
		Respawn:        1,
		Survival:       true,
		Version:        "1.0",
		WM:             true,
	}

	invalidInputJson, _ := json.Marshal(newServerInputWithMissingFields)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestInvalidInputOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()

	invalidInput := InvalidInput{
		SomeValue: "test",
	}

	invalidInputJson, _ := json.Marshal(invalidInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongCountryOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.Country = createLongString(types.MaxCountryLength + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongMapNameOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.CurrentMap = createLongString(types.MaxMapSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongGameStyleOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.GameStyle = createLongString(types.MaxGameStyleSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongInfoOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.Info = createLongString(types.MaxInfoSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongNameOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, createLongString(types.MaxNameSize+1))
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongOSOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.OS = createLongString(types.MaxOSSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongPlayerNameOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.Players[1] = createLongString(types.MaxPlayerNameSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func TestTooLongVersionOnRegisteringNewServer(t *testing.T) {
	router := setupRouter()
	registerServerInput := createRegisterServerInput(23073, "Test Name")
	registerServerInput.Version = createLongString(types.MaxVersionSize + 1)
	invalidInputJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", invalidInputJson)
	assert.Equal(t, http.StatusBadRequest, returnedCode)
}

func createLongString(size int) string {
	var stringBuilder strings.Builder
	stringBuilder.Grow(size)
	for i := 0; i < size; i++ {
		stringBuilder.WriteByte(0)
	}
	result := stringBuilder.String()
	return result
}

func registerServerAndAssert(t *testing.T, router *gin.Engine, port uint16, name string) types.Server {
	registerServerInput := createRegisterServerInput(port, name)
	serverJson, _ := json.Marshal(registerServerInput)
	returnedCode, _ := sendJsonToPostEndpoint(router, "/servers", serverJson)
	assert.Equal(t, http.StatusCreated, returnedCode)

	registeredServer := types.ConvertRegisterServerInputToServer(registerServerInput)
	return registeredServer
}

func createRegisterServerInput(port uint16, name string) types.RegisterServerInput {
	registerServerInput := types.RegisterServerInput{
		Advanced:       true,
		AntiCheatOn:    true,
		BonusFrequency: 10,
		Country:        "PL",
		CurrentMap:     "ctf_Ash",
		GameStyle:      "CTF",
		Info:           "Test Server Info",
		MaxPlayers:     32,
		Name:           name,
		NumBots:        1,
		OS:             "Linux",
		Players:        []string{"test_player_1", "another player"},
		Port:           port,
		Private:        true,
		Realistic:      true,
		Respawn:        1,
		Survival:       true,
		Version:        "1.0",
		WM:             true,
	}

	return registerServerInput
}

func getServersAndAssert(t *testing.T, router *gin.Engine, expected_servers []types.Server) {
	code, body := getServers(router)
	assert.Equal(t, http.StatusOK, code)
	expected_servers_json, _ := json.Marshal(expected_servers)
	assert.Equal(t, string(expected_servers_json), body)
}

func getServers(router *gin.Engine) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/servers", nil)
	router.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func sendJsonToPostEndpoint(router *gin.Engine, endpoint string, dataJson []byte) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(string(dataJson)))
	router.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}
