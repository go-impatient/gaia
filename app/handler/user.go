package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-impatient/gaia/internal/service"
	"net/http"
)

// UserHandler ...
type userHandler struct {
	userService service.UserService
}

// Login 用户登录
func (handler *userHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.userService.Login(c, "", "")

		c.JSON(http.StatusOK, gin.H{
			"text": "登录成功.",
		})
	}
}

// Register 注册
func (handler *userHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.userService.Register(c, "", "", "")

		c.JSON(http.StatusOK, gin.H{
			"text": "注册成功.",
		})
	}
}

func MakeUserHandler(r *gin.Engine, srv service.UserService) {
	handler := &userHandler{userService: srv}

	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", handler.Login())
		userGroup.POST("/register", handler.Register())
	}
}
