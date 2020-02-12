package model

import (
	"database/sql"
	"github.com/go-redis/redis"
)

func InitRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, err
}

func InitDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "./yuefan.db")
}
