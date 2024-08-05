package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nedik/spp-lobby/initializers"
	"github.com/nedik/spp-lobby/types"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: " + os.Args[0] + " <port>")
        os.Exit(1)
    }

    port64, err := strconv.ParseUint(os.Args[1], 10, 16)
    if err != nil {
        log.Fatal("Could not parse port from CLI params", err)
    }
    newServerPort := uint16(port64)

    config, err := initializers.LoadConfig(".")
    if err != nil {
        log.Fatal("Could not load environment variables", err)
    }

    serverInput := createRegisterServerInput(newServerPort, "Test Name")
    postBody, err := json.Marshal(serverInput)
    if err != nil {
        log.Fatal("Could not convert to json", err)
    }

    responseBody := bytes.NewBuffer(postBody)

    response, err := http.Post("http://localhost:" + config.ServerPort + "/servers", "application/json", responseBody)
    if err != nil {
        log.Fatalln(err)
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        log.Fatalln(err)
    }

    sb := string(body)
    log.Printf(sb)
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

