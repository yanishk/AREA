package userservice_service

import (
	"backend/src/entities"
	"backend/src/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type UserServiceService struct {
	UserServiceRepository storage.UserServiceRepository
	ServiceRepository     storage.ServiceRepository
	UserRepository        storage.UserRepository
}

func NewUserServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository, UserServiceRepository storage.UserServiceRepository) *UserServiceService {
	return &UserServiceService{
		ServiceRepository:     ServiceRepository,
		UserRepository:        UserRepository,
		UserServiceRepository: UserServiceRepository,
	}
}

func (self *UserServiceService) FindServiceByName(name string) (entities.Service, error) {
	return self.ServiceRepository.FindServiceByName(name)
}

func (self *UserServiceService) GetUser(emailStr string) (entities.User, error) {
	foundUser, errSeekMail := self.UserRepository.FindUserByEmail(emailStr)
	if errSeekMail != nil || len(foundUser.Email) <= 0 {
		return foundUser, fmt.Errorf("could not find requested user")
	}
	return foundUser, nil
}

func (self *UserServiceService) CreateOrUpdateUserService(idService string, userId string, tokenSpotify entities.ResultToken) {
	var retrievedService entities.UserService

	currentTime := time.Now()
	expiryTime := currentTime.Add(time.Duration(tokenSpotify.ExpiresIn) * time.Second).Format("2006-01-02 15:04:05-07:00")

	retrievedService, errUserService := self.UserServiceRepository.FindUserServiceByIdServiceandIdUser(idService, userId)

	if errUserService == nil {
		self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(retrievedService.Id, userId, tokenSpotify.AccessToken, tokenSpotify.RefreshToken, expiryTime)
	} else {
		self.UserServiceRepository.CreateUserService(userId, tokenSpotify.AccessToken, tokenSpotify.RefreshToken, expiryTime, idService)
	}
}

func (self *UserServiceService) GetAccesTokenFromCode(codeSpotify string) (entities.ResultToken, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	jsonBody := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s",
		codeSpotify,
		os.Getenv("CALLBACK_SPOTIFY"),
	)
	var tokensRes entities.ResultToken

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return tokensRes, errRequest
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_ID_CLIENT"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return tokensRes, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return tokensRes, err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&tokensRes)
	if errDecoder != nil {
		return tokensRes, errDecoder
	}

	return tokensRes, nil
}

func (self *UserServiceService) FindUserServiceByIdServiceandIdUser(idservice string, userId string) (entities.UserService, error) {
	return self.UserServiceRepository.FindUserServiceByIdServiceandIdUser(idservice, userId)
}

func (self *UserServiceService) RefreshToken(refreshToken string, userUuid string) (entities.ResultToken, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	jsonBody := fmt.Sprintf(
		"grant_type=refresh_token&refresh_token=%s",
		refreshToken,
	)
	var tokensRefreshRes entities.ResultToken
	var serviceFound entities.Service

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return tokensRefreshRes, errRequest
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_ID_CLIENT"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return tokensRefreshRes, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return tokensRefreshRes, err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&tokensRefreshRes)
	if errDecoder != nil {
		return tokensRefreshRes, errDecoder
	}

	currentTime := time.Now()
	expiryTime := currentTime.Add(time.Duration(tokensRefreshRes.ExpiresIn) * time.Second).Format("2006-01-02 15:04:05-07:00")

	serviceFound, errServiceFound := self.ServiceRepository.FindServiceByName("Spotify")
	if errServiceFound != nil {
		return tokensRefreshRes, errServiceFound
	}

	self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(serviceFound.Id, userUuid, tokensRefreshRes.AccessToken, tokensRefreshRes.RefreshToken, expiryTime)
	return tokensRefreshRes, err
}
