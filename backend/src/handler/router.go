package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"

	service_handler "backend/src/handler/service"
	user_handler "backend/src/handler/user"
	"backend/src/service"
)

type UserHandler interface {
	SetupRoutes(router *gin.Engine)
}

type ServiceHandler interface {
	SetupRoutes(router *gin.Engine)
}

type Handler struct {
	UserHandler    UserHandler
	ServiceHandler ServiceHandler
	router         *gin.Engine
}

func New(services *service.Service) *Handler {
	router := gin.Default()
	return &Handler{
		UserHandler:    user_handler.NewUserHandler(services.UserService, router),
		ServiceHandler: service_handler.NewServiceHandler(services.ServiceService, router),
		router:         router,
	}
}

func (self *Handler) Run() error {

	corsAllow := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := corsAllow.Handler(self.router)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return err
	}
	return nil
}
