package tasks

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/GabeMeister/vidfetcher/api"
	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/youtubedata"
)

const maxConcurrentGoRoutines = 50
const maxAPIResults = 50

// FetchYoutubeChannelInfoFromAPI fetches info from the Youtube API for the given youtube ids
func FetchYoutubeChannelInfoFromAPI(youtubeIDs []string) []youtubedata.Channel {
	var waitGroup sync.WaitGroup
	youtubeIDBatches := api.BreakYoutubeIDsIntoBatches(youtubeIDs, maxAPIResults)

	fmt.Println("Amount of api calls to make:", len(youtubeIDBatches))

	var channelsInBatch []chan youtubedata.Channel
	var count int
	var youtubeChannelData []youtubedata.Channel

	for batchStart := 0; batchStart < len(youtubeIDBatches); batchStart += maxConcurrentGoRoutines {
		channelsInBatch = nil
		batchSize := api.GetBatchSize(len(youtubeIDBatches), batchStart, maxConcurrentGoRoutines)

		for batchIndex := batchStart; batchIndex < batchStart+batchSize; batchIndex++ {
			ch := youtubedata.FetchChannelDataFromAPI(&waitGroup, youtubeIDBatches[batchIndex])
			channelsInBatch = append(channelsInBatch, ch)
		}
		mergedChannel := youtubedata.MergeChannels(channelsInBatch)

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
func FetchAllVideosForChannel(youtubeDB *sql.DB, youtubeChannel *youtubedata.Channel) {
	fmt.Println("The following channels are out of date: ")
	if AreVideosOutOfDate(youtubeDB, youtubeChannel) {
		fmt.Println(youtubeChannel.String())
	}
}

// AreVideosOutOfDate determines if there needs to be new videos fetched for a particular channel
func AreVideosOutOfDate(youtubeDB *sql.DB, channel *youtubedata.Channel) bool {
	// Get channel id from youtube id
	if channel.ChannelID == 0 {
		db.PopulateChannelIDFromYoutubeID(youtubeDB, channel)
	}

	// Get count from database
	dbVideoCount := db.SelectVideoCountOfChannel(youtubeDB, channel.ChannelID)
	apiVideoCount := channel.VideoCount()

	// Video counts are out of date if the count from the database
	// doesn't match the count from the api
	isOutOfDate := (dbVideoCount != apiVideoCount)
	if isOutOfDate {
		fmt.Printf("'%s' out of date. DB: %d API: %d\n", channel.Title(), dbVideoCount, apiVideoCount)
	}

	return isOutOfDate
}
