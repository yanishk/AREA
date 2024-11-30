package entities

type Service struct {
	Id           string
	Name         string
	AccessToken  string
	RefreshToken string
	IdUser       string
	ExpiryToken  string
}

type ResultToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
