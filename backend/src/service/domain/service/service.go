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

func GetGoogleOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var googleCallback string

	if callbackType == "service" {
		googleCallback = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		googleCallback = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf("Invalid callback type")
	}

	tokenUrl := "https://oauth2.googleapis.com/token"
	jsonBody := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code,
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		googleCallback,
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
		return nil, fmt.Errorf("Invalid callback type")
	}

	tokenUrl := "https://accounts.spotify.com/api/token"
	jsonBody := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s",
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
		"client_id=%s&client_secret=%s&grant_type=refresh_token&refresh_token=%s",
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
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
		"grant_type=refresh_token&refresh_token=%s",
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

func (self *ServiceService) GetResultTokenFromCode(code string, serviceName string, callbackType string) (entities.ResultToken, error) {
	var tokenRes entities.ResultToken
	var request *http.Request
	var errRequest error

	if serviceName == "Google" {
		request, errRequest = GetGoogleOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Spotify" {
		request, errRequest = GetSpotifyOAuth2AccessTokenRequest(code, callbackType)
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
	} else if serviceName == "Spotify" {
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
