package docs_userservice

type OAuth2Response struct {
	Msg string `json:"auth-url"example:"url-example.com"`
}

type ServiceCallBackSuccessResponse struct {
	Msg string `json:"success"example:"Token generated"`
}

type ServiceCallBackBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid code authorization-Could not find requested user-Could not find requested service"`
}

type ServiceCallBackInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to exchange token-Failed to get email"`
}
