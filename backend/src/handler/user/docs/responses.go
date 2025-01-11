package docs_user

type UserInfosExample struct {
	Email          string `json:"email"`
	CreatedAt      string `json:"createdat"`
	Timezone       string `json:"timezone"`
	ConnectionType string `json:"connectiontype"`
}

// General Responses
type UserInvalidBodyResponse struct {
	Msg string `json:"error"example:"Invalid request body"`
}

type UserInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Internal server error"`
}

// Register Responses
type UserRegisterSuccessResponse struct {
	Msg string `json:"success"example:"New user created"`
}

type UserRegisterConflictResponse struct {
	Msg string `json:"error"example:"Email address already used"`
}

type UserRegisterBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body-Connection type doesn't exist"`
}

// Login Responses
type UserLoginSuccessResponse struct {
	Msg string `json:"success"example:"Connection successful"`
}

type UserLoginUnauthorizedResponse struct {
	Msg string `json:"error"example:"Could not find requested user-Wrong password"`
}
type UserLoginTokenErrorResponse struct {
	Msg string `json:"error"example:"Error creating token"`
}

// Login Callback Responses
type UserLoginCallbackSuccessResponse struct {
	Msg string `json:"success"example:"Connection successful"`
}

type UserLoginCallbackBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body-Invalid code authorization"`
}

type UserLoginCallbackInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to connect with requested service"`
}

// Get User Responses
type UserGetUserSuccessResponse struct {
	User UserInfosExample
}

type UserGetUserUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string"`
}

type UserGetUserInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not find user"`
}

// Logout User Responses
type UserLogoutUserSuccessResponse struct {
	Msg string `json:"success"example:"Logout successful"`
}

type UserLogoutUserInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not find user"`
}

// Modify Password Responses
type UserModifyPasswordSuccessResponse struct {
	Msg string `json:"success"example:"Password modified"`
}

type UserModifyPasswordUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string-Connection type not found in token-Connection type is not a valid string"`
}

type UserModifyPasswordOldPasswordNotValidResponse struct {
	Msg string `json:"error"example:"Old password is incorrect"`
}

type UserModifyPasswordInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to hash password-Could not modify the password"`
}

type UserModifyPasswordBadRequestResponse struct {
	Msg string `json:"error"example:"Invalid request body-Could not find requested user"`
}

// Delete Account Responses
type UserDeleteAccountSuccessResponse struct {
	Msg string `json:"success"example:"Account deleted"`
}

type UserDeleteAccountUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string-Connection type not found in token-Connection type is not a valid string"`
}

type UserDeleteAccountInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not delete account"`
}
