package main

import (
	"fmt"

	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/tasks"
	"github.com/GabeMeister/vidfetcher/youtubedata"
)

const outFile = "output.txt"

func main() {
	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectAllChannelYoutubeIDs(youtubeDB)
	youtubeChannelData := tasks.FetchYoutubeChannelInfoFromAPI(youtubeIDs)
	sortedChannelData := youtubedata.GetYoutubeChannelsToFetch(youtubeChannelData)

	fmt.Println("The following channels are out of date:")
	count := 0
	for _, ytChan := range sortedChannelData {
		if tasks.AreVideosOutOfDate(youtubeDB, &ytChan) {
			count++
		}
	}
	fmt.Println(count, "channels out of date.")

}
