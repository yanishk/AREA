package docs_userservice

type ServiceCallbackSuccessResponse struct {
	Msg string `json:"success"example:"Token generated"`
}

type ServiceCallbackBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid code authorization"`
}

type ServiceCallbackInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to get email-Failed to update token"`
}
