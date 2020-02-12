package service

import (
	"database/sql"
	"fmt"
	"github.com/DRJ31/yuefan/model"
	"time"
)

/*
Insert information of restaurants
@Param: db (*sql.DB)
@Param: user (User)
@Param: restaurant (Restaurant)
@Return: (error)
*/
func InsertRestaurant(db *sql.DB, user model.User, restaurant model.Restaurant) error {
	var check model.Restaurant

	// Validate information of restaurant
	if len(restaurant.Name) == 0 {
		fmt.Println(restaurant)
		return &MyError{
			When: time.Now(),
			Err:  "Invalid Request",
		}
	}
	restaurant.UserId = user.Id
	row := db.QueryRow("SELECT ID from restaurant WHERE name=?", restaurant.Name)
	err := row.Scan(&check.Id)
	if err == nil {
		return &MyError{
			When: time.Now(),
			Err:  "Duplicated",
		}
	}

	// Insert restaurant information into database
	_, err = db.Exec(
		"INSERT INTO restaurant (name, user_id, custom) VALUES (?,?,?)",
		restaurant.Name,
		restaurant.UserId,
		restaurant.Custom,
	)
	if err != nil {
		return err
	}
	return nil
}

/*
Get information of all restaurants
@Param: db (*sql.DB)
@Return: restaurants ([]Restaurant)
@Return: (error)
*/
func GetAllRestaurants(db *sql.DB) ([]model.Restaurant, error) {
	restaurants := make([]model.Restaurant, 0)
	rows, err := db.Query("SELECT  * from restaurant where custom=0")
	if err != nil {
		return restaurants, err
	}
	for rows.Next() {
		var restaurant model.Restaurant
		if err = rows.Scan(&restaurant.Id, &restaurant.Name, &restaurant.UserId, &restaurant.Custom); err != nil {
			return restaurants, err
		}
		restaurants = append(restaurants, restaurant)
	}
	return restaurants, err
}

/*
Get information of all custom restaurants
@Param: db (*sql.DB)
@Param: user (User)
@Return: restaurants ([]Restaurant)
@Return: (error)
*/
func GetCustomRestaurants(db *sql.DB, user model.User) ([]model.Restaurant, error) {
	restaurants := make([]model.Restaurant, 0)
	rows, err := db.Query("SELECT  * from restaurant where custom=1 and user_id=?", user.Id)
	if err != nil {
		return restaurants, err
	}
	for rows.Next() {
		var restaurant model.Restaurant
		if err = rows.Scan(&restaurant.Id, &restaurant.Name, &restaurant.UserId, &restaurant.Custom); err != nil {
			return restaurants, err
		}
		restaurants = append(restaurants, restaurant)
	}
	return restaurants, err
}

/*
Delete a restaurant
@Param: db (*sql.DB)
@Param: user (User)
@Param: restaurant (Restaurant)
@Return: (error)
*/
func DeleteRestaurant(db *sql.DB, user model.User, restaurant model.Restaurant) error {
	if user.Role != 0 && restaurant.Custom != 1 {
		return &MyError{
			When: time.Now(),
			Err:  "You don't have permission to delete the restaurant.",
		}
	}
	_, err := db.Exec("DELETE from restaurant where ID=?", restaurant.Id)
	return err
}
