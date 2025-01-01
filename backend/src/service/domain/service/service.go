package service_service

import (
	"backend/src/entities"
	"backend/src/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
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

const InvalidCallbackTypeMessage = "Invalid callback type"
const UnknownServiceMessage = "Unknown service"
const APICallFailedMessage = "API call failed"

const grantTypeRefreshToken = "refresh_token"
const grantTypeAuthorization = "authorization_code"
const oneYearSecond = 31536000

const clientIdParam = "client_id="
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

func (self *ServiceService) OAuth2Service(serviceName string, callbackType string, appType string) (string, error) {
	var authUrl string

	if serviceName == "Google" {
		authUrl = self.OAuth2Google(callbackType, appType)
	} else if serviceName == "Spotify" {
		authUrl = self.OAuth2Spotify(callbackType)
	} else if serviceName == "Discord" {
		authUrl = self.OAuth2Discord(callbackType)
	} else if serviceName == "Github" {
		authUrl = self.OAuth2Github(callbackType)
	} else if serviceName == "Reddit" {
		authUrl = self.OAuth2Reddit(callbackType)
	} else if serviceName == "Asana" {
		authUrl = self.OAuth2Asana(callbackType)
	} else if serviceName == "Linkedin" {
		authUrl = self.OAuth2Linkedin(callbackType)
	} else {
		return authUrl, fmt.Errorf(UnknownServiceMessage)
	}
	return authUrl, nil
}

func (self *ServiceService) GetCallbackAndClientId(callbackType string, serviceName string, isIdNecessary bool) (string, string) {
	var callbackLink string
	var clientID string

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

func (self *ServiceService) OAuth2Google(callbackType string, appType string) string {
	var callbackLink string
	var appTypeLink string

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

func (self *ServiceService) OAuth2Spotify(callbackType string) string {
	callbackLink, _ := self.GetCallbackAndClientId(callbackType, "SPOTIFY", false)

	return ("https://accounts.spotify.com/authorize?" +
		clientIdParam + os.Getenv("SPOTIFY_CLIENT_ID") +
		codeResponseType +
		redirectUriParam + callbackLink +
		"&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state user-library-modify")
}

func (self *ServiceService) OAuth2Discord(callbackType string) string {
	callbackLink, _ := self.GetCallbackAndClientId(callbackType, "DISCORD", false)
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

func (self *ServiceService) OAuth2Github(callbackType string) string {
	callbackLink, clientID := self.GetCallbackAndClientId(callbackType, "GITHUB", true)

	return ("https://github.com/login/oauth/authorize?" +
		clientIdParam + clientID +
		redirectUriParam + callbackLink +
		"&scope=repo admin:org user" +
		stateParam + os.Getenv("GITHUB_STATE"))
}

func (self *ServiceService) OAuth2Reddit(callbackType string) string {
	callbackLink, clientID := self.GetCallbackAndClientId(callbackType, "REDDIT", true)

	return ("https://www.reddit.com/api/v1/authorize?" +
		clientIdParam + clientID +
		codeResponseType +
		stateParam + os.Getenv("REDDIT_STATE") +
		redirectUriParam + callbackLink +
		"&duration=permanent" +
		"&scope=identity read submit vote mysubreddits history account edit privatemessages")
}

func (self *ServiceService) OAuth2Asana(callbackType string) string {
	callbackLink, _ := self.GetCallbackAndClientId(callbackType, "ASANA", false)

	return ("https://app.asana.com/-/oauth_authorize?" +
		clientIdParam + os.Getenv("ASANA_CLIENT_ID") +
		redirectUriParam + callbackLink +
		codeResponseType +
		stateParam + os.Getenv("ASANA_STATE") +
		"&scope=default email profile")
}

func (self *ServiceService) OAuth2Linkedin(callbackType string) string {
	callbackLink, _ := self.GetCallbackAndClientId(callbackType, "LINKEDIN", false)

	return ("https://www.linkedin.com/oauth/v2/authorization?" +
		clientIdParam + url.QueryEscape(os.Getenv("LINKEDIN_CLIENT_ID")) +
		redirectUriParam + callbackLink +
		codeResponseType +
		"&scope=openid email profile w_member_social") +
		stateParam + os.Getenv("LINKEDIN_STATE")
}

func GetGoogleOAuth2AccessTokenWebRequestBody(code string, callback string) string {
	return fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s",
		code,
		os.Getenv("GOOGLE_WEB_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		callback,
		grantTypeAuthorization,
	)
}

func GetGoogleOAuth2AccessTokenMobileRequestBody(code string, callback string) string {
	return fmt.Sprintf(
		"code=%s&client_id=%s&redirect_uri=%s&grant_type=%s",
		code,
		os.Getenv("GOOGLE_MOBILE_CLIENT_ID"),
		callback,
		grantTypeAuthorization,
	)
}

func GetGoogleOAuth2AccessTokenRequest(code string, callbackType string, appType string) (*http.Request, error) {
	var googleCallback string
	var jsonBody string

	if callbackType == "service" {
		googleCallback = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		googleCallback = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	if appType == "web" {
		jsonBody = GetGoogleOAuth2AccessTokenWebRequestBody(code, googleCallback)
	} else if appType == "mobile" {
		jsonBody = GetGoogleOAuth2AccessTokenMobileRequestBody(code, googleCallback)
	} else {
		return nil, fmt.Errorf("Invalid app type")
	}

	tokenUrl := "https://oauth2.googleapis.com/token"
	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(ContentType, ContentTypeUrlEncoded)
	return request, nil
}

func GetGithubOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var githubCallback string
	var githubClientID string
	var githubClientSecret string

	if callbackType == "service" {
		githubCallback = os.Getenv("GITHUB_SERVICE_CALLBACK")
		githubClientID = os.Getenv("GITHUB_SERVICE_CLIENT_ID")
		githubClientSecret = os.Getenv("GITHUB_SERVICE_CLIENT_SECRET")
	} else if callbackType == "login" {
		githubCallback = os.Getenv("GITHUB_LOGIN_CALLBACK")
		githubClientID = os.Getenv("GITHUB_LOGIN_CLIENT_ID")
		githubClientSecret = os.Getenv("GITHUB_LOGIN_CLIENT_SECRET")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
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

	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func GetAsanaOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var asanaCallback string

	if callbackType == "service" {
		asanaCallback = os.Getenv("ASANA_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		asanaCallback = os.Getenv("ASANA_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	tokenUrl := "https://app.asana.com/-/oauth_token"
	jsonBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&grant_type=%s",
		os.Getenv("ASANA_CLIENT_ID"),
		os.Getenv("ASANA_CLIENT_SECRET"),
		code,
		asanaCallback,
		grantTypeAuthorization,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func GetLinkedinOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var callback string

	if callbackType == "service" {
		callback = os.Getenv("LINKEDIN_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		callback = os.Getenv("LINKEDIN_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	tokenUrl := "https://www.linkedin.com/oauth/v2/accessToken"
	jsonBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&grant_type=%s",
		url.QueryEscape(os.Getenv("LINKEDIN_CLIENT_ID")),
		url.QueryEscape(os.Getenv("LINKEDIN_CLIENT_SECRET")),
		code,
		callback,
		grantTypeAuthorization,
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
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://accounts.spotify.com/api/token", code, spotifyCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func GetDiscordOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var discordCallback string

	if callbackType == "service" {
		discordCallback = os.Getenv("DISCORD_SERVICE_CALLBACK")
	} else if callbackType == "login" {
		discordCallback = os.Getenv("DISCORD_LOGIN_CALLBACK")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://discord.com/api/v10/oauth2/token", code, discordCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func GetRedditOAuth2AccessTokenRequest(code string, callbackType string) (*http.Request, error) {
	var redditCallback string
	var redditClientID string
	var redditClientSecret string

	if callbackType == "service" {
		redditCallback = os.Getenv("REDDIT_SERVICE_CALLBACK")
		redditClientID = os.Getenv("REDDIT_SERVICE_CLIENT_ID")
		redditClientSecret = os.Getenv("REDDIT_SERVICE_CLIENT_SECRET")
	} else if callbackType == "login" {
		redditCallback = os.Getenv("REDDIT_LOGIN_CALLBACK")
		redditClientID = os.Getenv("REDDIT_LOGIN_CLIENT_ID")
		redditClientSecret = os.Getenv("REDDIT_LOGIN_CLIENT_SECRET")
	} else {
		return nil, fmt.Errorf(InvalidCallbackTypeMessage)
	}

	request, err := genericAccessTokenRequest("https://www.reddit.com/api/v1/access_token", code, redditCallback)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(redditClientID, redditClientSecret)
	request.Header.Set("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func genericAccessTokenRequest(tokenUrl string, code string, callbackUrl string) (*http.Request, error) {
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

func GetOAuth2UserEmailRequest(userInfoUrl string, accessToken string) (*http.Request, error) {
	request, errRequest := http.NewRequest("GET", userInfoUrl, nil)
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
		"client_id=%s&client_secret=%s&grant_type=%s&refresh_token=%s",
		os.Getenv("GOOGLE_WEB_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		grantTypeRefreshToken,
		refreshToken,
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return request, errRequest
	}

	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetAsanaRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://app.asana.com/-/oauth_token"
	jsonBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=%s&refresh_token=%s",
		os.Getenv("ASANA_CLIENT_ID"),
		os.Getenv("ASANA_CLIENT_SECRET"),
		grantTypeRefreshToken,
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
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://discord.com/api/v10/oauth2/token"
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)

	return request, nil
}

func (self *ServiceService) GetRedditRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	tokenUrl := "https://www.reddit.com/api/v1/access_token"
	request, err := self.genericRefreshTokenRequest(tokenUrl, refreshToken)
	if err != nil {
		return request, err
	}

	request.SetBasicAuth(os.Getenv("REDDIT_SERVICE_CLIENT_ID"), os.Getenv("REDDIT_SERVICE_CLIENT_SECRET"))
	request.Header.Set(ContentType, ContentTypeUrlEncoded)
	request.Header.Set("User-Agent", os.Getenv("REDDIT_SERVICE_USER_AGENT"))

	return request, nil
}

func (self *ServiceService) genericRefreshTokenRequest(tokenUrl string, refreshToken string) (*http.Request, error) {
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

func (self *ServiceService) ExecuteRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(APICallFailedMessage)
	}
	return res, nil
}

func (self *ServiceService) ExecuteApiRequest(url string, method string, typeToken string, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", typeToken+accessToken)
	req.Header.Set("Content-Type", "application/json")

	return self.ExecuteRequest(req)
}

func (self *ServiceService) GetResultTokenFromCode(code string, serviceName string, callbackType string, appType string) (entities.ResultToken, error) {
	var tokenRes entities.ResultToken
	var request *http.Request
	var errRequest error

	if serviceName == "Google" {
		request, errRequest = GetGoogleOAuth2AccessTokenRequest(code, callbackType, appType)
	} else if serviceName == "Spotify" {
		request, errRequest = GetSpotifyOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Discord" {
		request, errRequest = GetDiscordOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Github" {
		request, errRequest = GetGithubOAuth2AccessTokenRequest(code, callbackType)
		request.Header.Set("Accept", "application/json")
		tokenRes.ExpiresIn = oneYearSecond
	} else if serviceName == "Reddit" {
		codeStripped := strings.TrimSuffix(code, "#_")
		request, errRequest = GetRedditOAuth2AccessTokenRequest(codeStripped, callbackType)
	} else if serviceName == "Asana" {
		request, errRequest = GetAsanaOAuth2AccessTokenRequest(code, callbackType)
	} else if serviceName == "Linkedin" {
		request, errRequest = GetLinkedinOAuth2AccessTokenRequest(code, callbackType)
	} else {
		return tokenRes, fmt.Errorf(UnknownServiceMessage)
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

func decodeUserInfoWithGithub(res *http.Response) (entities.UserInfo, error) {
	var userInfos []entities.UserInfo

	err := json.NewDecoder(res.Body).Decode(&userInfos)
	if err != nil {
		return userInfos[0], err
	}
	return userInfos[0], nil
}

func (self *ServiceService) GetUserInfoFromService(accessToken string, serviceName string) (entities.UserInfo, error) {
	var userInfo entities.UserInfo
	var request *http.Request
	var err error

	if serviceName == "Google" {
		request, err = GetOAuth2UserEmailRequest("https://www.googleapis.com/oauth2/v3/userinfo", accessToken)
	} else if serviceName == "Spotify" {
		request, err = GetOAuth2UserEmailRequest("https://api.spotify.com/v1/me", accessToken)
	} else if serviceName == "Discord" {
		request, err = GetOAuth2UserEmailRequest("https://discord.com/api/users/@me", accessToken)
	} else if serviceName == "Github" {
		request, err = GetOAuth2UserEmailRequest(githubBaseUrl+"user/emails", accessToken)
	} else {
		return userInfo, fmt.Errorf(UnknownServiceMessage)
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
