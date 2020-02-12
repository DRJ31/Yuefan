package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DRJ31/yuefan/model"
	"github.com/DRJ31/yuefan/service"
	"github.com/go-redis/redis"
	"github.com/kataras/iris"
)

func userRouter(db *sql.DB, client *redis.Client, app *iris.Application) {
	app.Post("/api/user/register", func(ctx iris.Context) {
		var ru model.UserRegister
		ctx.ReadJSON(&ru)

		// Insert User
		err := service.InsertUser(db, ru)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": err,
			})
		} else {
			ctx.JSON(iris.Map{
				"message": "Successfully Registered",
			})
		}
	})

	app.Post("/api/user", func(ctx iris.Context) {
		var ul model.UserLogin
		ctx.ReadJSON(&ul)

		// User login
		token, has := service.UserLogin(db, client, ul)
		if !has {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{
				"message": "Wrong username or password",
			})
		} else {
			ctx.JSON(iris.Map{
				"message": "Login succeeded",
				"token":   token,
			})
		}
	})

	app.Get("/api/user", func(ctx iris.Context) {
		var user model.User
		token := ctx.URLParam("token")

		// Get current user from redis
		data, err := client.Get(fmt.Sprintf("token_%v", token)).Result()
		err = json.Unmarshal([]byte(data), &user)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": err,
			})
			return
		}

		ctx.JSON(iris.Map{
			"admin":    user.Role == 1,
			"username": user.Username,
		})
	})

	app.Post("/api/user/logout", func(ctx iris.Context) {
		var token model.Token
		ctx.ReadJSON(&token)

		// Get current user from redis
		user, err := getUser(client, token.Token)
		if err != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{
				"message": err,
			})
			return
		}

		// Logout
		service.UserLogout(client, user, token.Token)
		ctx.JSON(iris.Map{
			"message": "Successfully logout",
		})
	})

	app.Post("/api/user/password", func(ctx iris.Context) {
		var pc model.PasswordChange
		var newUser model.User
		ctx.ReadJSON(&pc)

		// Get current user from redis
		user, err := getUser(client, pc.Token)
		if err != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{
				"message": err,
			})
		}

		// Get new user object after change password
		newUser, err = service.ChangePassword(db, user, pc.Old, pc.New)
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": err,
			})
		} else {
			userStr, err := json.Marshal(newUser)
			if err != nil {
				ctx.StatusCode(500)
				ctx.JSON(iris.Map{
					"message": err,
				})
				return
			}

			client.Set(fmt.Sprintf("token_%v", pc.Token), userStr, 0)
			ctx.JSON(iris.Map{
				"message": "Successfully change password",
			})
		}
	})
}

func getUser(client *redis.Client, token string) (model.User, error) {
	var user model.User

	val, err := client.Get(fmt.Sprintf("token_%v", token)).Result()
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}
