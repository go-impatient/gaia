package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-impatient/gaia/app/handler"
	sdHandle "github.com/go-impatient/gaia/app/handler/sd"
	"github.com/go-impatient/gaia/app/middleware"
	"github.com/go-impatient/gaia/internal/model"
	"github.com/go-impatient/gaia/internal/service"
	"net/http"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to api app.",
	})
}

// NotFound creates a gin middleware for handling page not found.
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, &model.ErrorResponseType{
			Error:     http.StatusText(http.StatusNotFound),
			ErrorCode: http.StatusNotFound,
			Message:   "page not found",
		})
	}
}

// NewRouter ...
func NewRouter(router *gin.Engine, services *service.Services) {
	// 使用中间件.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Handler())
	router.Use(middleware.NoCache)
	router.Use(middleware.Options)
	router.Use(middleware.Secure)
	router.Use(middleware.RequestId())

	// 404 Handler.
	router.NoRoute(NotFound())

	router.GET("/", rootHandler)

	// The health check handlers
	svcd := router.Group("/sd")
	{
		svcd.GET("/health", sdHandle.HealthCheck)
		svcd.GET("/disk", sdHandle.DiskCheck)
		svcd.GET("/cpu", sdHandle.CPUCheck)
		svcd.GET("/ram", sdHandle.RAMCheck)
	}

	handler.MakeUserHandler(router, services.User)
}
