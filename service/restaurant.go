package service

import (
	"github.com/DRJ31/yuefan/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InsertRestaurant(ctx *gin.Context) {
	var restaurantForm model.RestaurantForm
	var check model.Restaurant

	if err := ctx.ShouldBindJSON(&restaurantForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	restaurant := restaurantForm.Restaurant

	// Initialize Redis
	client, err := model.InitRedis()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis initialization failed"})
		return
	}

	// Get current user information
	user, err := getUser(client, restaurantForm.Token)
	client.Close()
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User login is required",
		})
		return
	}

	// Validate information of restaurant
	if len(restaurant.Name) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant name"})
		return
	}
	restaurant.UserId = user.Id

	db, err := model.InitDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database initialization failed"})
		return
	}
	row := db.QueryRow("SELECT ID from restaurant WHERE name=?", restaurant.Name)
	err = row.Scan(&check.Id)
	if err == nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Duplicated name"})
		return
	}

	// Insert restaurant information into database
	_, err = db.Exec(
		"INSERT INTO restaurant (name, user_id, custom) VALUES (?,?,?)",
		restaurant.Name,
		restaurant.UserId,
		restaurant.Custom,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Insert succeeded"})
}

func GetAllRestaurants(ctx *gin.Context) {
	restaurants := make([]model.Restaurant, 0)

	// Initialize database
	db, err := model.InitDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database initialization failed"})
		return
	}

	rows, err := db.Query("SELECT  * from restaurant where custom=0")
	db.Close()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var restaurant model.Restaurant
		if err = rows.Scan(&restaurant.Id, &restaurant.Name, &restaurant.UserId, &restaurant.Custom); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		restaurants = append(restaurants, restaurant)
	}
	ctx.JSON(http.StatusOK, gin.H{"restaurants": restaurants})
}

func GetCustomRestaurants(ctx *gin.Context) {
	var token model.Token
	if err := ctx.ShouldBindJSON(&token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize Redis
	client, err := model.InitRedis()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis initialization failed"})
		return
	}

	// Get current user information
	user, err := getUser(client, token.Token)
	client.Close()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not exist"})
		return
	}

	// Initialize database
	db, err := model.InitDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database initialization failed"})
		return
	}

	restaurants := make([]model.Restaurant, 0)
	rows, err := db.Query("SELECT  * from restaurant where custom=1 and user_id=?", user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for rows.Next() {
		var restaurant model.Restaurant
		if err = rows.Scan(&restaurant.Id, &restaurant.Name, &restaurant.UserId, &restaurant.Custom); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		restaurants = append(restaurants, restaurant)
	}
	ctx.JSON(http.StatusOK, gin.H{"restaurants": restaurants})
}

func DeleteRestaurant(ctx *gin.Context) {
	var restaurant model.Restaurant
	token := ctx.Query("token")

	// Initialize Redis
	client, err := model.InitRedis()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis initialization failed"})
		return
	}

	// Check current user
	user, err := getUser(client, token)
	client.Close()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not exist"})
		return
	}

	id := ctx.Query("id")

	// Initialize database
	db, err := model.InitDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database initialization failed"})
		return
	}

	// Get restaurant from database
	row := db.QueryRow("SELECT * FROM restaurant WHERE ID=?", id)
	err = row.Scan(&restaurant.Id, &restaurant.Name, &restaurant.UserId, &restaurant.Custom)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant does not exist"})
		db.Close()
		return
	}

	if user.Role != 1 && restaurant.Custom != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to do that"})
		db.Close()
		return
	}

	_, err = db.Exec("DELETE from restaurant where ID=?", restaurant.Id)
	db.Close()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
