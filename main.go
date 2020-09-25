package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-impatient/gaia/pkg/http/ginhttp"
)

func main() {
	server := ginhttp.NewServer(ginhttp.Addr(":5000"))
	server.Router().GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"text": "Welcome to api app.",
		})
	})
	server.Run()
}
