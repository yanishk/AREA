package docs_userservice

type OAuth2SuccessResponse struct {
	Msg string `json:"auth-url"example:"url-example.com"`
}

type OAuth2BadRequestResponse struct {
	Msg string `json:"error"example:"Invalid callback type"`
}

type ServiceCallBackSuccessResponse struct {
	Msg string `json:"success"example:"Token generated"`
}

type ServiceCallBackBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid code authorization"`
}

type ServiceCallBackInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to get email-Failed to update token"`
}
