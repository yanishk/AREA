package service_service

import (
	"backend/src/entities"
	"backend/src/storage"
	"bytes"
	"encoding/json"
	"fmt"
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
const grantTypeRefreshToken = "refresh_token"
const grantTypeAuthorization = "authorization_code"

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

func (self *ServiceService) FindActionByServiceId(serviceId string) ([]entities.Action, error) {
	return self.ActionRepository.FindActionByServiceId(serviceId)
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

func (self *ServiceService) FindWorkflowByActionId(actionId string) ([]entities.Workflow, error) {
	return self.WorkflowRepository.FindWorkflowByActionId(actionId)
}

func (self *ServiceService) FindReactionByID(id string) (entities.Reaction, error) {
	return self.ReactionRepository.FindReactionByID(id)
}

func (self *ServiceService) RetrieveActionsServices() ([]entities.ServiceUsefulInformations, error) {
	return self.ServiceRepository.FindActionsServices()
}

func (self *ServiceService) RetrieveReactionsServices() ([]entities.ServiceUsefulInformations, error) {
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

func GetSpotifyOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var spotifyCallback string

	if callbackType == "service" {
		spotifyCallback = os.Getenv("SPOTIFY_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		spotifyCallback = os.Getenv("SPOTIFY_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	tokenUrl := "https://accounts.spotify.com/api/token"
	jsonBody := fmt.Sprintf(
		"grant_type=%s&code=%s&redirect_uri=%s",
		grantTypeAuthorization,
		code,
		spotifyCallback,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
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

	tokenUrl := "https://discord.com/api/v10/oauth2/token"
	jsonBody := fmt.Sprintf(
		"grant_type=%s&code=%s&redirect_uri=%s",
		grantTypeAuthorization,
		code,
		discordCallback,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
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
		githubClientID = os.Getenv("GITHUB_SERVICE_CLIENT_ID")
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

func GetGoogleOAuth2UserEmailRequest(accessToken string) (*http.Request, error) {
	tokenUrl := "https://www.googleapis.com/oauth2/v3/userinfo"

	request, errRequest := http.NewRequest("GET", tokenUrl, nil)
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
	jsonBody := fmt.Sprintf(
		"grant_type=%s&refresh_token=%s",
		grantTypeRefreshToken,
		refreshToken,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://discord.com/api/v10/oauth2/token"
	jsonBody := fmt.Sprintf(
		"grant_type=%s&refresh_token=%s",
		grantTypeRefreshToken,
		refreshToken,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetGithubRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://github.com/login/oauth/access_token"
	jsonBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=%s&refresh_token=%s",
		os.Getenv("GITHUB_SERVICE_CLIENT_ID"),
		os.Getenv("GITHUB_SERVICE_CLIENT_SECRET"),
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

func (self *ServiceService) RequestService(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to reach requested service API")
	}
	return res, nil
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
	} else {
		return tokenRes, fmt.Errorf("Unknown service")
	}

	if errRequest != nil {
		return tokenRes, errRequest
	}

	res, err := self.RequestService(request)
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

func (self *ServiceService) GetUserInfoFromService(accessToken string, serviceName string) (entities.UserInfo, error) {
	var userInfo entities.UserInfo
	var request *http.Request
	var errRequest error

	if serviceName == "Google" {
		request, errRequest = GetGoogleOAuth2UserEmailRequest(accessToken)
	} else {
		return userInfo, fmt.Errorf("Unknown service")
	}

	if errRequest != nil {
		return userInfo, errRequest
	}

	res, err := self.RequestService(request)
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	errDecoder := json.NewDecoder(res.Body).Decode(&userInfo)
	if errDecoder != nil {
		return userInfo, errDecoder
	}
	return userInfo, nil
}
