package docs_userservice

// Service Callback Responses
type UserServiceServiceCallbackSuccessResponse struct {
	Msg string `json:"success"example:"Token generated"`
}

type UserServiceServiceCallbackBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid code authorization-Invalid app type"`
}

type UserServiceServiceCallbackInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to get email-Failed to update token"`
}
