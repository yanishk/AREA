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

	"github.com/gin-gonic/gin"
)

type UserServiceService struct {
	UserServiceRepository storage.UserServiceRepository
	ServiceRepository     storage.ServiceRepository
	UserRepository        storage.UserRepository
}

const ContentType = "Content-Type"
const ContentTypeUrlEncoded = "application/x-www-form-urlencoded"
const FormattingDate = "2006-01-02 15:04:05-07:00"

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
		return foundUser, fmt.Errorf("Could not find requested user")
	}
	return foundUser, nil
}

func (self *UserServiceService) CreateOrUpdateUserService(userId string, serviceId string, token entities.ResultToken) {
	currentTime := time.Now()
	expiryDate := currentTime.Add(time.Duration(token.ExpiresIn) * time.Second).Format(FormattingDate)

	_, errUserService := self.UserServiceRepository.FindUserServiceByServiceIdandUserId(userId, serviceId)

	if errUserService == nil {
		self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(userId, token.AccessToken, token.RefreshToken, expiryDate, serviceId)
	} else {
		self.UserServiceRepository.CreateUserService(userId, token.AccessToken, token.RefreshToken, expiryDate, serviceId)
	}
}

func GetGoogleOAuth2Request(code string) (*http.Request, error) {
	tokenUrl := "https://oauth2.googleapis.com/token"
	jsonBody := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code,
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_CALLBACK"),
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(ContentType, ContentTypeUrlEncoded)
	return request, nil
}

func GetSpotifyOAuth2Request(code string) (*http.Request, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	jsonBody := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s",
		code,
		os.Getenv("SPOTIFY_CALLBACK"),
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error) {
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

func GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error) {
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

func (self *UserServiceService) GetAccessTokenFromCode(code string, serviceName string) (entities.ResultToken, error) {
	var tokenRes entities.ResultToken
	var request *http.Request
	var errRequest error

	if serviceName == "Google" {
		request, errRequest = GetGoogleOAuth2Request(code)
	} else if serviceName == "Spotify" {
		request, errRequest = GetSpotifyOAuth2Request(code)
	}

	if errRequest != nil {
		return tokenRes, errRequest
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return tokenRes, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return tokenRes, err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&tokenRes)
	if errDecoder != nil {
		return tokenRes, errDecoder
	}
	return tokenRes, nil
}

func (self *UserServiceService) FindUserServiceByServiceIdandUserId(userId string, serviceId string) (entities.UserService, error) {
	return self.UserServiceRepository.FindUserServiceByServiceIdandUserId(userId, serviceId)
}

func (self *UserServiceService) RefreshToken(refreshToken string, userId string, serviceName string) (entities.ResultToken, error) {
	var tokenRes entities.ResultToken
	var serviceFound entities.Service
	var request *http.Request
	var errRequest error

	if serviceName == "Google" {
		request, errRequest = GetGoogleRefreshTokenRequest(refreshToken)
	} else if serviceName == "Spotify" {
		request, errRequest = GetSpotifyRefreshTokenRequest(refreshToken)
	}

	if errRequest != nil {
		return tokenRes, errRequest
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return tokenRes, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return tokenRes, err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&tokenRes)
	if errDecoder != nil {
		return tokenRes, errDecoder
	}

	currentTime := time.Now()
	expiryDate := currentTime.Add(time.Duration(tokenRes.ExpiresIn) * time.Second).Format(FormattingDate)

	serviceFound, errServiceFound := self.ServiceRepository.FindServiceByName("Spotify")
	if errServiceFound != nil {
		return tokenRes, errServiceFound
	}

	self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(userId, tokenRes.AccessToken, tokenRes.RefreshToken, expiryDate, serviceFound.Id)
	return tokenRes, err
}

func (self *UserServiceService) GetEmailFromContext(context *gin.Context) (string, bool) {
	email, err := context.Get("email")

	if !err {
		return "", err
	}

	emailStr, errEmail := email.(string)
	if !errEmail {
		return "", errEmail
	}
	return emailStr, true
}

func (self *UserServiceService) CallApiAndRefresh(email string, serviceName string) (string, error) {
	foundUser, errGetUser := self.GetUser(email)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		return "", errGetUser
	}

	foundService, errServiceFound := self.FindServiceByName(serviceName)
	if errServiceFound != nil {
		return "", errServiceFound
	}

	foundUserService, errUserService := self.FindUserServiceByServiceIdandUserId(foundUser.Id, foundService.Id)
	token := foundUserService.AccessToken
	if errUserService != nil {
		return "", errUserService
	}

	expiryDate, errTimestamp := time.Parse(FormattingDate, foundUserService.ExpiryDate)
	if errTimestamp != nil {
		return "", errTimestamp
	}

	if expiryDate.Before(time.Now()) {
		tokenFromRefresh, errRefresh := self.RefreshToken(foundUserService.RefreshToken, foundUser.Id, serviceName)
		token = tokenFromRefresh.AccessToken
		if errRefresh != nil {
			return "", errRefresh
		}
	}
	return token, nil
}

func (self *UserServiceService) UpdateTokenForService(code string, serviceName string, email string) error {
	_, errorServcie := self.ServiceRepository.FindServiceByName(serviceName)
	if errorServcie != nil {
		return errorServcie
	}

	token, errToken := self.GetAccessTokenFromCode(code, serviceName)
	if errToken != nil {
		return errToken
	}

	foundUser, errGetUser := self.GetUser(email)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		return errGetUser
	}

	foundService, errFoundService := self.ServiceRepository.FindServiceByName(serviceName)
	if errFoundService != nil {
		return errFoundService
	}

	self.CreateOrUpdateUserService(foundUser.Id, foundService.Id, token)
	return nil
}
