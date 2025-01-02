package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	about_handler "backend/src/handler/about"
	service_handler "backend/src/handler/service"
	user_handler "backend/src/handler/user"
	user_service_handler "backend/src/handler/userservice"
	workflow_handler "backend/src/handler/workflow"
	"backend/src/service"
)

type UserHandler interface {
}

type ServiceHandler interface {
}

type UserServiceHandler interface {
}

type WorkflowHandler interface {
}

type AboutHandler interface {
}

type Handler struct {
	UserHandler        UserHandler
	ServiceHandler     ServiceHandler
	UserServiceHandler UserServiceHandler
	WorkflowHandler    WorkflowHandler
	AboutHandler       AboutHandler
	router             *gin.Engine
}

func New(services *service.Service) *Handler {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return &Handler{
		UserHandler:        user_handler.NewUserHandler(services.UserService, router),
		UserServiceHandler: user_service_handler.NewUserServiceHandler(services.UserServiceService, router),
		ServiceHandler:     service_handler.NewServiceHandler(services.ServiceService, services.UserService, router),
		WorkflowHandler:    workflow_handler.NewWorkflowHandler(services.WorkflowService, services.UserService, router),
		AboutHandler:       about_handler.NewAboutHandler(services.AboutService, router),
		router:             router,
	}
}

func (self *Handler) Run() error {
	corsAllow := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := corsAllow.Handler(self.router)
	err := http.ListenAndServeTLS(":8080", os.Getenv("CERTIFICATE"), os.Getenv("KEY"), handler)
	if err != nil {
		return err
	}
	return nil
}
