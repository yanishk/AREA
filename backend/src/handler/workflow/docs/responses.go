package docs_workflow

type WorkflowExample struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	OwnerId       string                 `json:"ownerid"`
	ActionId      string                 `json:"actionid"`
	ReactionId    string                 `json:"reactionid"`
	IsActivated   bool                   `json:"isactivated"`
	CreatedAt     string                 `json:"createdat"`
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
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string-Connection type not found in token-Connection type is not a valid string"`
}

type WorkflowCreateWorkflowInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not create workflow"`
}

// Retrieve User's Workflow Responses
type WorkflowRetrieveUserWorkflowsSuccessResponse struct {
	Workflows []WorkflowExample
}

type WorkflowRetrieveUserWorkflowsUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string-Connection type not found in token-Connection type is not a valid string"`
}

type WorkflowRetrieveUserWorkflowsInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve user's workflows"`
}

// Update Workflow Responses
type WorkflowUpdateWorkflowSuccessResponse struct {
	Msg string `json:"success"example:"Workflow successfully updated"`
}

type WorkflowUpdateWorkflowBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body"`
}

type WorkflowUpdateWorkflowInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not update workflow"`
}

// Delete Workflow Responses
type WorkflowDeleteWorkflowSuccessResponse struct {
	Msg string `json:"success"example:"Workflow deleted"`
}

type WorkflowDeleteWorkflowUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string-Connection type not found in token-Connection type is not a valid string"`
}

type WorkflowDeleteWorkflowInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not delete workflow"`
}
