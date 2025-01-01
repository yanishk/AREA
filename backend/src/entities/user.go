package entities

type User struct {
	Id             string
	Email          string
	Password       string
	CreatedAt      string
	Timezone       string
	ConnectionType string
}

type UserInfos struct {
	Email          string `json:"email"`
	CreatedAt      string `json:"createdat"`
	Timezone       string `json:"timezone"`
	ConnectionType string `json:"connectiontype"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModifyPassword struct {
	OldPassword string `json:"oldpassword"`
	Password    string `json:"password"`
}
