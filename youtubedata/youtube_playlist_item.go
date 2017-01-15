package youtubedata

import (
	"log"

	"github.com/GabeMeister/vidfetcher/api"
	youtube "google.golang.org/api/youtube/v3"
)

// PlaylistItem represents an item in a youtube channel uploads playlist
type PlaylistItem struct {
	APIPlaylistItem *youtube.PlaylistItem
}

// FetchChannelUploads fetches up to 50 videos of a channel and returns
// the response
func FetchChannelUploads(youtubeChannel *Channel, pageToken string) *youtube.PlaylistItemListResponse {
	service := api.GetYoutubeService()
	call := service.PlaylistItems.
		List("snippet").
		PlaylistId(youtubeChannel.UploadsPlaylistID()).
		PageToken(pageToken).
		MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	return response
}
