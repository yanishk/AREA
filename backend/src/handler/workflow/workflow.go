package workflow_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/src/entities"
	"backend/src/handler/middleware"
	_ "backend/src/handler/workflow/docs"
	"backend/src/service"
)

type WorkflowHandler struct {
	WorkflowService service.WorkflowService
	UserService     service.UserService
}

const invalidRequestBodyMessage = "Invalid request body"

func NewWorkflowHandler(WorkflowService service.WorkflowService, UserService service.UserService, router *gin.Engine) *WorkflowHandler {
	handler := &WorkflowHandler{
		WorkflowService: WorkflowService,
		UserService:     UserService,
	}
	handler.setupRoutes(router)
	return handler
}

func (self *WorkflowHandler) setupRoutes(router *gin.Engine) {
	self.publicRoutes(router)
	self.privateRoutes(router)
}

func (self *WorkflowHandler) publicRoutes(router *gin.Engine) {
	webhook := router.Group("/webhooks")
	{
		webhook.POST("/:service", self.receiveServiceWebhook)
	}
}

func (self *WorkflowHandler) privateRoutes(router *gin.Engine) {
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	workflow := private.Group("/workflows")
	{
		workflow.POST("", self.createWorkflow)
		workflow.GET("", self.getUserWorkflows)
		workflow.PUT("/:id", self.updateWorkflow)
		workflow.DELETE("/:id", self.deleteWorkflow)
	}
}

// // @Summary		Receive Service Webhook
// // @Description	Receive webhook from a service
// // @Tags			Webhooks
// // @Produce		json
// @Success		200		{object}	docs_workflow.WorkflowReceiveServiceWebhookSuccessResponse
// @Success		400		{object}	docs_workflow.WorkflowReceiveServiceWebhookBadRequestResponse
// // @Router		/webhooks/{service} [post]
func (self *WorkflowHandler) receiveServiceWebhook(context *gin.Context) {
	serviceName := context.Param("service")

	err := self.WorkflowService.CheckWebhooksWorkflows(serviceName, context.Request)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": invalidRequestBodyMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Webhook received",
	})
}

// @Summary		Create Workflow
// @Description	Create a workflow
// @Tags			Workflows
// @Produce		json
// @Param			workflow	body		entities.NewWorkflow	true	"Workflow informations"
// @Success		200		{object}	docs_workflow.WorkflowCreateWorkflowSuccessResponse
// @Failure		400		{object}	docs_workflow.WorkflowCreateWorkflowBadRequestResponse
// @Failure		401		{object}	docs_workflow.WorkflowCreateWorkflowUnauthorizedResponse
// @Failure		500		{object}	docs_workflow.WorkflowCreateWorkflowInternalServerErrorResponse
// @Router			/workflows [post]
func (self *WorkflowHandler) createWorkflow(context *gin.Context) {
	var newWorkflow entities.NewWorkflow
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	err := context.ShouldBindJSON(&newWorkflow)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": invalidRequestBodyMessage,
		})
		return
	}

	errCreationWorkflow := self.WorkflowService.CreateWorkflow(email, connectionType, newWorkflow)
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

// @Summary		Retrieve User's Workflows
// @Description	Retrieve the workflows of the user actually connected
// @Tags			Workflows
// @Produce		json
// @Success		200		{object}	docs_workflow.WorkflowRetrieveUserWorkflowsSuccessResponse
// @Failure		401		{object}	docs_workflow.WorkflowRetrieveUserWorkflowsUnauthorizedResponse
// @Failure		500		{object}	docs_workflow.WorkflowRetrieveUserWorkflowsInternalServerErrorResponse
// @Router			/workflows [get]
func (self *WorkflowHandler) getUserWorkflows(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")
	retrievedWorkflow, errRetrievedWorkflow := self.WorkflowService.GetUserWorkflows(email, connectionType)

	if errRetrievedWorkflow != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve user's workflows",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"workflows": retrievedWorkflow,
	})
}

// @Summary		Update Workflow
// @Description	Update a user's workflow by specifying the workflow id
// @Tags			Workflows
// @Produce		json
// @Param        id     path     string  true  "Workflow id"
// @Param			workflow	body		entities.UpdatedWorkflow	true	"Workflow informations"
// @Success		200		{object}	docs_workflow.WorkflowUpdateWorkflowSuccessResponse
// @Success		400		{object}	docs_workflow.WorkflowUpdateWorkflowBadRequestResponse
// @Success		500		{object}	docs_workflow.WorkflowUpdateWorkflowInternalServerErrorResponse
// @Router			/workflows/{id} [put]
func (self *WorkflowHandler) updateWorkflow(context *gin.Context) {
	var workflow entities.UpdatedWorkflow
	workflowId := context.Param("id")

	err := context.ShouldBindJSON(&workflow)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": invalidRequestBodyMessage,
		})
		return
	}

	err = self.WorkflowService.UpdateWorkflow(workflowId, workflow)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not update workflow",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Workflow successfully updated",
	})
}

// @Summary		Delete Workflow
// @Description	Delete a user's workflow by specifying the workflow id
// @Tags			Workflows
// @Produce		json
// @Param        id     path     string  true  "Workflow id"
// @Success		200		{object}	docs_workflow.WorkflowDeleteWorkflowSuccessResponse
// @Failure		401 	{object}	docs_workflow.WorkflowDeleteWorkflowUnauthorizedResponse
// @Failure		500		{object}	docs_workflow.WorkflowDeleteWorkflowInternalServerErrorResponse
// @Router			/workflows/{id} [delete]
func (self *WorkflowHandler) deleteWorkflow(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")
	workflowId := context.Param("id")

	err := self.WorkflowService.DeleteWorkflow(email, connectionType, workflowId)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not delete workflow",
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Workflow deleted",
	})
}
