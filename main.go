package main

import (
	"log"
	"os"
	"time"

	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/tasks"
	"github.com/GabeMeister/vidfetcher/util"
	"github.com/GabeMeister/vidfetcher/youtubedata"
)

const outFile = "output.txt"

func main() {
	// fetchAllChannelsToDownload()
	go fetchAllSubscriptions()

	time.Sleep(time.Minute * 14)
	log.Println("Exiting...")
	os.Exit(0)
}

func fetchAllChannels() {
	logFile := util.InitLoggerAndLogFile("all_channels.log")
	defer logFile.Close()

	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectAllChannelYoutubeIDs(youtubeDB)
	channels := tasks.FetchYoutubeChannelInfoFromAPI(youtubeIDs)
	channelsToFetch := tasks.GetOutOfDateChannels(youtubeDB, channels)

	log.Println(len(channelsToFetch), "are out of date")
	tasks.FetchNewVideosForChannels(channelsToFetch)
}

func fetchAllSubscriptions() {
	logFile := util.InitLoggerAndLogFile("subscriptions.log")
	defer logFile.Close()

	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectSubscriptionChannelYoutubeIDs(youtubeDB)
	channels := tasks.FetchYoutubeChannelInfoFromAPI(youtubeIDs)
	db.PopulateChannelIDsFromYoutubeIDs(youtubeDB, channels)

	log.Println(len(channels), "subscription channels are out of date")
	tasks.FetchNewVideosForChannels(channels)
}

func fetchAllChannelsToDownload() {
	logFile := util.InitLoggerAndLogFile("channels_to_download.log")
	defer logFile.Close()

	youtubeDB := db.CreateDBInstance()
	defer youtubeDB.Close()

	youtubeIDs := db.SelectChannelsToDownloadYoutubeIDs(youtubeDB)
	channels := tasks.FetchYoutubeChannelInfoFromAPI(youtubeIDs)
	db.PopulateChannelIDsFromYoutubeIDs(youtubeDB, channels)

	log.Println(len(channels), "channels to download are out of date")
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
