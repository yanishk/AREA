package userservice_service

import (
	"backend/src/entities"
	"backend/src/service"
	"backend/src/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserServiceService struct {
	UserServiceRepository storage.UserServiceRepository
	ServiceRepository     storage.ServiceRepository
	UserRepository        storage.UserRepository
	ServiceService        service.ServiceService
}

const FormattingDate = "2006-01-02 15:04:05-07:00"

func NewUserServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository,
	UserServiceRepository storage.UserServiceRepository, ServiceService service.ServiceService) *UserServiceService {
	return &UserServiceService{
		ServiceRepository:     ServiceRepository,
		UserRepository:        UserRepository,
		UserServiceRepository: UserServiceRepository,
		ServiceService:        ServiceService,
	}
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

func (self *UserServiceService) GetConnectionTypeFromContext(context *gin.Context) (string, bool) {
	connectionType, err := context.Get("connectionType")
	if !err {
		return "", err
	}

	connectionTypeStr, errConnectionType := connectionType.(string)
	if !errConnectionType {
		return "", errConnectionType
	}
	return connectionTypeStr, true
}

func (self *UserServiceService) GetUser(emailStr string, connectionTypeStr string) (entities.User, error) {
	foundUser, errSeekMail := self.UserRepository.FindUserByEmail(emailStr, connectionTypeStr)
	if errSeekMail != nil || len(foundUser.Email) <= 0 {
		return foundUser, fmt.Errorf("Could not find requested user")
	}
	return foundUser, nil
}

func (self *UserServiceService) RefreshToken(refreshToken string, userId string, serviceName string) (entities.ResultToken, error) {
	var tokenRes entities.ResultToken
	var serviceFound entities.Service
	var request *http.Request
	var errRequest error

	if serviceName == "Google" {
		request, errRequest = self.ServiceService.GetGoogleRefreshTokenRequest(refreshToken)
	} else if serviceName == "Spotify" {
		request, errRequest = self.ServiceService.GetSpotifyRefreshTokenRequest(refreshToken)
	} else if serviceName == "Discord" {
		request, errRequest = self.ServiceService.GetDiscordRefreshTokenRequest(refreshToken)
	} else {
		return tokenRes, fmt.Errorf("Unknown service")
	}

	if errRequest != nil {
		return tokenRes, errRequest
	}

	res, err := self.ServiceService.RequestService(request)
	if err != nil {
		return tokenRes, err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&tokenRes)
	if errDecoder != nil {
		return tokenRes, errDecoder
	}

	currentTime := time.Now()
	expiryDate := currentTime.Add(time.Duration(tokenRes.ExpiresIn) * time.Second).Format(FormattingDate)

	serviceFound, errServiceFound := self.ServiceRepository.FindServiceByName(serviceName)
	if errServiceFound != nil {
		return tokenRes, errServiceFound
	}

	self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(userId, tokenRes.AccessToken, tokenRes.RefreshToken, expiryDate, serviceFound.Id)
	return tokenRes, err
}

func (self *UserServiceService) CallApiAndRefresh(email string, connectionType string, serviceName string) (string, error) {
	foundUser, errGetUser := self.GetUser(email, connectionType)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		return "", errGetUser
	}

	foundService, errServiceFound := self.ServiceRepository.FindServiceByName(serviceName)
	if errServiceFound != nil {
		return "", errServiceFound
	}

	foundUserService, errUserService := self.UserServiceRepository.FindUserServiceByServiceIdandUserId(foundUser.Id, foundService.Id)
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

func (self *UserServiceService) UpdateTokenForService(code string, serviceName string, appType string, email string, connectionType string) error {
	_, errorService := self.ServiceRepository.FindServiceByName(serviceName)
	if errorService != nil {
		return errorService
	}

	token, errToken := self.ServiceService.GetResultTokenFromCode(code, serviceName, "service", appType)
	if errToken != nil {
		return errToken
	}

	foundUser, errGetUser := self.GetUser(email, connectionType)
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
