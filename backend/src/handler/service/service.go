package service_handler

import (
	"backend/src/entities"
	_ "backend/src/handler/service/docs"
	"backend/src/service"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceService     service.ServiceService
	UserServiceService service.UserServiceService
	UserService        service.UserService
}

const InternalServerErrorMessage = "Internal server error"
const SpotifyAPICallFailed = "Spotify API call failed"

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
		authentication.GET("/Spotify/:callback-type", self.OAuth2Spotify)
		authentication.GET("/Google/:callback-type", self.OAuth2Google)
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
}

func (self *ServiceHandler) spotifyRequest(url string, method string, accessToken string) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(SpotifyAPICallFailed)
	}
	return nil
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

	if reactionFound.Name == "Start playback" {
		return self.spotifyRequest("https://api.spotify.com/v1/me/player/play", "PUT", accessToken)
	}
	if reactionFound.Name == "Pause playback" {
		return self.spotifyRequest("https://api.spotify.com/v1/me/player/pause", "PUT", accessToken)
	}
	if reactionFound.Name == "Skip track" {
		return self.spotifyRequest("https://api.spotify.com/v1/me/player/next", "POST", accessToken)
	}
	if reactionFound.Name == "Add track to playback queue" {
		uri, uriErr := workflow.ReactionParam["uri"]
		if !uriErr {
			return fmt.Errorf("Incorrect URI")
		}
		urlQueue := "https://api.spotify.com/v1/me/player/queue?uri=" + uri.(string)
		return self.spotifyRequest(urlQueue, "POST", accessToken)
	}
	if reactionFound.Name == "Save a track" {
		id, idErr := workflow.ReactionParam["id"]
		if !idErr {
			return fmt.Errorf("Incorrect ID")
		}
		urlSave := "https://api.spotify.com/v1/me/tracks?ids=" + id.(string)
		return self.spotifyRequest(urlSave, "PUT", accessToken)
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
		return self.executeSpotifyReactions(workflow)
	}
	return nil
}

func (self *ServiceHandler) executeActiveActionsEveryHour(timeRes entities.TimeResponse, workflow entities.Workflow) error {
	minute, minuteExists := workflow.ActionParam["minute"]
	if !minuteExists {
		return fmt.Errorf("Error in the minute sent")
	}

	if int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
		return self.executeSpotifyReactions(workflow)
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
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Google/{callback-type} [get]
func (self *ServiceHandler) OAuth2Google(context *gin.Context) {
	var callbackLink string
	callbackType := context.Param("callback-type")

	if (callbackType == "login") {
		callbackLink = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else if (callbackType == "service") {
		callbackLink = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid callback type",
		})
		return
	}

	googleAuthUrl := "https://accounts.google.com/o/oauth2/auth?client_id=" +
		os.Getenv("GOOGLE_CLIENT_ID") + "&response_type=code&access_type=offline&redirect_uri=" +
		callbackLink + "&scope=email"

	context.IndentedJSON(http.StatusOK, gin.H{
		"auth-url": googleAuthUrl,
	})
}

// @Summary		Spotify OAuth2
// @Description	Get Spotify OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Spotify/{callback-type} [get]
func (self *ServiceHandler) OAuth2Spotify(context *gin.Context) {
	var callbackLink string
	callbackType := context.Param("callback-type")

	if (callbackType == "login") {
		callbackLink = os.Getenv("SPOTIFY_LOGIN_CALLBACK")
	} else if (callbackType == "service") {
		callbackLink = os.Getenv("SPOTIFY_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid callback type",
		})
		return
	}

	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_CLIENT_ID") + "&response_type=code&redirect_uri=" +
		callbackLink + "&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state user-library-modify"

	context.IndentedJSON(http.StatusOK, gin.H{
		"auth-url": spotifyAuthUrl,
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
