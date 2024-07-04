package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nedik/spp-lobby/routes"

	"github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
    router := gin.Default()
    routes.InitServerRoutes(router)
    return router
}

func TestGetServersEmptyList(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()

    req, _ := http.NewRequest("GET", "/servers", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "[]", w.Body.String())
}

