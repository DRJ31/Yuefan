package model

type UserRegister struct {
	Username string
	Password string
	Code     string
	Key      string
}

type UserLogin struct {
	Username string
	Password string
}

type PasswordChange struct {
	Token string
	Old   string
	New   string
}

type RestaurantForm struct {
	Token      string
	Restaurant Restaurant
}

type Token struct {
	Token string
}
