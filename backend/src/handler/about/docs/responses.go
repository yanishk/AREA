package docs_about

import "backend/src/entities"

// Get About Responses
type AboutGetAboutSuccessResponse struct {
	About entities.About
}

type AboutGetAboutInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not retrieve informations about the server"`
}
