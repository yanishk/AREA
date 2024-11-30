package service_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"backend/src/entities"
	"backend/src/storage"
)

type ServiceService struct {
	ServiceRepository storage.ServiceRepository
	UserRepository    storage.UserRepository
}

func NewServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository) *ServiceService {
	return &ServiceService{
		ServiceRepository: ServiceRepository,
		UserRepository:    UserRepository,
	}
}

func (self *ServiceService) GetUser(emailStr string) (entities.User, error) {
	foundUser, errSeekMail := self.UserRepository.FindUserByEmail(emailStr)
	if errSeekMail != nil || len(foundUser.Email) <= 0 {
		return foundUser, fmt.Errorf("could not find requested user")
	}
	return foundUser, nil
}

func (self *ServiceService) CreateOrUpdateSpotify(userId string, tokenSpotify entities.ResultToken) {
	currentTime := time.Now()
	expiryTime := currentTime.Add(time.Duration(tokenSpotify.ExpiresIn) * time.Second).Format("2006-01-02 15:04:05-07:00")

	_, errFindingService := self.ServiceRepository.FindUserServiceByName("spotify", userId)
	if errFindingService == nil {
		self.ServiceRepository.UpdateServiceByUserAndName("spotify", tokenSpotify.AccessToken, tokenSpotify.RefreshToken, userId, expiryTime)
	} else {
		self.ServiceRepository.CreateService("spotify", tokenSpotify.AccessToken, tokenSpotify.RefreshToken, userId, expiryTime)
	}
}

func (self *ServiceService) GetAccesTokenFromCode(codeSpotify string) (entities.ResultToken, error) {
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

func (self *ServiceService) FindUserServiceByName(name string, userId string) (entities.Service, error) {
	return self.ServiceRepository.FindUserServiceByName(name, userId)
}

func (self *ServiceService) RefreshToken(refreshToken string, userUuid string) (entities.ResultToken, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	jsonBody := fmt.Sprintf(
		"grant_type=refresh_token&refresh_token=%s",
		refreshToken,
	)
	var tokensRefreshRes entities.ResultToken

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

	fmt.Println("TOKEN REFRESH -> " + tokensRefreshRes.RefreshToken)
	self.ServiceRepository.UpdateServiceByUserAndName("spotify", tokensRefreshRes.AccessToken, tokensRefreshRes.RefreshToken, userUuid, expiryTime)
	return tokensRefreshRes, err
}
