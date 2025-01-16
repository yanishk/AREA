package user_service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"backend/src/entities"
	"backend/src/service"
	"backend/src/storage"
)

type UserService struct {
	UserRepository        storage.UserRepository
	ServiceRepository     storage.ServiceRepository
	UserServiceRepository storage.UserServiceRepository
	WorkflowRepository    storage.WorkflowRepository
	ServiceService        service.ServiceService
}

const basicConnectionType = "basic"

func NewUserService(UserRepository storage.UserRepository, ServiceRepository storage.ServiceRepository,
	UserServiceRepository storage.UserServiceRepository, WorkflowRepository storage.WorkflowRepository, ServiceService service.ServiceService) *UserService {
	return &UserService{
		UserRepository:        UserRepository,
		ServiceRepository:     ServiceRepository,
		UserServiceRepository: UserServiceRepository,
		WorkflowRepository:    WorkflowRepository,
		ServiceService:        ServiceService,
	}
}

func (self *UserService) CreateUser(userEmail, userPassword, userConnectionType string) error {
	var errorUser error

	if userConnectionType == basicConnectionType {
		hash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("Internal server error")
		}
		errorUser = self.UserRepository.CreateUser(userEmail, string(hash), userConnectionType)
	} else {
		_, errorService := self.ServiceRepository.FindServiceByName(userConnectionType)
		if errorService != nil {
			return fmt.Errorf("Connection type doesn't exist")
		}
		errorUser = self.UserRepository.CreateUser(userEmail, "", userConnectionType)
	}
	if errorUser != nil {
		return fmt.Errorf("Email address already used")
	}
	return nil
}

func createToken(email, connectionType string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email":          email,
			"connectionType": connectionType,
			"exp":            time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (self *UserService) LoginAuthentication(userEmail, userPassword, userConnectionType string) (string, error) {
	if userConnectionType != basicConnectionType && userPassword != "" {
		return "", fmt.Errorf("Wrong password")
	}

	foundUser, err := self.UserRepository.FindUserByEmail(userEmail, userConnectionType)
	if err != nil || len(foundUser.Email) <= 0 {
		return "", fmt.Errorf("Could not find requested user")
	}

	if userConnectionType == basicConnectionType {
		resBcrypt := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userPassword))
		if resBcrypt != nil {
			return "", fmt.Errorf("Wrong password")
		}
	}
	tokenString, errToken := createToken(foundUser.Email, foundUser.ConnectionType)
	if errToken != nil {
		return "", fmt.Errorf("Error creating token")
	}
	return tokenString, nil
}

func (self *UserService) LoginWithService(code, serviceName, appType string) (string, error) {
	_, errorService := self.ServiceRepository.FindServiceByName(serviceName)
	if errorService != nil {
		return "", errorService
	}

	resultToken, errToken := self.ServiceService.GetResultTokenFromCode(code, serviceName, "login", appType)
	if errToken != nil {
		return "", errToken
	}

	userInfo, errUserInfo := self.ServiceService.GetUserInfoFromService(resultToken.AccessToken, serviceName)
	if errUserInfo != nil {
		return "", errUserInfo
	}

	_, errFindingUser := self.UserRepository.FindUserByEmail(userInfo.Email, serviceName)
	if errFindingUser != nil {
		errCreatingUser := self.CreateUser(userInfo.Email, "", serviceName)
		if errCreatingUser != nil {
			return "", fmt.Errorf("Email address already used")
		}
	}

	token, errLogin := self.LoginAuthentication(userInfo.Email, "", serviceName)
	if errLogin != nil {
		return "", errLogin
	}
	return token, nil
}

func (self *UserService) GetUser(userEmail, userConnectionType string) (entities.UserInfos, error) {
	user, err := self.UserRepository.FindUserByEmail(userEmail, userConnectionType)
	if err != nil {
		return entities.UserInfos{}, err
	}

	return entities.UserInfos{Email: user.Email, CreatedAt: user.CreatedAt, Timezone: user.Timezone, ConnectionType: user.ConnectionType}, err
}

func (self *UserService) ModifyPassword(userEmail, userConnectionType string, newPassword entities.UserModifyPassword) error {
	if userConnectionType != basicConnectionType {
		return fmt.Errorf("Could not modify the password")
	}

	foundUser, errSeekMail := self.UserRepository.FindUserByEmail(userEmail, userConnectionType)
	if errSeekMail != nil || len(foundUser.Email) <= 0 {
		return fmt.Errorf("Could not find requested user")
	}

	resBcrypt := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(newPassword.OldPassword))
	if resBcrypt != nil {
		return fmt.Errorf("Old password is incorrect")
	}

	hash, errHash := bcrypt.GenerateFromPassword([]byte(newPassword.Password), bcrypt.DefaultCost)
	if errHash != nil {
		return fmt.Errorf("Failed to hash password")
	}

	errUpdate := self.UserRepository.UpdateUser(userEmail, string(hash), userConnectionType)
	if errUpdate != nil {
		return fmt.Errorf("Could not modify the password")
	}
	return nil
}

func (self *UserService) DeleteAccount(userEmail, userConnectionType string) error {
	user, err := self.UserRepository.FindUserByEmail(userEmail, userConnectionType)
	if err != nil {
		return err
	}

	err = self.WorkflowRepository.DeleteWorkflowByOwnerId(user.Id)
	if err != nil {
		return err
	}

	err = self.UserServiceRepository.DeleteUserServiceByUserId(user.Id)
	if err != nil {
		return err
	}

	err = self.UserRepository.DeleteUser(userEmail, userConnectionType)
	if err != nil {
		return fmt.Errorf("Could not delete account")
	}
	return nil
}

func (self *UserService) FindUserById(userId string) (entities.User, error) {
	return self.UserRepository.FindUserById(userId)
}
