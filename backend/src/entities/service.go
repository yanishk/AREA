package entities

type Service struct {
	Id    string
	Name  string
	Color string
	Logo  string
}

// type UserService struct {
// Id           string
// IdUser       string
// AccessToken  string
// RefreshToken string
// ExpiryToken  string
// IdService    string
// }

type ResultToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
