package api

import (
	"github.com/GabeMeister/vidfetcher/youtubedata"
)

// GetYoutubeIDsFromPlaylistItems converts a list of playlist items to
// just their youtube ids
func GetYoutubeIDsFromPlaylistItems(playlistItems []youtubedata.PlaylistItem) []string {
	youtubeIDs := []string{}
	for _, item := range playlistItems {
		youtubeIDs = append(youtubeIDs, item.YoutubeID())
	}
	return youtubeIDs
}
