package main

import (
	"fmt"

	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/tasks"
)

const outFile = "output.txt"

func main() {
	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectAllChannelYoutubeIDs(youtubeDB)
	channels := tasks.FetchYoutubeChannelInfoFromAPI(youtubeIDs)
	channelsToFetch := tasks.GetOutOfDateChannels(youtubeDB, channels)

	fmt.Println(len(channelsToFetch), "are out of date:")
	for _, channel := range channelsToFetch {
		fmt.Println(channel.String())
	}

}
