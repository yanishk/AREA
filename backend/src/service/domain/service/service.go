package service_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"backend/src/entities"
	"backend/src/storage"
)

type ServiceService struct {
	ServiceRepository  storage.ServiceRepository
	UserRepository     storage.UserRepository
	ActionRepository   storage.ActionRepository
	WorkflowRepository storage.WorkflowRepository
	ReactionRepository storage.ReactionRepository
}

const contentType = "Content-Type"
const contentTypeUrlEncoded = "application/x-www-form-urlencoded"

const invalidCallbackTypeMessage = "Invalid callback type"
const unknownServiceMessage = "Unknown service"
const apiCallFailedMessage = "API call failed"

const grantTypeRefreshToken = "refresh_token"
const grantTypeAuthorization = "authorization_code"
const oneYearSecond = 31536000

const clientIdParam = "client_id="
const clientSecretParam = "&client_secret="
const codeParam = "&code="
const grantTypeParam = "&grant_type="
const redirectUriParam = "&redirect_uri="
const stateParam = "&state="
const codeResponseType = "&response_type=code"

const githubBaseUrl = "https://api.github.com/"

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

func (self *ServiceService) FindServiceById(id string) (entities.Service, error) {
	return self.ServiceRepository.FindServiceById(id)
}

func (self *ServiceService) FindServiceByActionId(id string) (entities.Service, error) {
	action, err := self.ActionRepository.FindActionById(id)
	if err != nil {
		return entities.Service{}, err
	}
	return self.ServiceRepository.FindServiceById(action.ServiceId)
}

func (self *ServiceService) FindServiceByReactionId(id string) (entities.Service, error) {
	reaction, err := self.ReactionRepository.FindReactionById(id)
	if err != nil {
		return entities.Service{}, err
	}
	return self.ServiceRepository.FindServiceById(reaction.ServiceId)
}

func (self *ServiceService) RetrieveActionsFromService(serviceName string) ([]entities.Action, error) {
	foundService, errFoundService := self.ServiceRepository.FindServiceByName(serviceName)
	if errFoundService != nil {
		return nil, errFoundService
	}

	actions, errFindActions := self.ActionRepository.FindActionsByServiceId(foundService.Id)
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

	reactions, errFindReactions := self.ReactionRepository.FindReactionsByServiceId(foundService.Id)
	if errFindReactions != nil {
		return nil, errFindReactions
	}
	return reactions, errFindReactions
}

func (self *ServiceService) RetrieveActionsServices() ([]entities.Service, error) {
	return self.ServiceRepository.FindActionsServices()
}

func (self *ServiceService) RetrieveReactionsServices() ([]entities.Service, error) {
	return self.ServiceRepository.FindReactionsServices()
}

func getCallbackAndClientId(callbackType, serviceName string, isIdNecessary bool) (string, string) {
	var callbackLink, clientID string

	if callbackType == "login" {
		if isIdNecessary {
			clientID = os.Getenv(serviceName + "_LOGIN_CLIENT_ID")
		}
		callbackLink = os.Getenv(serviceName + "_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		if isIdNecessary {
			clientID = os.Getenv(serviceName + "_SERVICE_CLIENT_ID")
		}
		callbackLink = os.Getenv(serviceName + "_SERVICE_CALLBACK")
	}
	return callbackLink, clientID
}

func oauth2Google(callbackType, appType string) string {
	var callbackLink, appTypeLink string

	if callbackType == "login" {
		callbackLink = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		callbackLink = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	}

	if appType == "web" {
		appTypeLink = os.Getenv("GOOGLE_WEB_CLIENT_ID")
	} else if appType == "mobile" {
		appTypeLink = os.Getenv("GOOGLE_MOBILE_CLIENT_ID")
	}

	return ("https://accounts.google.com/o/oauth2/auth?" +
		clientIdParam + appTypeLink +
		codeResponseType +
		"&access_type=offline" +
		redirectUriParam + callbackLink +
		"&scope=email")
}

func oauth2Spotify(callbackType string) string {
	callbackLink, _ := getCallbackAndClientId(callbackType, "SPOTIFY", false)

	return ("https://accounts.spotify.com/authorize?" +
		clientIdParam + os.Getenv("SPOTIFY_CLIENT_ID") +
		codeResponseType +
		redirectUriParam + callbackLink +
		"&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state user-library-modify")
}

