package controller

import (
	"github.com/DRJ31/yuefan/model"
	"github.com/kataras/iris"
)

func Router(app *iris.Application) {
	// Initialize redis client
	client, err := model.InitRedis()
	if err != nil {
		panic(err)
	}
	// Initialize database client
	db, err := model.InitDB()
	if err != nil {
		panic(err)
	}

	// Router of user api
	userRouter(db, client, app)

	// Router of restaurant api
	restaurantRouter(db, client, app)
}
