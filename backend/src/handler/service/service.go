package service_handler

import (
	"backend/src/entities"
	_ "backend/src/handler/service/docs"
	"backend/src/service"
	"encoding/json"
	"fmt"
	"net/http"

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
	self.PrivateRoutes(router)
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

func (self *ServiceHandler) executeSpotifyReactions(reactionTitle string, userEmail string, workflow entities.Workflow) error {

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(userEmail, "Spotify")
	if errRefreshing != nil {
		return errRefreshing
	}

	if reactionTitle == "Start playback" {
		return self.spotifyRequest("https://api.spotify.com/v1/me/player/play", "PUT", accessToken)
	}
	if reactionTitle == "Pause playback" {
		return self.spotifyRequest("https://api.spotify.com/v1/me/player/pause", "PUT", accessToken)
	}
	if reactionTitle == "Skip track" {
		return self.spotifyRequest("https://api.spotify.com/v1/me/player/next", "POST", accessToken)
	}
	if reactionTitle == "Add track to playback queue" {
		uri, uriErr := workflow.ReactionParam["uri"]
		if !uriErr {
			return fmt.Errorf("Incorrect URI")
		}
		urlQueue := "https://api.spotify.com/v1/me/player/queue?uri=" + uri.(string)
		return self.spotifyRequest(urlQueue, "POST", accessToken)
	}
	if reactionTitle == "Save a track" {
		id, idErr := workflow.ReactionParam["id"]
		if !idErr {
			return fmt.Errorf("Incorrect ID")
		}
		urlSave := "https://api.spotify.com/v1/me/tracks?ids=" + id.(string)
		return self.spotifyRequest(urlSave, "PUT", accessToken)
	}
	return fmt.Errorf("Unknown reaction")
}

func (self *ServiceHandler) executeActiveActionsEveryDay(timeRes entities.TimeResponse, AllworkflowFound []entities.Workflow) error {
	for _, workflow := range AllworkflowFound {
		hour, hourExists := workflow.ActionParam["hour"]
		minute, minuteExists := workflow.ActionParam["minute"]
		if !hourExists || !minuteExists {
			continue
		}

		if int(hour.(float64)) == timeRes.Hour && int(minute.(float64)) == timeRes.Minute && workflow.IsActivated {
			reactionFound, errReaction := self.ServiceService.FindReactionByID(workflow.ReactionId)
			foundUser, errUser := self.UserService.FindUserById(workflow.OwnerId)
			if errReaction != nil || errUser != nil {
				continue
			}

			self.executeSpotifyReactions(reactionFound.Name, foundUser.Email, workflow)
		}
	}
	return nil
}

func (self *ServiceHandler) checkNameOfAction(timeRes entities.TimeResponse, action entities.Action) error {
	if action.Name == "Every day at" {

		AllworkflowFound, errWorkflow := self.ServiceService.FindWorkflowByActionId(action.Id)
		if errWorkflow != nil {
			return errWorkflow
		}
		return self.executeActiveActionsEveryDay(timeRes, AllworkflowFound)
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
