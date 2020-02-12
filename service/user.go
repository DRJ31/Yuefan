package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DRJ31/yuefan/model"
	"github.com/go-redis/redis"
	"time"
)

/*
Insert user while registration
@Param: db (*sql.DB)
@Param: ru (UserRegister) -> model/post.go
@Return: (error)
*/
func InsertUser(db *sql.DB, ru model.UserRegister) error {
	if len(ru.Username) == 0 || len(ru.Password) != 32 {
		return &MyError{
			time.Now(),
			"Bad Request",
		}
	}
	if checkKey(ru.Username, ru.Password, ru.Key) {
		return &MyError{
			time.Now(),
			"Wrong Key",
		}
	}
	row := db.QueryRow("SELECT * FROM user WHERE username=?", ru.Username)
	var user model.User
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err == nil {
		return &MyError{
			When: time.Now(),
			Err:  "User exist",
		}
	}

	var role int
	if ru.Code == CODE {
		role = 1
	} else if len(ru.Code) > 0 {
		return &MyError{
			When: time.Now(),
			Err:  "Wrong invitation code",
		}
	} else {
		role = 0
	}

	_, err = db.Exec(
		"INSERT INTO user (username,password,role) VALUES (?,?,?)",
		ru.Username,
		ru.Password,
		role,
	)
	if err != nil {
		return &MyError{
			time.Now(),
			"Insert User Failed",
		}
	}
	return nil
}

/*
Main logic of user login
@Param: db (*sql.DB)
@Param: client (*redis.Client)
@Param: ul (UserLogin) -> model/post.go
@Return: token (string)
@Return: status (bool)
*/
func UserLogin(db *sql.DB, client *redis.Client, ul model.UserLogin) (string, bool) {
	var check model.User

	// Check if the user has login
	val, err := client.Get(fmt.Sprintf("user_%v", ul.Username)).Result()
	if err == nil {
		data, _ := client.Get(fmt.Sprintf("token_%v", val)).Result()
		err = json.Unmarshal([]byte(data), &check)
		if check.Password != ul.Password {
			return "", false
		}
		client.Del(fmt.Sprintf("token_%v", val))
		token := getToken(ul.Username)
		client.Set(fmt.Sprintf("token_%v", token), data, 0)
		client.Set(fmt.Sprintf("user_%v", ul.Username), token, 0)
		return token, true
	} else if err == redis.Nil {
		row := db.QueryRow("SELECT * FROM user where username=? and password=?", ul.Username, ul.Password)
		err = row.Scan(&check.Id, &check.Username, &check.Password, &check.Role)
		var token string
		if err == nil {
			token = getToken(ul.Username)
			data, err := json.Marshal(check)
			if err != nil {
				fmt.Println("User serialize failed")
				return token, false
			}
			client.Set(fmt.Sprintf("token_%v", token), data, 0)
			client.Set(fmt.Sprintf("user_%v", ul.Username), token, 0)
		}
		return token, err == nil
	} else {
		return "", false
	}
}

/*
Main logic of user logout
@Param: client (*redis.Client)
@Param: token (string)
@Return: (error)
*/
func UserLogout(client *redis.Client, user model.User, token string) {
	client.Del(fmt.Sprintf("user_%v", user.Username))
	client.Del(fmt.Sprintf("token_%v", token))
}

/*
Main logic of change password
@Param: db (*sql.DB)
@Param: user (User)
@Param: old (string)
@Param: new (string)
@Return: (error)
*/
func ChangePassword(db *sql.DB, user model.User, old string, new string) (model.User, error) {
	var newUser model.User

	if user.Password != old {
		return user, &MyError{
			When: time.Now(),
			Err:  "Wrong password",
		}
	}

	_, err := db.Exec("UPDATE user SET password=? WHERE username=?", new, user.Username)
	row := db.QueryRow("SELECT * FROM user WHERE username=?", user.Username)
	row.Scan(&newUser.Id, &newUser.Username, &newUser.Password, &newUser.Role)
	return newUser, err
}
