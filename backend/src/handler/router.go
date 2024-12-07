package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	service_handler "backend/src/handler/service"
	user_handler "backend/src/handler/user"
	user_service_handler "backend/src/handler/userservice"
	"backend/src/service"
)

type UserHandler interface {
}

type ServiceHandler interface {
	CallTimeApi() error
}

type UserServiceHandler interface {
}

type Handler struct {
	UserHandler        UserHandler
	ServiceHandler     ServiceHandler
	UserServiceHandler UserServiceHandler
	router             *gin.Engine
}

func New(services *service.Service) *Handler {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return &Handler{
		UserHandler:        user_handler.NewUserHandler(services.UserService, router),
		UserServiceHandler: user_service_handler.NewUserServiceHandler(services.UserServiceService, router),
		ServiceHandler:     service_handler.NewServiceHandler(services.ServiceService, services.UserServiceService, services.UserService, router),
		router:             router,
	}
}

func (self *Handler) Run() error {

	corsAllow := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := corsAllow.Handler(self.router)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return err
	}
	return nil
}
