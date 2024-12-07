package entities

type Service struct {
	Id           string
	Name         string
	Color        string
	Logo         string
	HasActions   bool
	HasReactions bool
}

type ServiceUsefulInformations struct {
	Name  string
	Color string
	Logo  string
}

type ResultToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type TimeResponse struct {
	Timezone  string `json:"timezone"`
	Formatted string `json:"formatted"`
	Timestamp int    `json:"timestamp"`
	WeekDay   int    `json:"week_day"`
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
}
