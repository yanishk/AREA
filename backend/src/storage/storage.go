package storage

import (
	"backend/src/entities"
)

type UserRepository interface {
	CreateUser(email string, password string, connectionType string) error
	FindUserByEmail(email string, connectionType string) (entities.User, error)
	FindUserById(userId string) (entities.User, error)
	UpdateUser(email string, password string, connectionType string) error
	DeleteUser(email string, connectionType string) error
}

type ServiceRepository interface {
	CreateService(name string, color string, logo string) error
	FindServiceById(id string) (entities.Service, error)
	FindServiceByName(name string) (entities.Service, error)
	FindAllServices() ([]entities.Service, error)
	FindActionsServices() ([]entities.Service, error)
	FindReactionsServices() ([]entities.Service, error)
}

type UserServiceRepository interface {
	CreateUserService(userId string, token string, tokenRefresh string, expiryDate string, serviceId string) error
	FindUserServiceByServiceIdandUserId(userId, serviceId string) (entities.UserService, error)
	UpdateUserServiceByServiceIdAndUserId(userId string, accessToken string, refreshToken string, expiryDate string, serviceId string) error
	DeleteUserServiceByUserId(userId string) error
}

type ReactionRepository interface {
	CreateReaction(name string, description string, serviceId string, nbParam int) error
	FindReactionById(id string) (entities.Reaction, error)
	FindReactionByName(name string) (entities.Reaction, error)
	FindReactionsByServiceId(serviceId string) ([]entities.Reaction, error)
}

type ActionRepository interface {
	CreateAction(name string, description string, serviceId string, nbParam int) error
	FindActionById(id string) (entities.Action, error)
	FindActionByName(name string) (entities.Action, error)
	FindActionsByServiceId(serviceId string) ([]entities.Action, error)
}

type WorkflowRepository interface {
	CreateWorkflow(name, ownerId string, actionId string, reactionId string, actionParam map[string]interface{}, reactionParam map[string]interface{}, actionData map[string]interface{}) error
	FindWorkflowById(id string) (entities.Workflow, error)
	FindWorkflowsByActionId(actionId string) ([]entities.Workflow, error)
	FindWorkflowsByOwnerId(ownerId string) ([]entities.Workflow, error)
	UpdateWorkflow(id string, updatedWorkflow entities.Workflow) error
	DeleteWorkflow(id string, ownerId string) error
	DeleteWorkflowByOwnerId(ownerId string) error
}

type Repository struct {
	UserRepository        UserRepository
	ServiceRepository     ServiceRepository
	UserServiceRepository UserServiceRepository
	ReactionRepository    ReactionRepository
	ActionRepository      ActionRepository
	WorkflowRepository    WorkflowRepository
}
