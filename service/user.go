package service

import (
	"encoding/json"
	"fmt"
	"github.com/DRJ31/yuefan/model"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
)

func InsertUser(ctx *gin.Context) {
	// Receive data from request
	var ur model.UserRegister
	if err := ctx.ShouldBindJSON(&ur); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request Parameter"})
		return
	}

	// Initialize database
	db, err := model.InitDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database initialize error"})
		return
	}

	// Validate request parameter
	if len(ur.Username) == 0 || len(ur.Password) != 32 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request Parameter"})
		db.Close()
		return
	}
	if !checkKey(ur.Username, ur.Password, ur.Key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Key"})
		db.Close()
		return
	}

	// Get data from database
	row := db.QueryRow("SELECT * FROM user WHERE username=?", ur.Username)
	var user model.User
	err = row.Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err == nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User exist"})
		db.Close()
		return
	}

	// Validate invitation code
	var role int
	if ur.Code == CODE {
		role = 1
	} else if len(ur.Code) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Key"})
		db.Close()
		return
	} else {
		role = 0
	}

	// Insert into database
	_, err = db.Exec(
		"INSERT INTO user (username,password,role) VALUES (?,?,?)",
		ur.Username,
		ur.Password,
		role,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		db.Close()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully registered"})
	db.Close()
}

func UserLogin(ctx *gin.Context) {
	var check model.User
	var ul model.UserLogin

	// Read data from request
	if err := ctx.ShouldBindJSON(&ul); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request Parameter"})
		return
	}

	client, err := model.InitRedis()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis initialization failed"})
		return
	}

	// Check if the user has login
	val, err := client.Get(fmt.Sprintf("user_%v", ul.Username)).Result()
	if err == nil {
		data, _ := client.Get(fmt.Sprintf("token_%v", val)).Result()
		err = json.Unmarshal([]byte(data), &check)

		if check.Password != ul.Password && err == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password"})
			client.Close()
			return
		}
		client.Del(fmt.Sprintf("token_%v", val))
		token := getToken(ul.Username)
		client.Set(fmt.Sprintf("token_%v", token), data, 0)
		client.Set(fmt.Sprintf("user_%v", ul.Username), token, 0)
		ctx.JSON(http.StatusOK, gin.H{"token": token, "message": "Login Succeeded"})
		client.Close()
	} else if err == redis.Nil {
		// Initialize Databases
		db, err := model.InitDB()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database initialization failed"})
			client.Close()
			return
		}

		row := db.QueryRow("SELECT * FROM user where username=? and password=?", ul.Username, ul.Password)
		err = row.Scan(&check.Id, &check.Username, &check.Password, &check.Role)
		db.Close()
		var token string
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		token = getToken(ul.Username)
		data, err := json.Marshal(check)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User serialize failed"})
			client.Close()
			return
		}
		client.Set(fmt.Sprintf("token_%v", token), data, 0)
		client.Set(fmt.Sprintf("user_%v", ul.Username), token, 0)
		client.Close()
		ctx.JSON(http.StatusOK, gin.H{"token": token, "message": "Login Succeeded"})
	} else {
		client.Close()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error happened with redis"})
	}
}

func UserLogout(ctx *gin.Context) {
	var token model.Token
	if err := ctx.ShouldBindJSON(&token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request Parameter"})
		return
	}

	client, err := model.InitRedis()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis init failed"})
		return
	}

	user, err := getUser(client, token.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User did not login"})
		client.Close()
		return
	}

	client.Del(fmt.Sprintf("user_%v", user.Username))
	client.Del(fmt.Sprintf("token_%v", token.Token))
	client.Close()
}

func ChangePassword(ctx *gin.Context) {
	var newUser model.User
	var pc model.PasswordChange

	// Get data from request
	if err := ctx.ShouldBindJSON(&pc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request Parameter"})
		return
	}

	// Initialize Redis
	client, err := model.InitRedis()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis initialization failed"})
		return
	}

	// Get current user from redis
	user, err := getUser(client, pc.Token)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User did not login"})
		client.Close()
		return
	}

	// Compare old password
	if user.Password != pc.Old {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Old Password do not match"})
		client.Close()
		return
	}

	// Initialize database
	db, err := model.InitDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database initialization failed"})
		client.Close()
		return
	}

	// Database operation
	_, err = db.Exec("UPDATE user SET password=? WHERE username=?", pc.New, user.Username)
	row := db.QueryRow("SELECT * FROM user WHERE username=?", user.Username)
	err = row.Scan(&newUser.Id, &newUser.Username, &newUser.Password, &newUser.Role)
	db.Close()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		client.Close()
		return
	}

	data, err := json.Marshal(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		db.Close()
		return
	}

	client.Set(fmt.Sprintf("token_%v", pc.Token), data, 0)
	db.Close()
	ctx.JSON(http.StatusOK, gin.H{"error": "Successfully change password"})
}

func GetCurrentUser(ctx *gin.Context) {
	var user model.User
	token := ctx.Query("token")

	client, err := model.InitRedis()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis init failed"})
		return
	}

	// Get current user from redis
	data, err := client.Get(fmt.Sprintf("token_%v", token)).Result()
	err = json.Unmarshal([]byte(data), &user)
	client.Close()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"admin":    user.Role == 1,
		"username": user.Username,
	})
}

func getUser(client *redis.Client, token string) (model.User, error) {
	var user model.User

	data, err := client.Get(fmt.Sprintf("token_%v", token)).Result()
	if err != nil {
		return user, err
	}

	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
