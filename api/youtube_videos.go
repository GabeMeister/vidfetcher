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
		log.Printf("Video info #%d for %s\n", i+1, youtubeChannel.Title())

		attempts := 0
		success := false

		// Attempt up to 3 times to fetch video data. If not able to fetch, we just log and continue
		for attempts < 3 && success == false {
			service := GetYoutubeService()

			call := service.Videos.
				List("snippet,contentDetails,statistics").
				Id(idBatch).
				MaxResults(50)
			response, err := call.Do()
			if err != nil {
				log.Printf("Couldn't fetch video info %v\n", err)
				attempts++
			} else {
				success = true
				for _, vid := range response.Items {
					youtubeVideos = append(youtubeVideos, youtubedata.Video{APIVideo: vid, ChannelID: youtubeChannel.ChannelID})
				}
			}
		}

	}

	return youtubeVideos
}
