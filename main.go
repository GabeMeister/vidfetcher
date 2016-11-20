package main

import (
	"fmt"
	"sync"
)

var waitGroup sync.WaitGroup

const maxGoProcs = 25
const maxBatchSize = 50

func main() {
	db := CreateDBInstance()
	defer db.Close()
	youtubeIDs := SelectColumnFromChannels(db, "YoutubeID", "Channels")

	youtubeIDBatchArr := BreakYoutubeIDsIntoBatches(youtubeIDs, maxBatchSize)

	var allChannels []chan string
	var count int
	for poolIndex := 0; poolIndex < len(youtubeIDBatchArr); poolIndex += maxGoProcs {
		allChannels = nil
		batchSize := GetBatchSize(len(youtubeIDBatchArr), poolIndex, maxGoProcs)

		for batchIndex := poolIndex; batchIndex < poolIndex+batchSize; batchIndex++ {
			ch := FetchChannelUploadCount(&waitGroup, youtubeIDBatchArr[batchIndex])
			allChannels = append(allChannels, ch)
		}
		mergedChannel := merge(allChannels)

		for item := range mergedChannel {
			count++
			fmt.Println(count, ".)", item)
		}

		waitGroup.Wait()

	}

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
