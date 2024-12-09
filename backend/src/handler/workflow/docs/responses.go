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
