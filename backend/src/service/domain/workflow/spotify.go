package workflow_service

import (
	"fmt"
	"net/http"

	"backend/src/entities"
)

func spotifyReactions() []string {
	return []string{
		"Start playback",
		"Pause playback",
		"Activate playback shuffle",
		"Deactivate playback shuffle",
		"Skip to next track",
		"Skip to previous track",
		"Add track to playback queue",
		"Save a track",
		"Save an album",
		"Save an audiobook",
		"Save an episode",
		"Save a show",
		"Set playback volume",
		"Follow a playlist",
		"Unfollow a playlist",
	}
}

func spotifyReactionsToRequestParameters() map[string]entities.SpotifyRequestParameters {
	return map[string]entities.SpotifyRequestParameters{
		"Start playback":              {"https://api.spotify.com/v1/me/player/play", "", "", "PUT"},
		"Pause playback":              {"https://api.spotify.com/v1/me/player/pause", "", "", "PUT"},
		"Activate playback shuffle":   {"https://api.spotify.com/v1/me/player/shuffle?state=true", "", "", "PUT"},
		"Deactivate playback shuffle": {"https://api.spotify.com/v1/me/player/shuffle?state=false", "", "", "PUT"},
		"Skip to previous track":      {"https://api.spotify.com/v1/me/player/next", "", "", "POST"},
		"Skip to next track":          {"https://api.spotify.com/v1/me/player/next", "", "", "POST"},
		"Add track to playback queue": {"https://api.spotify.com/v1/me/player/queue?uri=", "uri", "", "POST"},
		"Save a track":                {"https://api.spotify.com/v1/me/tracks?ids=", "id", "", "PUT"},
		"Save an album":               {"https://api.spotify.com/v1/me/albums?ids=", "id", "", "PUT"},
		"Save an audiobook":           {"https://api.spotify.com/v1/me/audiobooks?ids=", "id", "", "PUT"},
		"Save an episode":             {"https://api.spotify.com/v1/me/episodes?ids=", "id", "", "PUT"},
		"Save a show":                 {"https://api.spotify.com/v1/me/shows?ids=", "id", "", "PUT"},
		"Set playback volume":         {"https://api.spotify.com/v1/me/player/volume?volume_percent=", "volume", "", "PUT"},
		"Follow a playlist":           {"https://api.spotify.com/v1/playlists/", "id", "/followers", "PUT"},
		"Unfollow a playlist":         {"https://api.spotify.com/v1/playlists/", "id", "/followers", "DELETE"},
	}
}

func getSpotifyUrl(workflow entities.Workflow, baseUrl, parameterName, endUrl string) (string, error) {
	parameter, parameterExists := workflow.ReactionParam[parameterName]
	if !parameterExists {
		return "", fmt.Errorf("Incorrect parameter")
	}

	return (baseUrl + parameter.(string) + endUrl), nil
}

func (self *WorkflowService) executeSpotifyRequest(workflow entities.Workflow, requestParameters entities.SpotifyRequestParameters, accessToken string) (*http.Response, error) {
	var err error
	url := requestParameters.BaseUrl

	if requestParameters.ParameterName != "" {
		url, err = getSpotifyUrl(workflow, requestParameters.BaseUrl, requestParameters.ParameterName, requestParameters.EndUrl)
		if err != nil {
			return nil, err
		}
	}
	return self.ServiceService.ExecuteApiRequest(url, requestParameters.Method, bearerType, accessToken, nil)
}

func (self *WorkflowService) checkSpotifyReactions(workflow entities.Workflow) error {
	reactionFound, errReaction := self.ReactionRepository.FindReactionById(workflow.ReactionId)
	if errReaction != nil {
		return fmt.Errorf(errorRetrievingReaction)
	}

	accessToken, err := self.refreshTokenForService("Spotify", reactionFound.Name, spotifyReactions(), workflow)
	if err != nil {
		return fmt.Errorf(errorUpdatingToken)
	}

	spotifyReactionsToRequestParameters := spotifyReactionsToRequestParameters()
	requestParameters, requestParametersExists := spotifyReactionsToRequestParameters[reactionFound.Name]
	if !requestParametersExists {
		return fmt.Errorf("Unknown reaction")
	}

	resp, err := self.executeSpotifyRequest(workflow, requestParameters, accessToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
