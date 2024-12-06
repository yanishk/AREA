package service_handler

import (
	"backend/src/service"

	"github.com/gin-gonic/gin"
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

// à documenter
func (self *ServiceHandler) RetrieveActionsServices(context *gin.Context) {
	//
}

func (self *ServiceHandler) RetrieveReactionsServices(context *gin.Context) {
	//
}
