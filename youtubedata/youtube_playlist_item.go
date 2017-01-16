package youtubedata

import (
	"log"

	youtube "google.golang.org/api/youtube/v3"
)

// PlaylistItem represents an item in a youtube channel uploads playlist
type PlaylistItem struct {
	APIPlaylistItem *youtube.PlaylistItem
}

// YoutubeID is the youtube id of the video for the playlist item
func (p *PlaylistItem) YoutubeID() string {
	if p.APIPlaylistItem == nil || p.APIPlaylistItem.Snippet == nil {
		log.Fatalln("playlist item is nil or doesn't contain snippet, cannot get youtube id")
	}

	return p.APIPlaylistItem.Snippet.ResourceId.VideoId
}
