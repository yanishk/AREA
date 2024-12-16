package docs_workflow

type WorkflowExample struct {
	Id            string
	Name          string
	OwnerId       string
	ActionId      string
	ReactionId    string
	IsActivated   bool
	CreatedAt     string
	ActionParam   map[string]interface{} `json:"actionparam"`
	ReactionParam map[string]interface{} `json:"reactionparam"`
}

// Create Workflow Responses
type WorkflowCreateWorkflowSuccessResponse struct {
	Msg string `json:"success"example:"Successful workflow creation"`
}

type WorkflowCreateWorkflowBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body"`
}

type WorkflowCreateWorkflowUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string"`
}

type WorkflowCreateWorkflowInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not create workflow"`
}

// Retrieve User's Workflow Responses
type WorkflowRetrieveUserWorkflowsSuccessResponse struct {
	Workflow []WorkflowExample
}

type WorkflowRetrieveUserWorkflowsUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string"`
}

type WorkflowRetrieveUserWorkflowsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve user's workflows"`
}

// Modify Workflow Activation Status Responses
type WorkflowModifyWorkflowActivationStatusSuccessResponse struct {
	Msg string `json:"success"example:"Successfully updates workflow's activation status"`
}

type WorkflowModifyWorkflowActivationStatusBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body"`
}

type WorkflowModifyWorkflowActivationStatusInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not update workflow's activation status"`
}
