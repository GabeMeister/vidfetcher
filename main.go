package main

import (
	"fmt"

	"github.com/GabeMeister/vidfetcher/api"
	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/tasks"
)

const outFile = "output.txt"

func main() {

	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectAllChannelYoutubeIDs(youtubeDB)
	youtubeChannelData := tasks.FetchAPIYoutubeChannelInfo(youtubeIDs)
	sortedChannelData := api.GetYoutubeChannelsToFetch(youtubeChannelData)

	fmt.Println("The following channels are out of date:")
	count := 0
	for _, ytChan := range sortedChannelData {
		if tasks.AreVideosOutOfDate(youtubeDB, &ytChan) {
			count++
		}
	}
	fmt.Println(count, "channels out of date.")

}
