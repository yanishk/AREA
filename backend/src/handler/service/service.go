package service_handler

import (
	"backend/src/entities"
	_ "backend/src/handler/service/docs"
	"backend/src/service"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceService     service.ServiceService
	UserServiceService service.UserServiceService
	UserService        service.UserService
}

const callbackTypeParamName = "callback-type"
const appTypeParamName = "app-type"

const botBearer = "Bot "
const accessBearer = "Bearer "

const InternalServerErrorMessage = "Internal server error"
const InvalidCallbackTypeMessage = "Invalid callback type"
const InvalidAppTypeMessage = "Invalid app type"
const APICallFailedMessage = "API call failed"
const getUserServersFailedMessage = "Could not retrieve the user's servers"
const getServerChannelsFailedMessage = "Could not retrieve the server's channels"

const authURL = "auth-url"

func NewServiceHandler(ServiceService service.ServiceService, UserServiceService service.UserServiceService, UserService service.UserService, router *gin.Engine) *ServiceHandler {
	handler := &ServiceHandler{
		ServiceService:     ServiceService,
		UserServiceService: UserServiceService,
		UserService:        UserService,
	}
	handler.SetupRoutes(router)
	return handler
}

func (self *ServiceHandler) SetupRoutes(router *gin.Engine) {
	self.PublicRoutes(router)
	self.PrivateRoutes(router)
}

func (self *ServiceHandler) PublicRoutes(router *gin.Engine) {
	authentication := router.Group("/authentication")
	{
		authentication.GET("/Spotify/:callback-type/:app-type", self.OAuth2Spotify)
		authentication.GET("/Google/:callback-type/:app-type", self.OAuth2Google)
		authentication.GET("/Discord/:callback-type/:app-type", self.OAuth2Discord)
		authentication.GET("/Github/:callback-type/:app-type", self.OAuth2Github)
	}
}

func (self *ServiceHandler) PrivateRoutes(router *gin.Engine) {
	service := router.Group("/:service")
	{
		service.GET("/actions", self.RetrieveActionsFromService)
		service.GET("/reactions", self.RetrieveReactionsFromService)
	}
	router.GET("/actions-services", self.RetrieveActionsServices)
	router.GET("/reactions-services", self.RetrieveReactionsServices)
	router.GET("/Discord/user/servers", self.getDiscordUserServers)
	router.GET("/Discord/server/channels", self.getDiscordGuildChannels)
}

func (self *ServiceHandler) apiRequest(url string, method string, typeToken string, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", typeToken+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 200 {
		return resp, fmt.Errorf(APICallFailedMessage)
	}
	return resp, nil
}

func (self *ServiceHandler) executeSpotifyReactions(workflow entities.Workflow) error {

	reactionFound, errReaction := self.ServiceService.FindReactionByID(workflow.ReactionId)
	foundUser, errUser := self.UserService.FindUserById(workflow.OwnerId)
	if errReaction != nil || errUser != nil {
		return fmt.Errorf("Error finding reaction")
	}
	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(foundUser.Email, foundUser.ConnectionType, "Spotify")
	if errRefreshing != nil {
		return errRefreshing
	}

	var resp *http.Response
	if reactionFound.Name == "Start playback" {
		resp, _ = self.apiRequest("https://api.spotify.com/v1/me/player/play", "PUT", accessBearer, accessToken, nil)
	}
	if reactionFound.Name == "Pause playback" {
		resp, _ = self.apiRequest("https://api.spotify.com/v1/me/player/pause", "PUT", accessBearer, accessToken, nil)
	}
	if reactionFound.Name == "Skip track" {
		resp, _ = self.apiRequest("https://api.spotify.com/v1/me/player/next", "POST", accessBearer, accessToken, nil)
	}
	if reactionFound.Name == "Add track to playback queue" {
		uri, uriErr := workflow.ReactionParam["uri"]
		if !uriErr {
			return fmt.Errorf("Incorrect URI")
		}
		urlQueue := "https://api.spotify.com/v1/me/player/queue?uri=" + uri.(string)
		resp, _ = self.apiRequest(urlQueue, "POST", accessBearer, accessToken, nil)
	}
	if reactionFound.Name == "Save a track" {
		id, idErr := workflow.ReactionParam["id"]
		if !idErr {
			return fmt.Errorf("Incorrect ID")
		}
		urlSave := "https://api.spotify.com/v1/me/tracks?ids=" + id.(string)
		resp, _ = self.apiRequest(urlSave, "PUT", accessBearer, accessToken, nil)
	}
	defer resp.Body.Close()
	return fmt.Errorf("Unknown reaction")
}

