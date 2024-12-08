package workflow_handler

import (
	"backend/src/entities"
	"backend/src/handler/middleware"
	"backend/src/service"
	"net/http"

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

func (self *WorkflowHandler) CreateWorkflow(context *gin.Context) {

	var newWorkflow entities.NewWorkflow

	err := context.ShouldBindJSON(&newWorkflow)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestBodyMessage,
		})
		return
	}

	email, errGetEmail := context.Get("email")
	if !errGetEmail {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not a valid string",
		})
		return
	}
	errCreationWorkflow := self.WorkflowService.CreateWorkflow(emailStr, newWorkflow)
	if errCreationWorkflow != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "could not create workflow",
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "successful workflow creation",
	})
}
