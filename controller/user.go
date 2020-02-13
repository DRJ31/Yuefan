package controller

import (
	"github.com/DRJ31/yuefan/service"
	"github.com/gin-gonic/gin"
)

func userRouter(router *gin.Engine) {
	router.POST("/api/user/register", service.InsertUser)

	router.POST("/api/user", service.UserLogin)

	router.GET("/api/user", service.GetCurrentUser)

	router.POST("/api/user/logout", service.UserLogout)

	router.POST("/api/user/password", service.ChangePassword)
}
