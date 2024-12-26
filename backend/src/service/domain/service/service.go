package service_service

import (
	"backend/src/entities"
	"backend/src/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ServiceService struct {
	ServiceRepository  storage.ServiceRepository
	UserRepository     storage.UserRepository
	ActionRepository   storage.ActionRepository
	WorkflowRepository storage.WorkflowRepository
	ReactionRepository storage.ReactionRepository
}

const ContentType = "Content-Type"
const ContentTypeUrlEncoded = "application/x-www-form-urlencoded"

const InvalidCallbackTypeMessage = "Invalid callback type"
const APICallFailedMessage = "API call failed"
const grantTypeRefreshToken = "refresh_token"
const grantTypeAuthorization = "authorization_code"
const oneYearSecond = 31536000

func NewServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository,
	ActionRepository storage.ActionRepository, WorkflowRepository storage.WorkflowRepository,
	ReactionRepository storage.ReactionRepository) *ServiceService {
	return &ServiceService{
		ServiceRepository:  ServiceRepository,
		UserRepository:     UserRepository,
		ActionRepository:   ActionRepository,
		WorkflowRepository: WorkflowRepository,
		ReactionRepository: ReactionRepository,
	}
}

func (self *ServiceService) FindServiceByName(name string) (entities.Service, error) {
	return self.ServiceRepository.FindServiceByName(name)
}

func (self *ServiceService) FindServiceById(id string) (entities.Service, error) {
	return self.ServiceRepository.FindServiceById(id)
}

func (self *ServiceService) RetrieveActionsFromService(serviceName string) ([]entities.Action, error) {
	foundService, errFoundService := self.ServiceRepository.FindServiceByName(serviceName)
	if errFoundService != nil {
		return nil, errFoundService
	}

	actions, errFindActions := self.ActionRepository.FindActionByServiceId(foundService.Id)
	if errFindActions != nil {
		return nil, errFindActions
	}
	return actions, nil
}

func (self *ServiceService) RetrieveReactionsFromService(serviceName string) ([]entities.Reaction, error) {
	foundService, errFoundService := self.ServiceRepository.FindServiceByName(serviceName)
	if errFoundService != nil {
		return nil, errFoundService
	}

	reactions, errFindReactions := self.ReactionRepository.FindReactionByServiceId(foundService.Id)
	if errFindReactions != nil {
		return nil, errFindReactions
	}
	return reactions, errFindReactions
}

func (self *ServiceService) RetrieveActionsServices() ([]entities.Service, error) {
	return self.ServiceRepository.FindActionsServices()
}

func (self *ServiceService) RetrieveReactionsServices() ([]entities.Service, error) {
	return self.ServiceRepository.FindReactionsServices()
}

func GetGoogleOAuth2AccessTokenWebRequestBody(code string, callback string) string {
	return fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s",
		code,
		os.Getenv("GOOGLE_WEB_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		callback,
		grantTypeAuthorization,
	)
}

func GetGoogleOAuth2AccessTokenMobileRequestBody(code string, callback string) string {
	return fmt.Sprintf(
		"code=%s&client_id=%s&redirect_uri=%s&grant_type=%s",
		code,
		os.Getenv("GOOGLE_MOBILE_CLIENT_ID"),
		callback,
		grantTypeAuthorization,
	)
}

