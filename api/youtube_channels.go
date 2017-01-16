package api

import (
	"log"
	"sync"

	youtube "google.golang.org/api/youtube/v3"

	"github.com/GabeMeister/vidfetcher/youtubedata"
)

// FetchChannelDataFromAPI - Fetches the number of uploads of a channel
func FetchChannelDataFromAPI(waitGroup *sync.WaitGroup, youtubeIDCommaText string) chan youtubedata.Channel {
	waitGroup.Add(1)
	defer waitGroup.Done()

	ch := make(chan youtubedata.Channel)

	go func() {
		service := GetYoutubeService()
		call := service.Channels.List("snippet,statistics,contentDetails").Id(youtubeIDCommaText)
		response, err := call.Do()
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range response.Items {
			youtubeChannel := youtubedata.Channel{APIChannel: item}
			ch <- youtubeChannel
		}

		close(ch)
	}()

	return ch
}

// FetchChannelUploads fetches up to 50 videos of a channel and returns the whole
//  playlist item list response
func FetchChannelUploads(youtubeChannel *youtubedata.Channel, pageToken string) *youtube.PlaylistItemListResponse {
	service := GetYoutubeService()
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