func oauth2Discord(callbackType string) string {
	callbackLink, _ := getCallbackAndClientId(callbackType, "DISCORD", false)
	var scope string

	if callbackType == "login" {
		scope = "identify+email"
	} else if callbackType == "service" {
		scope = "identify+email+guilds+bot"
	}

	return ("https://discord.com/oauth2/authorize?" +
		clientIdParam + os.Getenv("DISCORD_CLIENT_ID") +
		"&permissions=141376" +
		codeResponseType +
		redirectUriParam + callbackLink +
		"&scope=" + scope)
}

func oauth2Github(callbackType string) string {
	callbackLink, clientID := getCallbackAndClientId(callbackType, "GITHUB", true)

	return ("https://github.com/login/oauth/authorize?" +
		clientIdParam + clientID +
		redirectUriParam + callbackLink +
		"&scope=repo admin:org user" +
		stateParam + os.Getenv("GITHUB_STATE"))
}

func oauth2Reddit(callbackType string) string {
	callbackLink, clientID := getCallbackAndClientId(callbackType, "REDDIT", true)

	return ("https://www.reddit.com/api/v1/authorize?" +
		clientIdParam + clientID +
		codeResponseType +
		stateParam + os.Getenv("REDDIT_STATE") +
		redirectUriParam + callbackLink +
		"&duration=permanent" +
		"&scope=identity read submit vote mysubreddits history account edit privatemessages")
}

func oauth2Asana(callbackType string) string {
	callbackLink, _ := getCallbackAndClientId(callbackType, "ASANA", false)

	return ("https://app.asana.com/-/oauth_authorize?" +
		clientIdParam + os.Getenv("ASANA_CLIENT_ID") +
		redirectUriParam + callbackLink +
		codeResponseType +
		stateParam + os.Getenv("ASANA_STATE") +
		"&scope=default email profile")
}

func oauth2Linkedin(callbackType string) string {
	callbackLink, _ := getCallbackAndClientId(callbackType, "LINKEDIN", false)

	return ("https://www.linkedin.com/oauth/v2/authorization?" +
		clientIdParam + url.QueryEscape(os.Getenv("LINKEDIN_CLIENT_ID")) +
		redirectUriParam + callbackLink +
		codeResponseType +
		"&scope=openid email profile w_member_social") +
		stateParam + os.Getenv("LINKEDIN_STATE")
}

func oauth2Dropbox(callbackType string) string {
	callbackLink, _ := getCallbackAndClientId(callbackType, "DROPBOX", false)

	return ("https://www.dropbox.com/oauth2/authorize?" +
		clientIdParam + url.QueryEscape(os.Getenv("DROPBOX_CLIENT_ID")) +
		redirectUriParam + callbackLink +
		codeResponseType + "&token_access_type=offline")
}

func oauth2Gitlab(callbackType string) string {
	callbackLink, _ := getCallbackAndClientId(callbackType, "GITLAB", false)

	return ("https://gitlab.com/oauth/authorize?" +
		clientIdParam + os.Getenv("GITLAB_CLIENT_ID") +
		codeResponseType +
		redirectUriParam + callbackLink +
		stateParam + os.Getenv("GITLAB_STATE") +
		"&scope=api read_api read_user read_repository write_repository")
}

func (self *ServiceService) OAuth2Service(serviceName, callbackType, appType string) (string, error) {
	var authUrl string

	if serviceName == "Google" {
		authUrl = oauth2Google(callbackType, appType)
	} else if serviceName == "Spotify" {
		authUrl = oauth2Spotify(callbackType)
	} else if serviceName == "Discord" {
		authUrl = oauth2Discord(callbackType)
	} else if serviceName == "Github" {
		authUrl = oauth2Github(callbackType)
	} else if serviceName == "Reddit" {
		authUrl = oauth2Reddit(callbackType)
	} else if serviceName == "Asana" {
		authUrl = oauth2Asana(callbackType)
	} else if serviceName == "Linkedin" {
		authUrl = oauth2Linkedin(callbackType)
	} else if serviceName == "Dropbox" {
		authUrl = oauth2Dropbox(callbackType)
	} else if serviceName == "Gitlab" {
		authUrl = oauth2Gitlab(callbackType)
	} else {
		return authUrl, fmt.Errorf(unknownServiceMessage)
	}
	return authUrl, nil
}

func genericAccessTokenRequest(tokenUrl, code, callbackUrl string) (*http.Request, error) {
	jsonBody := fmt.Sprintf(
		"grant_type=%s&code=%s&redirect_uri=%s",
		grantTypeAuthorization,
		code,
		callbackUrl,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}
	return request, nil
}

