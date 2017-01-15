package db

import (
	"database/sql"
	"log"
	"strings"

	"github.com/GabeMeister/vidfetcher/youtubedata"
	// For postgres db
	"github.com/lib/pq"
)

// PopulateChannelIDsFromYoutubeIDs selects the corresponding channel ids from a
// slice of youtube ids
func PopulateChannelIDsFromYoutubeIDs(youtubeDB *sql.DB, channels []youtubedata.Channel) {
	youtubeIDs := make([]string, len(channels))
	for i := range channels {
		youtubeIDs[i] = channels[i].YoutubeID()
	}

	channelIDMap := SelectChannelIDsFromYoutubeIDs(youtubeDB, youtubeIDs)
	for i := range channels {
		channels[i].ChannelID = channelIDMap[channels[i].YoutubeID()]
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
	return SelectColumnFromTable(youtubeDB, "YoutubeID", "Channels", 100)
}

// SelectChannelIDsFromYoutubeIDs does a batch select of channel ids for the given youtube channels
func SelectChannelIDsFromYoutubeIDs(youtubeDB *sql.DB, youtubeIDs []string) map[string]int {
	rows, err := youtubeDB.Query(`select ChannelID, YoutubeID from Channels where YoutubeID = ANY($1);`, pq.Array(youtubeIDs))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	channelIDMap := make(map[string]int)
	for rows.Next() {
		var channelID int
		var youtubeID string

		err := rows.Scan(&channelID, &youtubeID)
		if err != nil {
			log.Fatal(err)
		}

		channelIDMap[strings.TrimSpace(youtubeID)] = channelID
	}

	return channelIDMap
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

// SelectChannelUploadsPlaylistID selects the uploads playlist id of the channel
// with id of channelID
func SelectChannelUploadsPlaylistID(youtubeDB *sql.DB, channelID int) string {
	rows, err := youtubeDB.Query("select UploadPlaylist from Channels where ChannelID=$1", channelID)
	if err != nil {
		log.Fatal(err)
	}

	uploadsPlaylistID := ""
	if rows.Next() {
		rows.Scan(&uploadsPlaylistID)
	}

	return uploadsPlaylistID
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

// GetVideoMapForChannel returns a map of video ids for the specified channel Youtube ID
func GetVideoMapForChannel(youtubeDB *sql.DB, channelID int) map[string]bool {
	youtubeIDs := SelectVideoYoutubeIDsMatchingChannelID(youtubeDB, channelID)

	youtubeIDMap := make(map[string]bool)
	for _, id := range youtubeIDs {
		youtubeIDMap[id] = true
	}

	return youtubeIDMap
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
	defer rows.Close()

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
