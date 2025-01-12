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

func getSpotifyUrl(workflow entities.Workflow, baseUrl, parameterName, endUrl string) (string, error) {
	parameter, parameterExists := workflow.ReactionParam[parameterName]
	if !parameterExists {
		return "", fmt.Errorf("Incorrect parameter")
	}

	return (baseUrl + parameter.(string) + endUrl), nil
}

func (self *WorkflowService) executeSpotifyRequestWithParam(workflow entities.Workflow, baseUrl, parameterName, endUrl, method, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, baseUrl, parameterName, endUrl)
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, method, bearerType, accessToken, nil)
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

	var resp *http.Response
	switch reactionFound.Name {
	case "Start playback":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/play", "PUT", bearerType, accessToken, nil)
	case "Pause playback":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/pause", "PUT", bearerType, accessToken, nil)
	case "Activate playback shuffle":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/shuffle?state=true", "PUT", bearerType, accessToken, nil)
	case "Deactivate playback shuffle":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/shuffle?state=false", "PUT", bearerType, accessToken, nil)
	case "Skip to next track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/next", "POST", bearerType, accessToken, nil)
	case "Skip to previous track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/previous", "POST", bearerType, accessToken, nil)
	case "Add track to playback queue":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/me/player/queue?uri=", "uri", "", "POST", accessToken)
	case "Save a track":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/me/tracks?ids=", "id", "", "PUT", accessToken)
	case "Save an album":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/me/albums?ids=", "id", "", "PUT", accessToken)
	case "Save an audiobook":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/me/audiobooks?ids=", "id", "", "PUT", accessToken)
	case "Save an episode":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/me/episodes?ids=", "id", "", "PUT", accessToken)
	case "Save a show":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/me/shows?ids=", "id", "", "PUT", accessToken)
	case "Set playback volume":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/me/player/volume?volume_percent=", "volume", "", "PUT", accessToken)
	case "Follow a playlist":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/playlists/", "id", "/followers", "PUT", accessToken)
	case "Unfollow a playlist":
		resp, _ = self.executeSpotifyRequestWithParam(workflow, "https://api.spotify.com/v1/playlists/", "id", "/followers", "DELETE", accessToken)
	}

	if resp != nil {
		defer resp.Body.Close()
		return nil
	}
	return fmt.Errorf("Unknown reaction")
}
