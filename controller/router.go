package controller

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {

	// Router of user api
	userRouter(router)

	// Router of restaurant api
	restaurantRouter(router)
}