// @Summary		Discord User Servers
// @Description	Retrieve the discord servers from the authenticated user
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceGetDiscordUserServersSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetDiscordUserServersInternalServerErrorResponse
// @Router			/Discord/user/servers [get]
func (self *ServiceHandler) getDiscordUserServers(context *gin.Context) {
	emailUser, err := context.Get("email")
	if !err {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the user's mail",
		})
		return
	}

	connectionTypeUser, err := context.Get("connectionType")
	if !err {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the user's connection type",
		})
		return
	}

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(emailUser.(string), connectionTypeUser.(string), "Discord")
	if errRefreshing != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve or refresh the token",
		})
		return
	}

	resp, errRequest := self.apiRequest("https://discord.com/api/v10/users/@me/guilds", "GET", accessBearer, accessToken, nil)
	if errRequest != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": getUserServersFailedMessage,
		})
		return
	}
	defer resp.Body.Close()

	var guilds []map[string]interface{}
	var ownedServers []map[string]interface{}
	errDecoder := json.NewDecoder(resp.Body).Decode(&guilds)
	if errDecoder != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": getUserServersFailedMessage,
		})
		return
	}

	for _, guild := range guilds {
		if guild["owner"] == true {
			ownedServers = append(ownedServers, guild)
		}
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"servers": ownedServers,
	})
}

// @Summary		Discord Server Channels
// @Description	Retrieve the channels from the selected server
// @Tags			Services
// @Produce		json
// @Param        id  query      string  true  "server id"
// @Success		200		{object}	docs_service.ServiceGetDiscordGuildChannelsSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetDiscordGuildChannelsInternalServerErrorResponse
// @Router /Discord/server/channels [get]
func (self *ServiceHandler) getDiscordGuildChannels(context *gin.Context) {
	tokenBot := os.Getenv("DISCORD_BOT_TOKEN")
	guildId := context.Query("id")
	getChannelUrl := "https://discord.com/api/v10/guilds/" + guildId + "/channels"

	resp, errRequest := self.apiRequest(getChannelUrl, "GET", botBearer, tokenBot, nil)
	if errRequest != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": getServerChannelsFailedMessage,
		})
		return
	}
	defer resp.Body.Close()

	var channels []map[string]interface{}
	errDecoder := json.NewDecoder(resp.Body).Decode(&channels)
	if errDecoder != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": getServerChannelsFailedMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"channels": channels,
	})
}

func (self *ServiceHandler) executeDiscordReactions(workflow entities.Workflow) error {
	tokenBot := os.Getenv("DISCORD_BOT_TOKEN")
	reactionFound, errReaction := self.ServiceService.FindReactionByID(workflow.ReactionId)
	foundUser, errUser := self.UserService.FindUserById(workflow.OwnerId)
	if errReaction != nil || errUser != nil {
		return fmt.Errorf("Error finding reaction")
	}
	_, errRefreshing := self.UserServiceService.CallApiAndRefresh(foundUser.Email, foundUser.ConnectionType, "Discord")
	if errRefreshing != nil {
		return errRefreshing
	}

	if reactionFound.Name == "Post a message to a channel" {
		message, messageErr := workflow.ReactionParam["message"]
		channelID, channelErr := workflow.ReactionParam["channelid"]
		if !messageErr || !channelErr {
			return fmt.Errorf("Failed to get message body or channel id")
		}

		urlDiscordAuth := "https://discord.com/api/v10/channels/" + channelID.(string) + "/messages"
		content, errContent := json.Marshal(map[string]string{
			"content": message.(string),
		})
		if errContent != nil {
			return fmt.Errorf("Failed to marshall message content")
		}

		resp, errRequest := self.apiRequest(urlDiscordAuth, "POST", botBearer, tokenBot, bytes.NewBuffer(content))
		if errRequest != nil {
			return errRequest
		}
		defer resp.Body.Close()
	}
	return fmt.Errorf("Unknown reaction")
}

func (self *ServiceHandler) executeActiveActionsEveryDay(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !hourExists || !minuteExists {
		return fmt.Errorf("Error in the hour and minute sent")
	}

	if int(hour.(float64)) == timeRes.Hour && int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
	}
	return nil
}

func (self *ServiceHandler) executeActiveActionsEveryHour(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists {
		return fmt.Errorf("Error in the minute sent")
	}

	if int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
	}
	return nil
}

