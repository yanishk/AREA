package workflow_handler

import (
	"backend/src/entities"
	"backend/src/handler/middleware"
	_ "backend/src/handler/workflow/docs"
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
		workflow.PUT("/modify-activation", self.ActivateDesactivateWorkflow)
	}

}

// @Summary		Create Workflow
// @Description	Create a Workflow
// @Tags			Workflow
// @Produce		json
// @Param			workflow	body		entities.NewWorkflow	true	"Workflow informations"
// @Success		200		{object}	docs_workflow.WorkflowCreateWorkflowSuccessResponse
// @Failure		400		{object}	docs_workflow.WorkflowCreateWorkflowBadRequestResponse
// @Failure		401		{object}	docs_workflow.WorkflowCreateWorkflowUnauthorizedResponse
// @Failure		500		{object}	docs_workflow.WorkflowCreateWorkflowInternalServerErrorResponse
// @Router			/workflow/create [post]
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
			"error": "Could not create workflow",
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Successful workflow creation",
	})
}

// @Summary		Modify Workflow Activation Status
// @Description	Modify the activation status of the specified workflow
// @Tags			Workflow
// @Produce		json
// @Param			workflow	body		entities.ChangeActivationWorkflow	true	"Workflow informations"
// @Success		200		{object}	docs_workflow.WorkflowModifyWorkflowActivationStatusSuccessResponse
// @Failure		400		{object}	docs_workflow.WorkflowModifyWorkflowActivationStatusBadRequestResponse
// @Failure		500		{object}	docs_workflow.WorkflowModifyWorkflowActivationStatusInternalServerErrorResponse
// @Router			/workflow/modify-activation [put]
func (self *WorkflowHandler) ActivateDesactivateWorkflow(context *gin.Context) {

	var workflowActivation entities.ChangeActivationWorkflow

	err := context.ShouldBindJSON(&workflowActivation)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestBodyMessage,
		})
		return
	}

	foundWorkflow, errFetchWorkflow := self.WorkflowService.FindWorkflowById(workflowActivation.WorkflowId)
	if errFetchWorkflow != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find specified workflow",
		})
		return
	}

	errUpdate := self.WorkflowService.ActivateDesactivateWorkflow(workflowActivation.WorkflowId, foundWorkflow.IsActivated)
	if errUpdate != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not update workflow's activation status",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Successfully updates workflow's activation status",
	})
}
