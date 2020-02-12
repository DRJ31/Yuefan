package controller

import (
	"database/sql"
	"github.com/DRJ31/yuefan/model"
	"github.com/DRJ31/yuefan/service"
	"github.com/go-redis/redis"
	"github.com/kataras/iris"
)

func restaurantRouter(db *sql.DB, client *redis.Client, app *iris.Application) {
	// Insert a restaurant
	app.Post("/api/restaurant", func(ctx iris.Context) {
		var restaurantForm model.RestaurantForm
		ctx.ReadJSON(&restaurantForm)

		// Get current user information
		user, err := getUser(client, restaurantForm.Token)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": "User not exist",
			})
		}

		// Insert restaurant
		err = service.InsertRestaurant(db, user, restaurantForm.Restaurant)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": err,
			})
		} else {
			ctx.JSON(iris.Map{
				"message": "Successfully inserted restaurant",
			})
		}
	})

	// Get all restaurants
	app.Get("/api/restaurants", func(ctx iris.Context) {
		restaurants, err := service.GetAllRestaurants(db)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": "Failed to get all restaurants",
			})
		} else {
			ctx.JSON(iris.Map{
				"restaurants": restaurants,
			})
		}
	})

	// Get custom restaurants
	app.Post("/api/restaurants", func(ctx iris.Context) {
		var token model.Token
		ctx.ReadJSON(&token)

		// Get current user information
		user, err := getUser(client, token.Token)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": "User not exist",
			})
			return
		}

		// Get custom restaurant
		restaurants, err := service.GetCustomRestaurants(db, user)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": err,
			})
		} else {
			ctx.JSON(iris.Map{
				"restaurants": restaurants,
			})
		}
	})

	// Delete restaurant
	app.Delete("/api/restaurant", func(ctx iris.Context) {
		var restaurant model.Restaurant
		token := ctx.URLParam("token")

		// Check current user
		user, err := getUser(client, token)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": "User not exist",
			})
			return
		}

		// Validate id
		id, err := ctx.URLParamInt("id")
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": "Format of id is invalid",
			})
			return
		}

		// Get restaurant from database
		row := db.QueryRow("SELECT * FROM restaurant WHERE ID=?", id)
		err = row.Scan(&restaurant.Id, &restaurant.Name, &restaurant.UserId, &restaurant.Custom)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": "Restaurant not exist",
			})
			return
		}

		// Delete restaurant
		err = service.DeleteRestaurant(db, user, restaurant)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": err,
			})
		} else {
			ctx.JSON(iris.Map{
				"message": "Successfully deleted",
			})
		}
	})
}
