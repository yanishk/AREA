package user_service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"backend/src/entities"
	"backend/src/storage"
)

type UserService struct {
	UserRepository storage.UserRepository
	ServiceRepository  storage.ServiceRepository
}

func NewUserService(UserRepository storage.UserRepository, ServiceRepository storage.ServiceRepository) *UserService {
	return &UserService{
		UserRepository: UserRepository,
		ServiceRepository: ServiceRepository,
	}
}

func (self *UserService) CreateUser(userEmail string, userPassword string, userConnectionType string) error {
	var errorUser error

	if (userConnectionType == "basic") {
		hash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("Internal server error")
		}
		errorUser = self.UserRepository.CreateUser(userEmail, string(hash), userConnectionType)
	} else {
		_ , errorService := self.ServiceRepository.FindServiceByName(userConnectionType)
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

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (self *UserService) LoginAuthentication(userEmail string, userPassword string) (string, error) {
	foundUser, err := self.UserRepository.FindUserByEmail(userEmail)
	if err != nil || len(foundUser.Email) <= 0 {
		return "", fmt.Errorf("Could not find requested user")
	}

	resBcrypt := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userPassword))
	if resBcrypt == nil {
		tokenString, errToken := createToken(foundUser.Email)
		if errToken != nil {
			return "", fmt.Errorf("Error creating token")
		}
		return tokenString, nil
	}
	return "", fmt.Errorf("Wrong password")
}

func (self *UserService) ModifyPassword(userEmail string, newPassword entities.UserModifyPassword) error {
	foundUser, errSeekMail := self.UserRepository.FindUserByEmail(userEmail)
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

	errUpdate := self.UserRepository.UpdateUser(userEmail, string(hash))
	if errUpdate != nil {
		return fmt.Errorf("Could not modify the password")
	}
	return nil
}

func (self *UserService) DeleteAccount(userEmail string) error {
	err := self.UserRepository.DeleteUser(userEmail)
	if err != nil {
		return fmt.Errorf("Could not delete account")
	}
	return nil
}

func (self *UserService) FindUserById(userId string) (entities.User, error) {
	return self.UserRepository.FindUserById(userId)
}
