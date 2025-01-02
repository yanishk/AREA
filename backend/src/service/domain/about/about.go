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

func getAboutAction(action entities.Action) entities.AboutAction {
	var aboutAction entities.AboutAction

	aboutAction.Name = action.Name
	aboutAction.Description = action.Description
	return aboutAction
}

func getAboutReaction(reaction entities.Reaction) entities.AboutReaction {
	var aboutReaction entities.AboutReaction

	aboutReaction.Name = reaction.Name
	aboutReaction.Description = reaction.Description
	return aboutReaction
}

func (self *AboutService) getAboutActions(serviceId string) ([]entities.AboutAction, error) {
	var aboutActions []entities.AboutAction

	actions, err := self.ActionRepository.FindActionsByServiceId(serviceId)
	if err != nil {
		return aboutActions, err
	}

	for _, action := range actions {
		aboutAction := getAboutAction(action)
		aboutActions = append(aboutActions, aboutAction)
	}
	return aboutActions, nil
}
func (self *AboutService) getAboutReactions(serviceId string) ([]entities.AboutReaction, error) {
	var aboutReactions []entities.AboutReaction

	reactions, err := self.ReactionRepository.FindReactionsByServiceId(serviceId)
	if err != nil {
		return aboutReactions, err
	}

	for _, reaction := range reactions {
		aboutReaction := getAboutReaction(reaction)
		aboutReactions = append(aboutReactions, aboutReaction)
	}
	return aboutReactions, nil
}

func (self *AboutService) getAboutService(service entities.Service) (entities.AboutService, error) {
	var aboutService entities.AboutService

	actions, err := self.getAboutActions(service.Id)
	if err != nil {
		return aboutService, err
	}

	reactions, err := self.getAboutReactions(service.Id)
	if err != nil {
		return aboutService, err
	}

	aboutService.Name = service.Name
	aboutService.Actions = actions
	aboutService.Reactions = reactions
	return aboutService, nil
}

func (self *AboutService) getAboutServices() ([]entities.AboutService, error) {
	var aboutServices []entities.AboutService

	services, err := self.ServiceRepository.FindAllServices()
	if err != nil {
		return aboutServices, err
	}

	for _, service := range services {
		aboutService, err := self.getAboutService(service)
		if err != nil {
			return aboutServices, err
		}
		aboutServices = append(aboutServices, aboutService)
	}
	return aboutServices, nil
}

func (self *AboutService) GetAboutServer(about entities.About) (entities.About, error) {
	about.Server.CurrentTime = time.Now().Unix()

	services, err := self.getAboutServices()
	if err != nil {
		return about, err
	}
	about.Server.Services = services
	return about, nil
}
