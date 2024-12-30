package entities

type Service struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
	HasActions   bool   `json:"hasactions"`
	HasReactions bool   `json:"hasreactions"`
}

type ResultToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	Email string `json:"email"`
}

type TimeResponse struct {
	Timezone  string `json:"timezone"`
	Formatted string `json:"formatted"`
	Timestamp int    `json:"timestamp"`
	WeekDay   int    `json:"weekDay"`
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
}

type RedditPostResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				ID     string `json:"id"`
				Title  string `json:"title"`
				Author string `json:"author"`
				URL    string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditCommentResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				ID     string `json:"id"`
				Author string `json:"author"`
				Body   string `json:"body"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditUsername struct {
	Name string `json:"name"`
}
