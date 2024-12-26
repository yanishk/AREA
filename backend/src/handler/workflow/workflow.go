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
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	workflow := private.Group("/workflow")
	{
		workflow.POST("/create", self.CreateWorkflow)
		workflow.GET("/retrieve", self.GetUserWorkflows)
		workflow.PUT("/update", self.UpdateWorkflow)
		workflow.DELETE("/delete", self.DeleteWorkflow)
	}
}

// @Summary		Create Workflow
// @Description	Create a Workflow
// @Tags			Workflows
// @Produce		json
// @Param			workflow	body		entities.NewWorkflow	true	"Workflow informations"
// @Success		200		{object}	docs_workflow.WorkflowCreateWorkflowSuccessResponse
// @Failure		400		{object}	docs_workflow.WorkflowCreateWorkflowBadRequestResponse
// @Failure		401		{object}	docs_workflow.WorkflowCreateWorkflowUnauthorizedResponse
// @Failure		500		{object}	docs_workflow.WorkflowCreateWorkflowInternalServerErrorResponse
// @Router			/workflow/create [post]
func (self *WorkflowHandler) CreateWorkflow(context *gin.Context) {
	var newWorkflow entities.NewWorkflow
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	err := context.ShouldBindJSON(&newWorkflow)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestBodyMessage,
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
// @Router			/workflow/retrieve [get]
func (self *WorkflowHandler) GetUserWorkflows(context *gin.Context) {
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
// @Param        id     query     string  true  "Workflow id"
// @Param			workflow	body		entities.UpdatedWorkflow	true	"Workflow informations"
// @Success		200		{object}	docs_workflow.WorkflowUpdateWorkflowSuccessResponse
// @Success		400		{object}	docs_workflow.WorkflowUpdateWorkflowBadRequestResponse
// @Success		500		{object}	docs_workflow.WorkflowUpdateWorkflowInternalServerErrorResponse
// @Router			/workflow/update [put]
func (self *WorkflowHandler) UpdateWorkflow(context *gin.Context) {
	var workflow entities.UpdatedWorkflow
	workflowId := context.Query("id")

	if workflowId == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workflow id",
		})
		return
	}

	err := context.ShouldBindJSON(&workflow)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestBodyMessage,
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
// @Param        id     query     string  true  "Workflow id"
// @Success		200		{object}	docs_workflow.WorkflowDeleteWorkflowSuccessResponse
// @Failure		400		{object}	docs_workflow.WorkflowDeleteWorkflowBadRequestResponse
// @Failure		401 	{object}	docs_workflow.WorkflowDeleteWorkflowUnauthorizedResponse
// @Failure		500		{object}	docs_workflow.WorkflowDeleteWorkflowInternalServerErrorResponse
// @Router			/workflow/delete [delete]
func (self *WorkflowHandler) DeleteWorkflow(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")
	workflowId := context.Query("id")

	if workflowId == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workflow id",
		})
		return
	}

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
