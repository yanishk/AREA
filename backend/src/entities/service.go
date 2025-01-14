package entities

type Service struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
	HasActions   bool   `json:"hasactions"`
	HasReactions bool   `json:"hasreactions"`
	IsAuthNeeded bool   `json:"isauthneeded"`
	Description  string `json:"description"`
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

type CallbackInformations struct {
	Service string `json:"service"`
	AppType string `json:"apptype"`
}

// Time & Date
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

// Reddit
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

type RedditVoteResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditUsername struct {
	Name string `json:"name"`
}

// Weather
type WeatherResponse struct {
	Current struct {
		Temperature float64 `json:"temp_c"`
	} `json:"current"`
	Forecast struct {
		ForecastDay []struct {
			Day struct {
				MaxTemperature float64 `json:"maxtemp_c"`
				MinTemperature float64 `json:"mintemp_c"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

// Github
type GithubWebhookTriggeredResponse struct {
	Repository struct {
		Name string `json:"full_name"`
	} `json:"repository"`
}

type GithubWebhookResponse struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Config struct {
		Url         string `json:"url"`
		ContentType string `json:"content_type"`
	} `json:"config"`
	Events []string `json:"events"`
}

type GithubRepository struct {
	Name string `json:"full_name"`
}

type GithubPullRequest struct {
	Number int `json:"number"`
}

type GithubBranch struct {
	Name string `json:"name"`
}

type GithubCommit struct {
	Message string `json:"message"`
}

type GithubIssue struct {
	Title string `json:"title"`
}

// Gitlab
type GitlabWebhookTriggeredResponse struct {
	Project struct {
		Id int `json:"id"`
	} `json:"project"`
}

type GitlabUser struct {
	Username string `json:"username"`
}

type GitlabProject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Linkedin
type LinkedinSub struct {
	Sub string `json:"sub"`
}

// Asana
type AsanaWorkspacesInfo struct {
	Data []struct {
		Gid  string `json:"gid"`
		Name string `json:"name"`
	} `json:"data"`
}

// Spotify
type SpotifyRequestParameters struct {
	BaseUrl       string
	ParameterName string
	EndUrl        string
	Method        string
}