func (self *ServiceHandler) executeActiveActionsEveryWeekAt(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	days, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists {
		return fmt.Errorf("Error in the minute, hour or day sent")
	}

	dayArray, err := days.([]interface{})
	if !err {
		return fmt.Errorf("Error parsing the day array")
	}

	for _, day := range dayArray {
		if int(day.(float64)) == timeRes.WeekDay && int(hour.(float64)) == timeRes.Hour &&
			int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
			self.executeSpotifyReactions(workflow)
			self.executeDiscordReactions(workflow)
		}
	}

	return nil
}

func (self *ServiceHandler) executeActiveActionsEveryMonthOnThe(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	day, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists {
		return fmt.Errorf("Error in the minute, hour or day sent")
	}

	if int(day.(float64)) == timeRes.Day && int(hour.(float64)) == timeRes.Hour &&
		int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
	}
	return nil
}

func (self *ServiceHandler) executeActiveActionsEveryYearOn(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	month, monthExists := workflow.ActionParam["month"]
	day, dayExists := workflow.ActionParam["day"]
	hour, hourExists := workflow.ActionParam["hour"]
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists || !hourExists || !dayExists || !monthExists {
		return fmt.Errorf("Error in the minute, hour, day or month sent")
	}

	if int(month.(float64)) == timeRes.Month && int(day.(float64)) == timeRes.Day && int(hour.(float64)) == timeRes.Hour &&
		int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		self.executeSpotifyReactions(workflow)
		self.executeDiscordReactions(workflow)
	}
	return nil
}

func (self *ServiceHandler) checkNameOfAction(timeRes entities.TimeResponse, action entities.Action) error {
	AllworkflowFound, errWorkflow := self.ServiceService.FindWorkflowByActionId(action.Id)
	if errWorkflow != nil {
		return errWorkflow
	}

	for _, workflow := range AllworkflowFound {
		if action.Name == "Every day at" {
			self.executeActiveActionsEveryDay(timeRes, workflow)
		}
		if action.Name == "Every hour at" {
			self.executeActiveActionsEveryHour(timeRes, workflow)
		}
		if action.Name == "Every day of the week at" {
			self.executeActiveActionsEveryWeekAt(timeRes, workflow)
		}
		if action.Name == "Every month on the" {
			self.executeActiveActionsEveryMonthOnThe(timeRes, workflow)
		}
		if action.Name == "Every year on" {
			self.executeActiveActionsEveryYearOn(timeRes, workflow)
		}
	}
	return nil
}

func (self *ServiceHandler) requestToTimeAPI() (entities.TimeResponse, error) {
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

func (self *ServiceHandler) CallTimeApi() error {
	var timeRes entities.TimeResponse

	timeRes, errRequest := self.requestToTimeAPI()
	if errRequest != nil {
		return errRequest
	}

	serviceFound, errFindingService := self.ServiceService.FindServiceByName("Time & Date")
	if errFindingService != nil {
		return errFindingService
	}

	allactionFound, errFindingaction := self.ServiceService.FindActionByServiceId(serviceFound.Id)
	if errFindingaction != nil {
		return errFindingaction
	}

	for _, action := range allactionFound {
		errCheckName := self.checkNameOfAction(timeRes, action)
		if errCheckName != nil {
			return errCheckName
		}
	}
	return nil
}

// @Summary		Google OAuth2
// @Description	Get Google OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Param        app-type  path      string  true  "App type (web or mobile)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Google/{callback-type}/{app-type} [get]
func (self *ServiceHandler) OAuth2Google(context *gin.Context) {
	var callbackLink string
	var appTypeLink string
	callbackType := context.Param(callbackTypeParamName)
	appType := context.Param(appTypeParamName)

	if callbackType == "login" {
		callbackLink = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		callbackLink = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	if appType == "web" {
		appTypeLink = os.Getenv("GOOGLE_WEB_CLIENT_ID")
	} else if appType == "mobile" {
		appTypeLink = os.Getenv("GOOGLE_MOBILE_CLIENT_ID")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	googleAuthUrl := "https://accounts.google.com/o/oauth2/auth?client_id=" +
		appTypeLink + "&response_type=code&access_type=offline&redirect_uri=" +
		callbackLink + "&scope=email"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: googleAuthUrl,
	})
}

// @Summary		Spotify OAuth2
// @Description	Get Spotify OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Param        app-type  path      string  true  "App type (web or mobile)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Spotify/{callback-type}/{app-type} [get]
func (self *ServiceHandler) OAuth2Spotify(context *gin.Context) {
	var callbackLink string
	callbackType := context.Param(callbackTypeParamName)
	appType := context.Param(appTypeParamName)

	if callbackType == "login" {
		callbackLink = os.Getenv("SPOTIFY_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		callbackLink = os.Getenv("SPOTIFY_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	if appType != "web" && appType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_CLIENT_ID") + "&response_type=code&redirect_uri=" +
		callbackLink + "&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state user-library-modify"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: spotifyAuthUrl,
	})
}

// @Summary		Discord OAuth2
// @Description	Get Discord OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Param        app-type  path      string  true  "App type (web or mobile)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Discord/{callback-type}/{app-type} [get]
func (self *ServiceHandler) OAuth2Discord(context *gin.Context) {
	var callbackLink string
	callbackType := context.Param(callbackTypeParamName)
	appType := context.Param(appTypeParamName)

	if callbackType == "login" {
		callbackLink = os.Getenv("DISCORD_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		callbackLink = os.Getenv("DISCORD_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	if appType != "web" && appType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}
	discordAuthUrl := "https://discord.com/oauth2/authorize?client_id=" +
		os.Getenv("DISCORD_CLIENT_ID") +
		"&permissions=141376&response_type=code&redirect_uri=" +
		callbackLink +
		"&scope=identify+email+guilds+bot"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: discordAuthUrl,
	})
}

