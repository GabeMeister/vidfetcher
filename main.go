package main

import (
	"fmt"
	"sync"

	"github.com/GabeMeister/vidfetcher/api"
	. "github.com/ahmetalpbalkan/go-linq"
)

var waitGroup sync.WaitGroup

const maxGoProcs = 25
const maxBatchSize = 50
const outputFilePath = "output.txt"

func main() {
	db := CreateDBInstance()
	defer db.Close()

	youtubeIDs := SelectColumnFromChannels(db, "YoutubeID", "Channels")
	youtubeIDBatchArr := BreakYoutubeIDsIntoBatches(youtubeIDs, maxBatchSize)

	fmt.Println("Amount of api calls to make:", len(youtubeIDBatchArr))

	var channelBatch []chan api.Channel
	var count int
	var youtubeChannelData []api.Channel

	for poolIndex := 0; poolIndex < len(youtubeIDBatchArr); poolIndex += maxGoProcs {
		channelBatch = nil
		batchSize := GetBatchSize(len(youtubeIDBatchArr), poolIndex, maxGoProcs)

		for batchIndex := poolIndex; batchIndex < poolIndex+batchSize; batchIndex++ {
			ch := FetchChannelUploadCount(&waitGroup, youtubeIDBatchArr[batchIndex])
			channelBatch = append(channelBatch, ch)
		}
		mergedChannel := merge(channelBatch)

		for item := range mergedChannel {
			count++
			// fmt.Println(count, item)
			youtubeChannelData = append(youtubeChannelData, item)
		}
		waitGroup.Wait()
		fmt.Println("Api Calls Made:", poolIndex+maxGoProcs)
	}

	var sortedChannelData []string
	From(youtubeChannelData).
		Where(func(x interface{}) bool { return x.(api.Channel).VideoCount > 0 }).
		OrderByDescending(func(x interface{}) interface{} { return x.(api.Channel).VideoCount }).
		Select(func(x interface{}) interface{} { return x.(api.Channel).String() }).
		ToSlice(&sortedChannelData)

	WriteLines(sortedChannelData, outputFilePath)

}

func doDatabaseStuff() {
	// db := CreateDBInstance()
	// defer db.Close()

	// Select all youtube ids of channels in collections
	// youtubeIDs := SelectUniqueCollectionItemyoutubeIDs(db)
	// channelIDs := SelectChannelIdsFromyoutubeIDs(db, youtubeIDs)

	// Select latest video youtube id and video count of each channel from database
	// for _, id := range channelIDs {
	// 	fmt.Println(id)
	// }

	// Select latest video youtube id and video count of each channel from api

	// Determine which channels need to be updated

	// Update all channels
}
