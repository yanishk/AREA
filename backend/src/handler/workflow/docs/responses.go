package docs_workflow

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

// Modify Workflow Activation Status Responses
type WorkflowModifyWorkflowActivationStatusSuccessResponse struct {
	Msg string `json:"success"example:"Successfully updates workflow's activation status"`
}

type WorkflowModifyWorkflowActivationStatusBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body-Could not find specified workflow"`
}

type WorkflowModifyWorkflowActivationStatusInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not update workflow's activation status"`
}
