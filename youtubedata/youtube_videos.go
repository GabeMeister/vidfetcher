package youtubedata

import (
	"encoding/json"
	"fmt"
	"log"

	youtube "google.golang.org/api/youtube/v3"
)

// Video contains Youtube video data
type Video struct {
	apiVideo  *youtube.Video
	ChannelID int
}

func (v *Video) String() string {
	return fmt.Sprintf("%s", v.Title())
}

// JSONString returns the json encoding of the youtubedata.Video object
func (v *Video) JSONString() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Fatalln("Unable to marshal youtube video:", v.String())
	}
	return string(bytes)
}

// Title is the title of the video
func (v *Video) Title() string {
	if !v.apiVideoExists() {
		log.Fatalln("apiVideo is nil, cannot access video title")
	}
	if !v.snippetExists() {
		log.Fatalln("apiVideo's snippet is nil, cannot access video title")
	}
	return v.apiVideo.Snippet.Title
}

// YoutubeID is the video's 32 character id string recognized by Youtube
func (v *Video) YoutubeID() string {
	if !v.apiVideoExists() {
		log.Fatalln("apiVideo is nil, cannot access YoutubeID")
	}
	if !v.snippetExists() {
		log.Fatalln("apiVideo's snippet is nil, cannot access YoutubeID")
	}
	return v.apiVideo.Id
}

func (v *Video) apiVideoExists() bool {
	return v.apiVideo != nil
}

func (v *Video) snippetExists() bool {
	return v.apiVideo != nil && v.apiVideo.Snippet != nil
}

// FetchVideoInfo fetches video data from playlist items found in playlistItems
func FetchVideoInfo(playlistItems []PlaylistItem) []Video {
	// TODO
	return []Video{}
}
