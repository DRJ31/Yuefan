package model

type UserRegister struct {
	Username string `json:""`
	Password string `json:""`
	Code     string `json:""`
	Key      string `json:""`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PasswordChange struct {
	Token string `json:"token"`
	Old   string `json:"old"`
	New   string `json:"new"`
}

type RestaurantForm struct {
	Token      string     `json:"token"`
	Restaurant Restaurant `json:"restaurant"`
}

type Token struct {
	Token string `json:"token"`
}
