package main

import (
	"fmt"

	"github.com/GabeMeister/vidfetcher/api"
	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/tasks"
	. "github.com/ahmetalpbalkan/go-linq"
)

const outFile = "output.txt"

func main() {

	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectAllChannelYoutubeIDs(youtubeDB)
	youtubeChannelData := tasks.FetchAPIYoutubeChannelInfo(youtubeIDs)

	var sortedChannelData []api.YoutubeChannel
	From(youtubeChannelData).
		Where(func(x interface{}) bool {
			return x.(api.YoutubeChannel).VideoCount > 0
		}).
		OrderByDescending(func(x interface{}) interface{} {
			return x.(api.YoutubeChannel).VideoCount
		}).
		ToSlice(&sortedChannelData)

	fmt.Println("The following channels are out of date: ")
	count := 0
	for _, ytChan := range sortedChannelData {
		if tasks.AreVideosOutOfDate(youtubeDB, &ytChan) {
			count++
		}
	}
	fmt.Println(count, "channels out of date.")

}
