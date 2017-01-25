package main

import (
	"log"
	"time"

	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/tasks"
	"github.com/GabeMeister/vidfetcher/youtubedata"
)

const outFile = "output.txt"

func main() {
	fetchAllChannelsToDownload()
}

func fetchAllChannels() {
	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectAllChannelYoutubeIDs(youtubeDB)
	channels := tasks.FetchYoutubeChannelInfoFromAPI(youtubeIDs)
	channelsToFetch := tasks.GetOutOfDateChannels(youtubeDB, channels)

	log.Println(len(channelsToFetch), "are out of date")
	time.Sleep(time.Second * 3)
	tasks.FetchNewVideosForChannels(channelsToFetch)
}

func fetchAllChannelsToDownload() {
	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectChannelsToDownloadYoutubeIDs(youtubeDB)
	channels := tasks.FetchYoutubeChannelInfoFromAPI(youtubeIDs)
	db.PopulateChannelIDsFromYoutubeIDs(youtubeDB, channels)

	log.Println(len(channels), "channels to download are out of date")
	time.Sleep(time.Second * 3)
	tasks.FetchNewVideosForChannelsToDownload(channels)
}

func fetchSingleChannel() {
	const youtubeID = "UCqG4BLgo8VyK4LK8tXRqDVg"

	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	channel := tasks.FetchYoutubeChannelInfoFromAPI([]string{youtubeID})[0]
	db.PopulateChannelIDFromYoutubeID(youtubeDB, &channel)
	tasks.FetchNewVideosForChannels([]youtubedata.Channel{channel})
}
