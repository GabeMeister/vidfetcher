package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/GabeMeister/vidfetcher/api"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

const developerKey = "AIzaSyC9uXxwF4PxYilaOvPTDLdXAnToBwFvXcs"

// FetchChannelUploadPlaylistIDs - Fetches Youtube video data from youtube video id
func FetchChannelUploadPlaylistIDs(waitGroup *sync.WaitGroup, youtubeIDs string) chan string {
	waitGroup.Add(1)
	defer waitGroup.Done()

	ch := make(chan string)

	go func() {
		service := getYoutubeService()
		call := service.Channels.List("contentDetails").Id(youtubeIDs)
		response, err := call.Do()
		if err != nil {
			log.Fatal(err)
		}

		if len(response.Items) > 0 {
			for _, item := range response.Items {
				ch <- item.ContentDetails.RelatedPlaylists.Uploads
			}
		}
		close(ch)
	}()

	return ch
}

// FetchLatestVideoID - Get the latest video id of channel uploads playlist
// 		playlistID - id of playlist to fetch latest video of
func FetchLatestVideoID(playlistID string) string {
	service := getYoutubeService()
	call := service.PlaylistItems.
		List("contentDetails").
		PlaylistId(playlistID).
		MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(len(response.Items))

	latestVideoID := ""
	if len(response.Items) > 0 {
		latestVideoID = response.Items[0].ContentDetails.VideoId
	}

	return latestVideoID

}

// FetchChannelUploadCount - Fetches the number of uploads of a channel
//		channelID - id of the channel to fetch
func FetchChannelUploadCount(waitGroup *sync.WaitGroup, youtubeID string) chan api.Channel {
	waitGroup.Add(1)
	defer waitGroup.Done()

	ch := make(chan api.Channel)

	go func() {
		service := getYoutubeService()
		call := service.Channels.List("snippet,statistics").Id(youtubeID)
		response, err := call.Do()
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range response.Items {
			youtubeChannel := api.Channel{Title: item.Snippet.Title, VideoCount: item.Statistics.VideoCount}
			ch <- youtubeChannel
		}

		close(ch)
	}()

	return ch
}

// BreakYoutubeIDsIntoBatches breaks a list of ids into smaller batches
func BreakYoutubeIDsIntoBatches(youtubeIDs []string, batchSize int) (batchArr []string) {
	for idIndex := 0; idIndex < len(youtubeIDs); idIndex += batchSize {
		batchSize := GetBatchSize(len(youtubeIDs), idIndex, batchSize)
		batchArr = append(batchArr, strings.Join(youtubeIDs[idIndex:idIndex+batchSize], ","))
	}
	return batchArr
}

// GetBatchSize determines how big a batch can be without going over the bounds of the array
func GetBatchSize(arrSize int, index int, maxBatchSize int) int {
	batchSize := maxBatchSize
	if index+maxBatchSize > arrSize {
		batchSize = arrSize - index
	}
	return batchSize
}

func getYoutubeService() *youtube.Service {
	client := getYoutubeClient()
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	return service
}

func getYoutubeClient() *http.Client {
	return &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
}
