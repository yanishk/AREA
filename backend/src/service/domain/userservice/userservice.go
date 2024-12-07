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

func (self *UserServiceService) CreateOrUpdateUserService(idUser string, idService string, token entities.ResultToken) {
	currentTime := time.Now()
	expiryTime := currentTime.Add(time.Duration(token.ExpiresIn) * time.Second).Format("2006-01-02 15:04:05-07:00")

	_, errUserService := self.UserServiceRepository.FindUserServiceByIdServiceandIdUser(idUser, idService)

	if errUserService == nil {
		self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(idUser, token.AccessToken, token.RefreshToken, expiryTime, idService)
	} else {
		self.UserServiceRepository.CreateUserService(idUser, token.AccessToken, token.RefreshToken, expiryTime, idService)
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

func (self *UserServiceService) FindUserServiceByIdServiceandIdUser(idUser string, idService string) (entities.UserService, error) {
	return self.UserServiceRepository.FindUserServiceByIdServiceandIdUser(idUser, idService)
}

func (self *UserServiceService) RefreshToken(refreshToken string, userUuid string, serviceName string) (entities.ResultToken, error) {
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
	expiryTime := currentTime.Add(time.Duration(tokenRes.ExpiresIn) * time.Second).Format("2006-01-02 15:04:05-07:00")

	serviceFound, errServiceFound := self.ServiceRepository.FindServiceByName("Spotify")
	if errServiceFound != nil {
		return tokenRes, errServiceFound
	}

	self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(userUuid, tokenRes.AccessToken, tokenRes.RefreshToken, expiryTime, serviceFound.Id)
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

func (self *UserServiceService) SpotifyCallApi(email string) (string, error) {
	foundUser, errGetUser := self.GetUser(email)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		return "", errGetUser
	}

	foundService, errServiceFound := self.FindServiceByName("Spotify")
	if errServiceFound != nil {
		return "", errServiceFound
	}

	foundUserService, errUserService := self.FindUserServiceByIdServiceandIdUser(foundService.Id, foundUser.Id)
	token := foundUserService.AccessToken
	if errUserService != nil {
		return "", errUserService
	}

	expiryTime, errTimestamp := time.Parse("2006-01-02 15:04:05-07:00", foundUserService.ExpiryToken)
	if errTimestamp != nil {
		return "", errTimestamp
	}

	if expiryTime.Before(time.Now()) {
		tokenFromRefresh, errRefresh := self.RefreshToken(foundUserService.RefreshToken, foundUser.Id)
		token = tokenFromRefresh.AccessToken
		if errRefresh != nil {
			return "", errRefresh
		}
	}
	return token, nil
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

func (self *UserServiceService) SpotifyCallApi(email string) (string, error) {
	foundUser, errGetUser := self.GetUser(email)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		return "", errGetUser
	}

	foundService, errServiceFound := self.FindServiceByName("Spotify")
	if errServiceFound != nil {
		return "", errServiceFound
	}

	foundUserService, errUserService := self.FindUserServiceByIdServiceandIdUser(foundService.Id, foundUser.Id)
	token := foundUserService.AccessToken
	if errUserService != nil {
		return "", errUserService
	}

	expiryTime, errTimestamp := time.Parse("2006-01-02 15:04:05-07:00", foundUserService.ExpiryToken)
	if errTimestamp != nil {
		return "", errTimestamp
	}

	if expiryTime.Before(time.Now()) {
		tokenFromRefresh, errRefresh := self.RefreshToken(foundUserService.RefreshToken, foundUser.Id)
		token = tokenFromRefresh.AccessToken
		if errRefresh != nil {
			return "", errRefresh
		}
	}
	return token, nil
}
