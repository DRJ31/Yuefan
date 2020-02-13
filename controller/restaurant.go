package controller

import (
	"github.com/DRJ31/yuefan/service"
	"github.com/gin-gonic/gin"
)

func restaurantRouter(router *gin.Engine) {
	// Insert a restaurant
	router.POST("/api/restaurant", service.InsertRestaurant)

	// Get all restaurants
	router.GET("/api/restaurants", service.GetAllRestaurants)

	// Get custom restaurants
	router.POST("/api/restaurants", service.GetCustomRestaurants)

	// Delete restaurant
	router.DELETE("/api/restaurant", service.DeleteRestaurant)
}
