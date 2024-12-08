package workflow_handler

import (
	"backend/src/handler/middleware"
	"backend/src/service"

	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	WorkflowService service.WorkflowService
	UserService     service.UserService
}

const InvalidRequestBodyMessage = "Invalid request body"

func NewWorkflowHandler(WorkflowService service.WorkflowService, UserService service.UserService, router *gin.Engine) *WorkflowHandler {
	handler := &WorkflowHandler{
		WorkflowService: WorkflowService,
		UserService:     UserService,
	}
	handler.SetupRoutes(router)
	return handler
}

func (self *WorkflowHandler) SetupRoutes(router *gin.Engine) {
	self.PrivateRoutes(router)
}

func (self *WorkflowHandler) PrivateRoutes(router *gin.Engine) {
	workflow := router.Group("/workflow", middleware.VerifyJWTCookie)
	{
		workflow.POST("/create", self.CreateWorkflow)
	}

}
