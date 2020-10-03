package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-impatient/gaia/internal/service"
)

// UserHandler ...
type userHandler struct {
	UserService service.UserService
}

// Login 用户登录
func (handler *userHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.UserService.Login(c, "", "")

		c.JSON(http.StatusOK, gin.H{
			"text": "登录成功.",
		})
	}
}

// Register 注册
func (handler *userHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.UserService.Register(c, "", "", "")

		c.JSON(http.StatusOK, gin.H{
			"text": "注册成功.",
		})
	}
}

func MakeUserHandler(r *gin.Engine, srv *service.Services) {
	handler := &userHandler{UserService: srv.User}

	userGroup := r.Group("/user")
	{
		userGroup.GET("/login", handler.Login())
		userGroup.POST("/register", handler.Register())
	}
}
