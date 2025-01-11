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
		"Skip track",
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

func (self *WorkflowService) executeAddTrackToPlaybackQueueSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/player/queue?uri=", "uri", "")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "POST", bearerType, accessToken, nil)
}

func (self *WorkflowService) executeSaveTrackSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/tracks?ids=", "id", "")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
}

func (self *WorkflowService) executeSaveAlbumSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/albums?ids=", "id", "")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
}

func (self *WorkflowService) executeSaveAudiobookSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/audiobooks?ids=", "id", "")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
}

func (self *WorkflowService) executeSaveEpisodeSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/episodes?ids=", "id", "")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
}

func (self *WorkflowService) executeSaveShowSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/shows?ids=", "id", "")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
}

func (self *WorkflowService) executeSetPlaybackVolumeSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/me/player/volume?volume_percent=", "volume", "")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
}

func (self *WorkflowService) executeFollowPlaylistSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/playlists/", "id", "/followers")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "PUT", bearerType, accessToken, nil)
}

func (self *WorkflowService) executeUnfollowPlaylistSpotifyRequest(workflow entities.Workflow, accessToken string) (*http.Response, error) {
	url, err := getSpotifyUrl(workflow, "https://api.spotify.com/v1/playlists/", "id", "/followers")
	if err != nil {
		return nil, err
	}
	return self.ServiceService.ExecuteApiRequest(url, "DELETE", bearerType, accessToken, nil)
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
	case "Skip track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/next", "POST", bearerType, accessToken, nil)
	case "Skip to previous track":
		resp, _ = self.ServiceService.ExecuteApiRequest("https://api.spotify.com/v1/me/player/previous", "POST", bearerType, accessToken, nil)
	case "Add track to playback queue":
		resp, _ = self.executeAddTrackToPlaybackQueueSpotifyRequest(workflow, accessToken)
	case "Save a track":
		resp, _ = self.executeSaveTrackSpotifyRequest(workflow, accessToken)
	case "Save an album":
		resp, _ = self.executeSaveAlbumSpotifyRequest(workflow, accessToken)
	case "Save an audiobook":
		resp, _ = self.executeSaveAudiobookSpotifyRequest(workflow, accessToken)
	case "Save an episode":
		resp, _ = self.executeSaveEpisodeSpotifyRequest(workflow, accessToken)
	case "Save a show":
		resp, _ = self.executeSaveShowSpotifyRequest(workflow, accessToken)
	case "Set playback volume":
		resp, _ = self.executeSetPlaybackVolumeSpotifyRequest(workflow, accessToken)
	case "Follow a playlist":
		resp, _ = self.executeFollowPlaylistSpotifyRequest(workflow, accessToken)
	case "Unfollow a playlist":
		resp, _ = self.executeUnfollowPlaylistSpotifyRequest(workflow, accessToken)
	}

	if resp != nil {
		defer resp.Body.Close()
		return nil
	}
	return fmt.Errorf("Unknown reaction")
}
