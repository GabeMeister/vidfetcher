package api

import (
	"log"

	"github.com/GabeMeister/vidfetcher/youtubedata"
)

// FetchVideoInfo fetches video data from playlist items found in playlistItems
// Note: this is called in a go routine, so it synchronously fetches video results
func FetchVideoInfo(youtubeIDs []string, youtubeChannel *youtubedata.Channel) []youtubedata.Video {
	youtubeVideos := []youtubedata.Video{}

	youtubeIDBatches := BreakYoutubeIDsIntoBatches(youtubeIDs, 50)

	for i, idBatch := range youtubeIDBatches {
		log.Printf("Fetching video info #%d for %s", i+1, youtubeChannel.Title())

		service := GetYoutubeService()

		call := service.Videos.
			List("snippet,contentDetails,statistics").
			Id(idBatch).
			MaxResults(50)
		response, err := call.Do()
		if err != nil {
			log.Fatalln("Couldn't fetch video info %v", err)
		}

		for _, vid := range response.Items {
			youtubeVideos = append(youtubeVideos, youtubedata.Video{APIVideo: vid, ChannelID: youtubeChannel.ChannelID})
		}
	}

	return youtubeVideos
}
