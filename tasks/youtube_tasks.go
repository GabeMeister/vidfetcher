package tasks

import (
	"database/sql"
	"fmt"
	"sort"
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

	fmt.Println("api calls to make:", len(youtubeIDBatches))

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

// FetchAllVideosForChannel fetches all the video uploads for the specified youtube channel
func FetchAllVideosForChannel(youtubeDB *sql.DB, youtubeChannel *youtubedata.Channel) {
	// TODO
}

// GetOutOfDateChannels returns a slice that contains only youtube channels that are
// out of date in the database when compared to the latest video uploads
func GetOutOfDateChannels(youtubeDB *sql.DB, channels []youtubedata.Channel) []youtubedata.Channel {
	// Only get channels with videos
	channels = youtubedata.GetOnlyChannelsWithVideos(channels)

	// If channel ids haven't been initialized with database yet, then populate the channel id
	db.PopulateChannelIDsFromYoutubeIDs(youtubeDB, channels)

	// Only get channels that don't have matching video counts
	channels = db.GetOutOfDateChannels(youtubeDB, channels)

	// Sort channels by video count descending
	sortedChannels := make(youtubedata.ChannelsByDescendingVideoCount, len(channels))
	copy(sortedChannels, channels)
	sort.Sort(sortedChannels)

	return sortedChannels
}

// AreVideosOutOfDate determines if there needs to be new videos fetched for a particular channel
func AreVideosOutOfDate(youtubeDB *sql.DB, channel *youtubedata.Channel) bool {
	// Get channel id from youtube id
	if !channel.IsChannelIDPopulated() {
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