func GetGoogleOAuth2AccessTokenRequest(code string, callbackType string, appType string) (*http.Request, error) {
	var googleCallback string
	var jsonBody string

	if callbackType == "service" {
		googleCallback = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		googleCallback = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	if appType == "web" {
		jsonBody = GetGoogleOAuth2AccessTokenWebRequestBody(code, googleCallback)
	} else if appType == "mobile" {
		jsonBody = GetGoogleOAuth2AccessTokenMobileRequestBody(code, googleCallback)
	} else {
		return nil, fmt.Errorf("Invalid app type")
	}

	tokenUrl := "https://oauth2.googleapis.com/token"
	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(ContentType, ContentTypeUrlEncoded)
	return request, nil
}

func GetGithubOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var githubCallback string
	var githubClientID string
	var githubClientSecret string

	if callbackType == "service" {
		githubCallback = os.Getenv("GITHUB_SERVICE_CALLBACK")
		githubClientID = os.Getenv("GITHUB_SERVICE_CLIENT_ID")
		githubClientSecret = os.Getenv("GITHUB_SERVICE_CLIENT_SECRET")
	} else if callbackType == "login" {
		githubCallback = os.Getenv("GITHUB_LOGIN_CALLBACK")
		githubClientID = os.Getenv("GITHUB_LOGIN_CLIENT_ID")
		githubClientSecret = os.Getenv("GITHUB_LOGIN_CLIENT_SECRET")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	tokenUrl := "https://github.com/login/oauth/access_token"
	jsonBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		githubClientID,
		githubClientSecret,
		code,
		githubCallback,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func GetSpotifyOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var spotifyCallback string

	if callbackType == "service" {
		spotifyCallback = os.Getenv("SPOTIFY_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		spotifyCallback = os.Getenv("SPOTIFY_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://accounts.spotify.com/api/token", code, spotifyCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func GetDiscordOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var discordCallback string

	if callbackType == "service" {
		discordCallback = os.Getenv("DISCORD_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		discordCallback = os.Getenv("DISCORD_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://discord.com/api/v10/oauth2/token", code, discordCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func GetRedditOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var redditCallback string
	var redditClientID string
	var redditClientSecret string

	if callbackType == "service" {
		redditCallback = os.Getenv("REDDIT_SERVICE_CALLBACK")
		redditClientID = os.Getenv("REDDIT_SERVICE_CLIENT_ID")
		redditClientSecret = os.Getenv("REDDIT_SERVICE_CLIENT_SECRET")
	} else if callbackType == "login" {
		redditCallback = os.Getenv("REDDIT_LOGIN_CALLBACK")
		redditClientID = os.Getenv("REDDIT_LOGIN_CLIENT_ID")
		redditClientSecret = os.Getenv("REDDIT_LOGIN_CLIENT_SECRET")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://www.reddit.com/api/v1/access_token", code, redditCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(redditClientID, redditClientSecret)
	request.Header.Set("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func genericAccessTokenRequest(tokenUrl string, code string, callbackUrl string) (*http.Request, error) {
	jsonBody := fmt.Sprintf(
		"grant_type=%s&code=%s&redirect_uri=%s",
		grantTypeAuthorization,
		code,
		callbackUrl,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}
	return request, nil
}

func GetOAuth2UserEmailRequest(userInfoUrl string, accessToken string) (*http.Request, error) {
	request, errRequest := http.NewRequest("GET", userInfoUrl, nil)
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set(ContentType, ContentTypeUrlEncoded)
	return request, nil
}

func (self *ServiceService) GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://oauth2.googleapis.com/token"
	jsonBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=%s&refresh_token=%s",
		os.Getenv("GOOGLE_WEB_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		grantTypeRefreshToken,
		refreshToken,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://discord.com/api/v10/oauth2/token"
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetRedditRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://www.reddit.com/api/v1/access_token"
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("REDDIT_SERVICE_CLIENT_ID"), os.Getenv("REDDIT_SERVICE_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)
	request.Header.Set("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))

	return request, nil
}

func (self *ServiceService) genericRefreshTokenRequest(tokenUrl string, refreshToken string) (*http.Request, error) {
	jsonBody := fmt.Sprintf(
		"grant_type=%s&refresh_token=%s",
		grantTypeRefreshToken,
		refreshToken,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}
	return request, nil
}

func (self *ServiceService) ExecuteRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(APICallFailedMessage)
	}
	return res, nil
}

func (self *ServiceService) ExecuteApiRequest(url string, method string, typeToken string, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", typeToken+accessToken)
	req.Header.Set("Content-Type", "application/json")

	return self.ExecuteRequest(req)
}

func (self *ServiceService) GetResultTokenFromCode(code string, serviceName string, callbackType string, appType string) (entities.ResultToken, error) {
	var tokenRes entities.ResultToken
	var request *http.Request
	var errRequest error

	if serviceName == "Google" {
		request, errRequest = GetGoogleOAuth2AccessTokenRequest(code, callbackType, appType)
	} else if serviceName == "Spotify" {
		request, errRequest = GetSpotifyOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Discord" {
		request, errRequest = GetDiscordOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Github" {
		request, errRequest = GetGithubOAuth2AccessTokenRequest(code, callbackType)
		request.Header.Set("Accept", "application/json")
		tokenRes.ExpiresIn = oneYearSecond
	} else if serviceName == "Reddit" {
		request, errRequest = GetRedditOAuth2AccessTokenRequest(code, callbackType)
	} else {
		return tokenRes, fmt.Errorf("Unknown service")
	}

	if errRequest != nil {
		return tokenRes, errRequest
	}

	res, err := self.ExecuteRequest(request)
	if err != nil {
		return tokenRes, err
	}
	defer res.Body.Close()

	errDecoder := json.NewDecoder(res.Body).Decode(&tokenRes)
	if errDecoder != nil {
		return tokenRes, errDecoder
	}
	return tokenRes, nil
}

func decodeUserInfoWithGithub(res *http.Response) (entities.UserInfo, error) {
	var userInfos []entities.UserInfo

	err := json.NewDecoder(res.Body).Decode(&userInfos)
	if err != nil {
		return userInfos[0], err
	}
	return userInfos[0], nil
}

func (self *ServiceService) GetUserInfoFromService(accessToken string, serviceName string) (entities.UserInfo, error) {
	var userInfo entities.UserInfo
	var request *http.Request
	var err error

	if serviceName == "Google" {
		request, err = GetOAuth2UserEmailRequest("https://www.googleapis.com/oauth2/v3/userinfo", accessToken)
	} else if serviceName == "Spotify" {
		request, err = GetOAuth2UserEmailRequest("https://api.spotify.com/v1/me", accessToken)
	} else if serviceName == "Discord" {
		request, err = GetOAuth2UserEmailRequest("https://discord.com/api/users/@me", accessToken)
	} else if serviceName == "Github" {
		request, err = GetOAuth2UserEmailRequest("https://api.github.com/user/emails", accessToken)
	} else {
		return userInfo, fmt.Errorf("Unknown service")
	}
	if err != nil {
		return userInfo, err
	}

	res, err := self.ExecuteRequest(request)
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	if serviceName == "Github" {
		userInfo, err = decodeUserInfoWithGithub(res)
	} else {
		err = json.NewDecoder(res.Body).Decode(&userInfo)
	}
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (self *ServiceService) RequestToTimeApi() (entities.TimeResponse, error) {
	timeUrl := "https://tools.aimylogic.com/api/now?tz=Europe/Paris"
	var timeRes entities.TimeResponse

	request, errRequest := http.NewRequest("GET", timeUrl, nil)
	if errRequest != nil {
		return timeRes, errRequest
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return timeRes, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return timeRes, fmt.Errorf("HTTP request failed: %s", res.Status)
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&timeRes)
	if errDecoder != nil {
		return timeRes, errDecoder
	}
	return timeRes, nil
}
