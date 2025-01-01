package about_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/src/entities"
	"backend/src/service"
)

type AboutHandler struct {
	AboutService service.AboutService
}

func NewAboutHandler(AboutService service.AboutService, router *gin.Engine) *AboutHandler {
	handler := &AboutHandler{AboutService: AboutService}
	handler.SetupRoutes(router)
	return handler
}

func (self *AboutHandler) SetupRoutes(router *gin.Engine) {
	self.PublicRoutes(router)
}

func (self *AboutHandler) PublicRoutes(router *gin.Engine) {
	router.GET("/about.json", self.GetAbout)
}

// @Summary		About.json
// @Description	Get about.json file
// @Tags			About.json
// @Produce		json
// @Success		200		{object}	docs_about.AboutGetAboutSuccessResponse
// @Failure		500		{object}	docs_about.AboutGetAboutInternalServerErrorResponse
// @Router			/about.json [get]
func (self *AboutHandler) GetAbout(context *gin.Context) {
	var about entities.About

	about.Client.Host = context.ClientIP()
	about, err := self.AboutService.GetAboutServer(about)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve informations about the server",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, about)
}