func getGoogleOAuth2AccessTokenWebRequestBody(code, callback string) string {
	return fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s",
		code,
		os.Getenv("GOOGLE_WEB_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		callback,
		grantTypeAuthorization,
	)
}

func getGoogleOAuth2AccessTokenMobileRequestBody(code, callback string) string {
	return fmt.Sprintf(
		"code=%s&client_id=%s&redirect_uri=%s&grant_type=%s",
		code,
		os.Getenv("GOOGLE_MOBILE_CLIENT_ID"),
		callback,
		grantTypeAuthorization,
	)
}

func getGoogleOAuth2AccessTokenRequest(code, callbackType, appType string) (*http.Request, error) {
	var googleCallback, jsonBody string

	if callbackType == "service" {
		googleCallback = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		googleCallback = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
	}

	if appType == "web" {
		jsonBody = getGoogleOAuth2AccessTokenWebRequestBody(code, googleCallback)
	} else if appType == "mobile" {
		jsonBody = getGoogleOAuth2AccessTokenMobileRequestBody(code, googleCallback)
	} else {
		return nil, fmt.Errorf("Invalid app type")
	}

	tokenUrl := "https://oauth2.googleapis.com/token"
	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(contentType, contentTypeUrlEncoded)
	return request, nil
}

func getGithubOAuth2AccessTokenRequest(code, callbackType string) (*http.Request, error) {
	var githubCallback, githubClientID, githubClientSecret string

	if callbackType == "service" {
		githubCallback = os.Getenv("GITHUB_SERVICE_CALLBACK")
		githubClientID = os.Getenv("GITHUB_SERVICE_CLIENT_ID")
		githubClientSecret = os.Getenv("GITHUB_SERVICE_CLIENT_SECRET")
	} else if callbackType == "login" {
		githubCallback = os.Getenv("GITHUB_LOGIN_CALLBACK")
		githubClientID = os.Getenv("GITHUB_LOGIN_CLIENT_ID")
		githubClientSecret = os.Getenv("GITHUB_LOGIN_CLIENT_SECRET")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
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

	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func getGitlabOAuth2AccessTokenRequest(code, callbackType string) (*http.Request, error) {
	var gitlabCallback string

	if callbackType == "service" {
		gitlabCallback = os.Getenv("GITLAB_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		gitlabCallback = os.Getenv("GITLAB_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
	}

	tokenUrl := "https://gitlab.com/oauth/token"
	jsonBody := fmt.Sprintf(
		clientIdParam + os.Getenv("GITLAB_CLIENT_ID") +
			clientSecretParam + os.Getenv("GITLAB_CLIENT_SECRET") +
			codeParam + code +
			redirectUriParam + gitlabCallback +
			grantTypeParam + grantTypeAuthorization,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func getAsanaOAuth2AccessTokenRequest(code, callbackType string) (*http.Request, error) {
	var asanaCallback string

	if callbackType == "service" {
		asanaCallback = os.Getenv("ASANA_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		asanaCallback = os.Getenv("ASANA_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
	}

	tokenUrl := "https://app.asana.com/-/oauth_token"
	jsonBody := fmt.Sprintf(
		clientIdParam + os.Getenv("ASANA_CLIENT_ID") +
			clientSecretParam + os.Getenv("ASANA_CLIENT_SECRET") +
			codeParam + code +
			redirectUriParam + asanaCallback +
			grantTypeParam + grantTypeAuthorization,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func getLinkedinOAuth2AccessTokenRequest(code, callbackType string) (*http.Request, error) {
	var callback string

	if callbackType == "service" {
		callback = os.Getenv("LINKEDIN_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		callback = os.Getenv("LINKEDIN_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
	}

	tokenUrl := "https://www.linkedin.com/oauth/v2/accessToken"
	jsonBody := fmt.Sprintf(
		clientIdParam + url.QueryEscape(os.Getenv("LINKEDIN_CLIENT_ID")) +
			clientSecretParam + url.QueryEscape(os.Getenv("LINKEDIN_CLIENT_SECRET")) +
			codeParam + code +
			redirectUriParam + callback +
			grantTypeParam + grantTypeAuthorization,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func getDropboxOAuth2AccessTokenRequest(code, callbackType string) (*http.Request, error) {
	var dropboxCallback string

	if callbackType == "service" {
		dropboxCallback = os.Getenv("DROPBOX_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		dropboxCallback = os.Getenv("DROPBOX_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://api.dropboxapi.com/oauth2/token", code, dropboxCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("DROPBOX_CLIENT_ID"), os.Getenv("DROPBOX_CLIENT_SECRET"))
	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func getSpotifyOAuth2AccessTokenRequest(code, callbackType string) (*http.Request, error) {
	var spotifyCallback string

	if callbackType == "service" {
		spotifyCallback = os.Getenv("SPOTIFY_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		spotifyCallback = os.Getenv("SPOTIFY_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://accounts.spotify.com/api/token", code, spotifyCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func getDiscordOAuth2AccessTokenRequest(code, callbackType string) (*http.Request, error) {
	var discordCallback string

	if callbackType == "service" {
		discordCallback = os.Getenv("DISCORD_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		discordCallback = os.Getenv("DISCORD_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://discord.com/api/v10/oauth2/token", code, discordCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func getRedditOAuth2AccessTokenRequest(code, callbackType string) (*http.Request, error) {
	var redditCallback, redditClientID, redditClientSecret string

	if callbackType == "service" {
		redditCallback = os.Getenv("REDDIT_SERVICE_CALLBACK")
		redditClientID = os.Getenv("REDDIT_SERVICE_CLIENT_ID")
		redditClientSecret = os.Getenv("REDDIT_SERVICE_CLIENT_SECRET")
	} else if callbackType == "login" {
		redditCallback = os.Getenv("REDDIT_LOGIN_CALLBACK")
		redditClientID = os.Getenv("REDDIT_LOGIN_CLIENT_ID")
		redditClientSecret = os.Getenv("REDDIT_LOGIN_CLIENT_SECRET")
	} else {
		return nil, fmt.Errorf(invalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://www.reddit.com/api/v1/access_token", code, redditCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(redditClientID, redditClientSecret)
	request.Header.Set("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))
	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) genericRefreshTokenRequest(tokenUrl, refreshToken string) (*http.Request, error) {
	jsonBody := fmt.Sprintf(
		"grant_type=%s&refresh_token=%s",
		grantTypeRefreshToken,
		refreshToken,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}
	return request, nil
}

func (self *ServiceService) genericJsonBodyRefreshToken(clientId, secretId, refreshToken string) string {
	return fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=%s&refresh_token=%s",
		clientId,
		secretId,
		grantTypeRefreshToken,
		refreshToken,
	)
}

func (self *ServiceService) GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://oauth2.googleapis.com/token"
	jsonBody := self.genericJsonBodyRefreshToken(os.Getenv("GOOGLE_WEB_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), refreshToken)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetAsanaRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://app.asana.com/-/oauth_token"
	jsonBody := self.genericJsonBodyRefreshToken(os.Getenv("ASANA_CLIENT_ID"), os.Getenv("ASANA_CLIENT_SECRET"), refreshToken)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetDropboxRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://api.dropbox.com/oauth2/token"
	jsonBody := self.genericJsonBodyRefreshToken(os.Getenv("DROPBOX_CLIENT_ID"), os.Getenv("DROPBOX_CLIENT_SECRET"), refreshToken)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://discord.com/api/v10/oauth2/token"
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetRedditRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://www.reddit.com/api/v1/access_token"
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("REDDIT_SERVICE_CLIENT_ID"), os.Getenv("REDDIT_SERVICE_CLIENT_SECRET"))
	request.Header.Set(contentType, contentTypeUrlEncoded)
	request.Header.Set("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))

	return request, nil
}

func (self *ServiceService) GetGitlabRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://gitlab.example.com/oauth/token"
	jsonBody := self.genericJsonBodyRefreshToken(os.Getenv("GITLAB_CLIENT_ID"), os.Getenv("GITLAB_CLIENT_SECRET"), refreshToken)
	jsonBody += "&redirect_uri=" + os.Getenv("GITLAB_SERVICE_CALLBACK")

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(contentType, contentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) ExecuteRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf(apiCallFailedMessage)
	}
	return res, nil
}

func (self *ServiceService) ExecuteApiRequest(url, method, typeToken, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", typeToken+accessToken)
	req.Header.Set("Content-Type", "application/json")

	return self.ExecuteRequest(req)
}

func (self *ServiceService) GetResultTokenFromCode(code, serviceName, callbackType, appType string) (entities.ResultToken, error) {
	var tokenRes entities.ResultToken
	var request *http.Request
	var errRequest error

	if serviceName == "Google" {
		request, errRequest = getGoogleOAuth2AccessTokenRequest(code, callbackType, appType)
	} else if serviceName == "Spotify" {
		request, errRequest = getSpotifyOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Discord" {
		request, errRequest = getDiscordOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Github" {
		request, errRequest = getGithubOAuth2AccessTokenRequest(code, callbackType)
		request.Header.Set("Accept", "application/json")
		tokenRes.ExpiresIn = oneYearSecond
	} else if serviceName == "Reddit" {
		codeStripped := strings.TrimSuffix(code, "#_")
		request, errRequest = getRedditOAuth2AccessTokenRequest(codeStripped, callbackType)
	} else if serviceName == "Asana" {
		request, errRequest = getAsanaOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Linkedin" {
		request, errRequest = getLinkedinOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Dropbox" {
		request, errRequest = getDropboxOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Gitlab" {
		request, errRequest = getGitlabOAuth2AccessTokenRequest(code, callbackType)
	} else {
		return tokenRes, fmt.Errorf(unknownServiceMessage)
	}

	if errRequest != nil {
		return tokenRes, errRequest
	}

	res, err := self.ExecuteRequest(request)
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

func getOAuth2UserEmailRequest(userInfoUrl, accessToken string) (*http.Request, error) {
	request, errRequest := http.NewRequest("GET", userInfoUrl, nil)
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set(contentType, contentTypeUrlEncoded)
	return request, nil
}

func decodeUserInfoWithGithub(res *http.Response) (entities.UserInfo, error) {
	var userInfos []entities.UserInfo

	err := json.NewDecoder(res.Body).Decode(&userInfos)
	if err != nil {
		return userInfos[0], err
	}
	return userInfos[0], nil
}

func (self *ServiceService) GetUserInfoFromService(accessToken, serviceName string) (entities.UserInfo, error) {
	var userInfo entities.UserInfo
	var request *http.Request
	var err error

	if serviceName == "Google" {
		request, err = getOAuth2UserEmailRequest("https://www.googleapis.com/oauth2/v3/userinfo", accessToken)
	} else if serviceName == "Spotify" {
		request, err = getOAuth2UserEmailRequest("https://api.spotify.com/v1/me", accessToken)
	} else if serviceName == "Discord" {
		request, err = getOAuth2UserEmailRequest("https://discord.com/api/users/@me", accessToken)
	} else if serviceName == "Github" {
		request, err = getOAuth2UserEmailRequest(githubBaseUrl+"user/emails", accessToken)
	} else {
		return userInfo, fmt.Errorf(unknownServiceMessage)
	}
	if err != nil {
		return userInfo, err
	}

	res, err := self.ExecuteRequest(request)
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	if serviceName == "Github" {
		userInfo, err = decodeUserInfoWithGithub(res)
	} else {
		err = json.NewDecoder(res.Body).Decode(&userInfo)
	}
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (self *ServiceService) RequestToTimeApi() (entities.TimeResponse, error) {
	timeUrl := "https://tools.aimylogic.com/api/now?tz=Europe/Paris"
	var timeRes entities.TimeResponse

	request, errRequest := http.NewRequest("GET", timeUrl, nil)
	if errRequest != nil {
		return timeRes, errRequest
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return timeRes, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return timeRes, fmt.Errorf("HTTP request failed: %s", res.Status)
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&timeRes)
	if errDecoder != nil {
		return timeRes, errDecoder
	}
	return timeRes, nil
}

func (self *ServiceService) RequestGithubUserRepositories(accessToken string) ([]entities.GithubRepository, error) {
	var repositories []entities.GithubRepository
	url := githubBaseUrl + "user/repos"

	resp, err := self.ExecuteApiRequest(url, "GET", "Bearer ", accessToken, nil)
	if err != nil {
		return repositories, err
	}

	err = json.NewDecoder(resp.Body).Decode(&repositories)
	if err != nil {
		return repositories, err
	}
	return repositories, nil
}

func (self *ServiceService) RetrieveDiscordGuildChannels(guildId string) ([]map[string]interface{}, error) {
	var channels []map[string]interface{}
	tokenBot := os.Getenv("DISCORD_BOT_TOKEN")
	getChannelUrl := "https://discord.com/api/v10/guilds/" + guildId + "/channels"

	resp, errRequest := self.ExecuteApiRequest(getChannelUrl, "GET", "Bot ", tokenBot, nil)
	if errRequest != nil {
		return channels, errRequest
	}
	defer resp.Body.Close()

	errDecoder := json.NewDecoder(resp.Body).Decode(&channels)
	if errDecoder != nil {
		return channels, errDecoder
	}
	return channels, nil
}
