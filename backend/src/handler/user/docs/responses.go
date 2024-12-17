package docs_user

// General Responses
type FailureInvalidBodyResponse struct {
	Msg string `json:"error"example:"Invalid request body"`
}

type FailureInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Internal server error"`
}

// Register Responses
type RegisterSuccessResponse struct {
	Msg string `json:"success"example:"New user created"`
}

type RegisterFailureConflictResponse struct {
	Msg string `json:"error"example:"Email address already used"`
}

type RegisterFailureBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body-Connection type doesn't exist"`
}

// Login Responses
type LoginSuccessResponse struct {
	Msg string `json:"success"example:"Connection successful"`
}

type LoginFailureUnauthorizedResponse struct {
	Msg string `json:"error"example:"Could not find requested user-Wrong password"`
}
type LoginFailureTokenErrorResponse struct {
	Msg string `json:"error"example:"Error creating token"`
}

// Modify Password Responses
type ModifyPasswordSuccessResponse struct {
	Msg string `json:"success"example:"Password modified"`
}

type ModifyPasswordFailureEmailNotValidResponse struct {
	Msg string `json:"error"example:"Email is not a valid string"`
}

type ModifyPasswordFailureOldPasswordNotValidResponse struct {
	Msg string `json:"error"example:"Old password is incorrect"`
}

type ModifyPasswordFailureInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to hash password-Could not modify the password"`
}

type ModifyPasswordBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body-Could not find requested user"`
}

// Delete Account Responses
type DeleteAccountSuccessResponse struct {
	Msg string `json:"success"example:"Account deleted"`
}

type DeleteAccountFailureUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string"`
}

type DeleteAccountFailureInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not delete account"`
}

// Get User Email Responses
type GetUserEmailSuccessResponse struct {
	Msg string `json:"email"example:"email@example.com"`
}

type GetUserEmailUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string"`
}
