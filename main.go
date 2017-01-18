package main

import (
	"log"

	"time"

	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/tasks"
)

const outFile = "output.txt"

func main() {
	const youtubeID = "UCUj5zbH960nGLW11dBo9RRQ"
	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectAllChannelYoutubeIDs(youtubeDB)
	channels := tasks.FetchYoutubeChannelInfoFromAPI(youtubeIDs)
	channelsToFetch := tasks.GetOutOfDateChannels(youtubeDB, channels)
	log.Println(len(channelsToFetch), "are out of date")
	time.Sleep(time.Second * 3)
	tasks.FetchNewVideosForChannels(channelsToFetch)

	// channel := tasks.FetchYoutubeChannelInfoFromAPI([]string{youtubeID})[0]
	// playlistItems := tasks.FetchNewVideosForChannel(youtubeDB, &channel)

}
