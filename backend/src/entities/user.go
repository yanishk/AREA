package entities

type User struct {
	Id             string
	Email          string
	Password       string
	CreatedAt      string
	Timezone       string
	ConnectionType string
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ConnectionType string `json:"connectiontype"`
}

type UserModifyPassword struct {
	OldPassword string `json:"oldpassword"`
	Password    string `json:"password"`
}
