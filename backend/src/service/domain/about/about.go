package about_service

import (
	"time"

	"backend/src/entities"
	_ "backend/src/handler/about/docs"
	"backend/src/storage"
)

type AboutService struct {
	ServiceRepository  storage.ServiceRepository
	ActionRepository   storage.ActionRepository
	ReactionRepository storage.ReactionRepository
}

func NewAboutService(ServiceRepository storage.ServiceRepository, ActionRepository storage.ActionRepository,
	ReactionRepository storage.ReactionRepository) *AboutService {
	return &AboutService{
		ServiceRepository:  ServiceRepository,
		ActionRepository:   ActionRepository,
		ReactionRepository: ReactionRepository,
	}
}

func (self *AboutService) GetAboutAction(action entities.Action) (entities.AboutAction) {
	var aboutAction entities.AboutAction

	aboutAction.Name = action.Name
	aboutAction.Description = action.Description
	return aboutAction
}

func (self *AboutService) GetAboutReaction(reaction entities.Reaction) (entities.AboutReaction) {
	var aboutReaction entities.AboutReaction

	aboutReaction.Name = reaction.Name
	aboutReaction.Description = reaction.Description
	return aboutReaction
}

func (self *AboutService) GetAboutActions(serviceId string) ([]entities.AboutAction, error) {
	var aboutActions []entities.AboutAction

	actions, err := self.ActionRepository.FindActionByServiceId(serviceId)
	if err != nil {
		return aboutActions, err
	}

	for _, action := range actions {
		aboutAction := self.GetAboutAction(action)
		aboutActions = append(aboutActions, aboutAction)
	}
	return aboutActions, nil
}
func (self *AboutService) GetAboutReactions(serviceId string) ([]entities.AboutReaction, error) {
	var aboutReactions []entities.AboutReaction

	reactions, err := self.ReactionRepository.FindReactionByServiceId(serviceId)
	if err != nil {
		return aboutReactions, err
	}

	for _, reaction := range reactions {
		aboutReaction := self.GetAboutReaction(reaction)
		aboutReactions = append(aboutReactions, aboutReaction)
	}
	return aboutReactions, nil
}

func (self *AboutService) GetAboutService(service entities.Service) (entities.AboutService, error) {
	var aboutService entities.AboutService
	var err error

	aboutService.Actions, err = self.GetAboutActions(service.Id)
	if err != nil {
		return aboutService, err
	}

	aboutService.Reactions, err = self.GetAboutReactions(service.Id)
	if err != nil {
		return aboutService, err
	}

	aboutService.Name = service.Name
	return aboutService, nil
}

func (self *AboutService) GetAboutServices() ([]entities.AboutService, error) {
	var aboutServices []entities.AboutService

	services, err := self.ServiceRepository.FindServices()
	if err != nil {
		return aboutServices, err
	}

	for _, service := range services {
		aboutService, err := self.GetAboutService(service)
		if err != nil {
			return aboutServices, err
		}
		aboutServices = append(aboutServices, aboutService)
	}
	return aboutServices, nil
}

func (self *AboutService) GetAboutServer(about entities.About) (entities.About, error) {
	about.Server.CurrentTime = time.Now().Unix()
	var err error

	about.Server.Services, err = self.GetAboutServices()
	if err != nil {
		return about, err
	}
	return about, nil
}
