package domain

import (
	"backend/src/service"
	about_service "backend/src/service/domain/about"
	service_service "backend/src/service/domain/service"
	user_service "backend/src/service/domain/user"
	user_service_service "backend/src/service/domain/userservice"
	workflow_service "backend/src/service/domain/workflow"
	"backend/src/storage"
)

func New(repositories *storage.Repository) *service.Service {
	serviceService := service_service.NewServiceService(repositories.ServiceRepository, repositories.UserRepository, repositories.ActionRepository, repositories.WorkflowRepository, repositories.ReactionRepository)
	userService := user_service.NewUserService(repositories.UserRepository, repositories.ServiceRepository, repositories.UserServiceRepository, repositories.WorkflowRepository, serviceService)
	userServiceService := user_service_service.NewUserServiceService(repositories.ServiceRepository, repositories.UserRepository, repositories.UserServiceRepository, serviceService)
	workflowService := workflow_service.NewWorkflowService(repositories.WorkflowRepository, repositories.UserRepository, repositories.ActionRepository, repositories.ReactionRepository, serviceService, userServiceService)
	aboutService := about_service.NewAboutService(repositories.ServiceRepository, repositories.ActionRepository, repositories.ReactionRepository)

	return &service.Service{
		ServiceService:     serviceService,
		UserService:        userService,
		UserServiceService: userServiceService,
		WorkflowService:    workflowService,
		AboutService:       aboutService,
	}
}