// @Summary		Github OAuth2
// @Description	Get Github OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Param        app-type  path      string  true  "App type (web or mobile)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Github/{callback-type}/{app-type} [get]
func (self *ServiceHandler) OAuth2Github(context *gin.Context) {
	var callbackLink string
	var clientID string
	callbackType := context.Param(callbackTypeParamName)
	appType := context.Param(appTypeParamName)

	if callbackType == "login" {
		clientID = os.Getenv("GITHUB_LOGIN_CLIENT_ID")
		callbackLink = os.Getenv("GITHUB_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		clientID = os.Getenv("GITHUB_SERVICE_CLIENT_ID")
		callbackLink = os.Getenv("GITHUB_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	if appType != "web" && appType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	githubAuthUrl := "https://github.com/login/oauth/authorize?" +
		"client_id=" + clientID +
		"&redirect_uri=" + callbackLink +
		"&scope=repo admin:repo_hook write:repo_hook read:repo_hook read:user user:email user:follow read:org" +
		"&state=" + os.Getenv("GITHUB_STATE")

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: githubAuthUrl,
	})
}

// @Summary		Actions Services
// @Description	Retrieve services with actions available
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceRetrieveActionsReactionsServicesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/actions-services [get]
func (self *ServiceHandler) RetrieveActionsServices(context *gin.Context) {
	actionsServices, err := self.ServiceService.RetrieveActionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, actionsServices)
}

// @Summary		Reactions Services
// @Description	Retrieve services with reactions available
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceRetrieveActionsReactionsServicesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/reactions-services [get]
func (self *ServiceHandler) RetrieveReactionsServices(context *gin.Context) {
	reactionsServices, err := self.ServiceService.RetrieveReactionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, reactionsServices)
}

// @Summary		Actions From Service
// @Description	Retrieve actions from requested service
// @Tags			Services
// @Produce		json
// @Param        service  path      string  true  "Service name (Google, Spotify, Time & Date..)"
// @Success		200		{object}	docs_service.ServiceRetrieveActionsFromServiceSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/{service}/actions [get]
func (self *ServiceHandler) RetrieveActionsFromService(context *gin.Context) {
	serviceName := context.Param("service")

	actions, errActions := self.ServiceService.RetrieveActionsFromService(serviceName)
	if errActions != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"actions": actions,
	})
}

// @Summary		Reactions From Service
// @Description	Retrieve reactions from requested service
// @Tags			Services
// @Produce		json
// @Param        service  path      string  true  "Service name (Google, Spotify, Time & Date..)"
// @Success		200		{object}	docs_service.ServiceRetrieveReactionsFromServiceSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/{service}/reactions [get]
func (self *ServiceHandler) RetrieveReactionsFromService(context *gin.Context) {
	serviceName := context.Param("service")

	reactions, errReactions := self.ServiceService.RetrieveReactionsFromService(serviceName)
	if errReactions != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"reactions": reactions,
	})
}
