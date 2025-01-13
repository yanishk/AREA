package userservice_service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"backend/src/entities"
	"backend/src/service"
	"backend/src/storage"
)

type UserServiceService struct {
	UserServiceRepository storage.UserServiceRepository
	ServiceRepository     storage.ServiceRepository
	UserRepository        storage.UserRepository
	ServiceService        service.ServiceService
}

const asanaWorkspaceRoute = "https://app.asana.com/api/1.0/workspaces/"

const formattingDate = "2006-01-02 15:04:05-07:00"

const bearerType = "Bearer "

func NewUserServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository,
	UserServiceRepository storage.UserServiceRepository, ServiceService service.ServiceService) *UserServiceService {
	return &UserServiceService{
		ServiceRepository:     ServiceRepository,
		UserRepository:        UserRepository,
		UserServiceRepository: UserServiceRepository,
		ServiceService:        ServiceService,
	}
}

func (self *UserServiceService) GetUser(email, connectionType string) (entities.User, error) {
	foundUser, errSeekMail := self.UserRepository.FindUserByEmail(email, connectionType)
	if errSeekMail != nil || len(foundUser.Email) <= 0 {
		return foundUser, fmt.Errorf("Could not find requested user")
	}
	return foundUser, nil
}

func (self *UserServiceService) RetrieveUserServiceAuthenticationStatus(email, connectionType, serviceName string) (bool, error) {
	user, err := self.UserRepository.FindUserByEmail(email, connectionType)
	if err != nil {
		return false, fmt.Errorf("Could not find requested user")
	}

	service, err := self.ServiceRepository.FindServiceByName(serviceName)
	if err != nil {
		return false, fmt.Errorf("Unknown service")
	}
	if !service.IsAuthNeeded {
		return true, nil
	}

	_, err = self.UserServiceRepository.FindUserServiceByServiceIdandUserId(user.Id, service.Id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (self *UserServiceService) refreshToken(refreshToken, userId, serviceName string) (entities.ResultToken, error) {
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
	} else if serviceName == "Reddit" {
		request, errRequest = self.ServiceService.GetRedditRefreshTokenRequest(refreshToken)
	} else if serviceName == "Asana" {
		request, errRequest = self.ServiceService.GetAsanaRefreshTokenRequest(refreshToken)
	} else if serviceName == "Dropbox" {
		request, errRequest = self.ServiceService.GetDropboxRefreshTokenRequest(refreshToken)
	} else if serviceName == "Gitlab" {
		request, errRequest = self.ServiceService.GetGitlabRefreshTokenRequest(refreshToken)
	} else {
		return tokenRes, fmt.Errorf("Unknown service")
	}

	if errRequest != nil {
		return tokenRes, errRequest
	}

	res, err := self.ServiceService.ExecuteRequest(request)
	if err != nil {
		return tokenRes, err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&tokenRes)
	if errDecoder != nil {
		return tokenRes, errDecoder
	}

	currentTime := time.Now()
	expiryDate := currentTime.Add(time.Duration(tokenRes.ExpiresIn) * time.Second).Format(formattingDate)

	serviceFound, errServiceFound := self.ServiceRepository.FindServiceByName(serviceName)
	if errServiceFound != nil {
		return tokenRes, errServiceFound
	}

	err = self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(userId, tokenRes.AccessToken, tokenRes.RefreshToken, expiryDate, serviceFound.Id)
	if err != nil {
		return tokenRes, err
	}
	return tokenRes, nil
}

func (self *UserServiceService) CallApiAndRefresh(email, connectionType, serviceName string) (string, error) {
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

	expiryDate, errTimestamp := time.Parse(formattingDate, foundUserService.ExpiryDate)
	if errTimestamp != nil {
		return "", errTimestamp
	}

	if expiryDate.Before(time.Now()) && serviceName != "Github" && serviceName != "Linkedin" {
		tokenFromRefresh, errRefresh := self.refreshToken(foundUserService.RefreshToken, foundUser.Id, serviceName)
		token = tokenFromRefresh.AccessToken
		if errRefresh != nil {
			return "", errRefresh
		}
	}
	return token, nil
}

func (self *UserServiceService) createOrUpdateUserService(userId, serviceId string, token entities.ResultToken) {
	currentTime := time.Now()
	expiryDate := currentTime.Add(time.Duration(token.ExpiresIn) * time.Second).Format(formattingDate)

	_, errUserService := self.UserServiceRepository.FindUserServiceByServiceIdandUserId(userId, serviceId)

	if errUserService == nil {
		self.UserServiceRepository.UpdateUserServiceByServiceIdAndUserId(userId, token.AccessToken, token.RefreshToken, expiryDate, serviceId)
	} else {
		self.UserServiceRepository.CreateUserService(userId, token.AccessToken, token.RefreshToken, expiryDate, serviceId)
	}
}

func (self *UserServiceService) UpdateTokenForService(code, serviceName, appType, email, connectionType string) error {
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

	self.createOrUpdateUserService(foundUser.Id, foundService.Id, token)
	return nil
}

func (self *UserServiceService) RetrieveGithubUserRepositories(email, connectionType string) ([]entities.GithubRepository, error) {
	accessToken, err := self.CallApiAndRefresh(email, connectionType, "Github")
	if err != nil {
		return []entities.GithubRepository{}, err
	}
	return self.ServiceService.RequestGithubUserRepositories(accessToken)
}

func (self *UserServiceService) RetrieveGitlabUserProjects(email, connectionType string) ([]entities.GitlabProject, error) {
	accessToken, err := self.CallApiAndRefresh(email, connectionType, "Gitlab")
	if err != nil {
		return []entities.GitlabProject{}, err
	}
	return self.ServiceService.RequestGitlabUserProjects(accessToken)
}

func (self *UserServiceService) RetrieveDiscordUserServers(email, connectionType string) ([]map[string]interface{}, error) {
	var ownedServers, guilds []map[string]interface{}

	accessToken, errRefreshing := self.CallApiAndRefresh(email, connectionType, "Discord")
	if errRefreshing != nil {
		return ownedServers, errRefreshing
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest("https://discord.com/api/v10/users/@me/guilds", "GET", bearerType, accessToken, nil)
	if errRequest != nil {
		return ownedServers, errRequest
	}
	defer resp.Body.Close()

	errDecoder := json.NewDecoder(resp.Body).Decode(&guilds)
	if errDecoder != nil {
		return ownedServers, errDecoder
	}

	for _, guild := range guilds {
		if guild["owner"] == true {
			ownedServers = append(ownedServers, guild)
		}
	}
	return ownedServers, nil
}

func (self *UserServiceService) RetrieveAsanaUserWorkspaces(email, connectionType string) (entities.AsanaWorkspacesInfo, error) {
	var workspaces entities.AsanaWorkspacesInfo
	accessToken, errRefreshing := self.CallApiAndRefresh(email, connectionType, "Asana")
	if errRefreshing != nil {
		return workspaces, errRefreshing
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(asanaWorkspaceRoute, "GET", bearerType, accessToken, nil)
	if errRequest != nil {
		return workspaces, errRequest
	}
	defer resp.Body.Close()

	errDecodeWorkspace := json.NewDecoder(resp.Body).Decode(&workspaces)
	if errDecodeWorkspace != nil {
		return workspaces, errDecodeWorkspace
	}
	return workspaces, nil
}

func (self *UserServiceService) decodeRequiredWorkspaceInfo(email, connectionType, url string) (entities.AsanaWorkspacesInfo, error) {
	var info entities.AsanaWorkspacesInfo

	accessToken, errRefreshing := self.CallApiAndRefresh(email, connectionType, "Asana")
	if errRefreshing != nil {
		return info, errRefreshing
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(url, "GET", bearerType, accessToken, nil)
	if errRequest != nil {
		return info, errRequest
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&info)
	if errDecode != nil || len(info.Data) == 0 {
		return info, errDecode
	}
	return info, nil
}

func (self *UserServiceService) RetrieveAsanaWorkspaceAssignees(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	url := asanaWorkspaceRoute + workspaceId + "/users"
	return self.decodeRequiredWorkspaceInfo(email, connectionType, url)
}

func (self *UserServiceService) RetrieveAsanaWorkspaceProjects(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	url := asanaWorkspaceRoute + workspaceId + "/projects"
	return self.decodeRequiredWorkspaceInfo(email, connectionType, url)
}

func (self *UserServiceService) RetrieveAsanaWorkspaceTags(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	url := asanaWorkspaceRoute + workspaceId + "/tags"
	return self.decodeRequiredWorkspaceInfo(email, connectionType, url)
}
