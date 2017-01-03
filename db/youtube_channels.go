package db

import (
	"database/sql"
	"log"

	"github.com/GabeMeister/vidfetcher/youtubedata"
	// For postgres db
	"github.com/lib/pq"
)

// PopulateChannelIDsFromYoutubeIDs selects the corresponding channel ids from a
// slice of youtube ids
func PopulateChannelIDsFromYoutubeIDs(youtubeDB *sql.DB, channels []youtubedata.Channel) {
	for i := range channels {
		if !channels[i].IsChannelIDPopulated() {
			channels[i].ChannelID = SelectChannelIDFromYoutubeID(youtubeDB, channels[i].YoutubeID())
		}
	}
}

// PopulateChannelIDFromYoutubeID sets the channel ID of a Channel from a Youtube ID
func PopulateChannelIDFromYoutubeID(youtubeDB *sql.DB, channel *youtubedata.Channel) {
	if channel.YoutubeID() == "" {
		log.Fatal("Youtube ID cannot be blank")
	}

	rows, err := youtubeDB.Query("select ChannelID from Channels where YoutubeID=$1", channel.YoutubeID())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rows.Next()

	if err = rows.Scan(&channel.ChannelID); err != nil {
		log.Fatal(err)
	}
}

// SelectAllChannelYoutubeIDs fetches all channel
func SelectAllChannelYoutubeIDs(youtubeDB *sql.DB) []string {
	// TODO: Remove the 250 video count limit
	return SelectColumnFromTable(youtubeDB, "YoutubeID", "Channels", 50)
}

// SelectChannelIDFromYoutubeID selects the corresponding channel id from
// a given youtube id
func SelectChannelIDFromYoutubeID(youtubeDB *sql.DB, youtubeID string) int {
	if youtubeID == "" {
		log.Fatal("Youtube ID cannot be blank")
	}

	rows, err := youtubeDB.Query("select ChannelID from Channels where YoutubeID=$1", youtubeID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rows.Next()

	var channelID int
	if err = rows.Scan(&channelID); err != nil {
		log.Fatal(err)
	}

	return channelID
}

// SelectVideoCountOfChannel gets the count of video uploads for a channel
func SelectVideoCountOfChannel(youtubeDB *sql.DB, channelID int) uint64 {
	rows, err := youtubeDB.Query(`select VideoCount from Channels where ChannelID=$1;`, channelID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count uint64
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}

	return count
}

// GetOutOfDateChannels returns only channels that are out of date in the database
func GetOutOfDateChannels(youtubeDB *sql.DB, channels []youtubedata.Channel) []youtubedata.Channel {
	var outOfDateChannels []youtubedata.Channel

	videoCountMap := getChannelIDToVideoCountMap(youtubeDB, channels)
	for i := range channels {
		if channels[i].VideoCount() != videoCountMap[channels[i].ChannelID] {
			outOfDateChannels = append(outOfDateChannels, channels[i])
		}
	}

	return outOfDateChannels
}

// getChannelIDToVideoCountMap creates a map of channel ids to video counts for the specified channels
func getChannelIDToVideoCountMap(youtubeDB *sql.DB, channels []youtubedata.Channel) map[int]uint64 {
	channelIDs := make([]int, len(channels))
	for i := range channels {
		channelIDs[i] = channels[i].ChannelID
	}

	rows, err := youtubeDB.Query("select ChannelID,VideoCount from Channels where ChannelID = ANY($1)", pq.Array(channelIDs))
	if err != nil {
		log.Fatal(err)
	}

	videoCountMap := make(map[int]uint64)
	for rows.Next() {
		var channelID int
		var videoCount uint64

		err := rows.Scan(&channelID, &videoCount)
		if err != nil {
			log.Fatal(err)
		}
		videoCountMap[channelID] = videoCount
	}

	return videoCountMap
}
