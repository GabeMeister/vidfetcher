package tasks

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/GabeMeister/vidfetcher/api"
	"github.com/GabeMeister/vidfetcher/db"
)

const maxConcurrentGoRoutines = 50
const maxAPIResults = 50

// FetchAPIYoutubeChannelInfo fetches info from the Youtube API for the given youtube ids
func FetchAPIYoutubeChannelInfo(youtubeIDs []string) []api.YoutubeChannel {
	var waitGroup sync.WaitGroup
	youtubeIDBatches := api.BreakYoutubeIDsIntoBatches(youtubeIDs, maxAPIResults)

	fmt.Println("Amount of api calls to make:", len(youtubeIDBatches))

	var channelBatches []chan api.YoutubeChannel
	var count int
	var youtubeChannelData []api.YoutubeChannel

	for batchStart := 0; batchStart < len(youtubeIDBatches); batchStart += maxConcurrentGoRoutines {
		channelBatches = nil
		batchSize := api.GetBatchSize(len(youtubeIDBatches), batchStart, maxConcurrentGoRoutines)

		for batchIndex := batchStart; batchIndex < batchStart+batchSize; batchIndex++ {
			ch := api.FetchChannelData(&waitGroup, youtubeIDBatches[batchIndex])
			channelBatches = append(channelBatches, ch)
		}
		mergedChannel := api.MergeChannels(channelBatches)

		for item := range mergedChannel {
			count++
			// fmt.Println(count, item)
			youtubeChannelData = append(youtubeChannelData, item)
		}
		waitGroup.Wait()
		fmt.Println("api calls made:", batchStart+batchSize)
	}

	return youtubeChannelData
}

// FetchAllVideosForChannel fetches all the video uploads of the receiver Youtube Channel
func FetchAllVideosForChannel(youtubeDB *sql.DB, youtubeChannel *api.YoutubeChannel) {
	fmt.Println("The following channels are out of date: ")
	if AreVideosOutOfDate(youtubeDB, youtubeChannel) {
		fmt.Println(youtubeChannel.String())
	}
}

// AreVideosOutOfDate determines if there needs to be new videos fetched for a particular channel
func AreVideosOutOfDate(youtubeDB *sql.DB, youtubeChannel *api.YoutubeChannel) bool {
	// Get channel id from youtube id
	id := db.SelectChannelIDFromYoutubeID(youtubeDB, youtubeChannel.YoutubeID)

	// Compare to YoutubeChannel object
	videoCountDB := db.SelectVideoCountOfChannel(youtubeDB, id)

	// Video counts are out of date if the count from the database
	// doesn't match the count from the api
	isOutOfDate := videoCountDB != youtubeChannel.VideoCount
	if isOutOfDate {
		fmt.Printf("'%s' out of date. DB: %d API: %d\n", youtubeChannel.Title, videoCountDB, youtubeChannel.VideoCount)
	}

	return isOutOfDate
}
