package service_handler

import (
	_ "backend/src/handler/service/docs"
	"backend/src/service"

	"github.com/gin-gonic/gin"
	"net/http"
)

type ServiceHandler struct {
	ServiceService service.ServiceService
}

func NewServiceHandler(ServiceService service.ServiceService, router *gin.Engine) *ServiceHandler {
	handler := &ServiceHandler{ServiceService: ServiceService}
	handler.SetupRoutes(router)
	return handler
}

func (self *ServiceHandler) SetupRoutes(router *gin.Engine) {
	self.PrivateRoutes(router)
}

func (self *ServiceHandler) PrivateRoutes(router *gin.Engine) {
	router.GET("/actionsServices", self.RetrieveActionsServices)
	router.GET("/reactionsServices", self.RetrieveReactionsServices)
}

// @Summary		Actions Services
// @Description	Get Retrieve Actions Services
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceActionsReactionsSuccessResponse
// @Success		500		{object}	docs_service.ServiceActionReactionInternalServerErrorResponse
// @Router			/actionsServices [get]
func (self *ServiceHandler) RetrieveActionsServices(context *gin.Context) {
	actionsServices, err := self.ServiceService.RetrieveActionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	} else {
		context.IndentedJSON(http.StatusOK, actionsServices)
	}
}

// @Summary		Reactions Services
// @Description	Get Retrieve Reactions Services
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceActionsReactionsSuccessResponse
// @Success		500		{object}	docs_service.ServiceActionReactionInternalServerErrorResponse
// @Router			/reactionsServices [get]
func (self *ServiceHandler) RetrieveReactionsServices(context *gin.Context) {
	reactionsServices, err := self.ServiceService.RetrieveReactionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	} else {
		context.IndentedJSON(http.StatusOK, reactionsServices)
	}
}
