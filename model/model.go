package model

type User struct {
	Id       int
	Username string
	Password string
	Role     int
}

type Restaurant struct {
	Id     int
	Name   string
	UserId int
	Custom int
}
