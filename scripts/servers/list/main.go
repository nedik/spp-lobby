package main

import (
    "io"
    "net/http"
    "log"

	"github.com/nedik/spp-lobby/initializers"
)

func main() {
    config, err := initializers.LoadConfig(".")
    if err != nil {
        log.Fatal("Could not load environment variables", err)
    }

    response, err := http.Get("http://localhost:" + config.ServerPort + "/servers")
    if err != nil {
        log.Fatalln(err)
    }

    body, err := io.ReadAll(response.Body)
    if err != nil {
        log.Fatalln(err)
    }

    sb := string(body)
    log.Printf(sb)
}
