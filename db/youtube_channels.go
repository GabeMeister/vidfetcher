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
		log.Fatalln("Youtube ID cannot be blank while selecting channel ID from YoutubeID")
	}

	rows, err := youtubeDB.Query("select ChannelID from Channels where YoutubeID=$1", channel.YoutubeID())
	if err != nil {
		log.Printf("unable to select channel id where YoutubeID=%s %v\n", channel.YoutubeID(), err)
		return
	}
	defer rows.Close()

	rows.Next()

	if err = rows.Scan(&channel.ChannelID); err != nil {
		log.Printf("unable to scan channel id where YoutubeID=%s %v\n", channel.YoutubeID(), err)
	}
}

// SelectAllChannelYoutubeIDs fetches all channel
func SelectAllChannelYoutubeIDs(youtubeDB *sql.DB) []string {
	return SelectColumnFromTable(youtubeDB, "YoutubeID", "Channels", 0)
}

// SelectChannelsToDownloadYoutubeIDs gets the youtube ids of all channels found
// in the ChannelsToDownload table
func SelectChannelsToDownloadYoutubeIDs(youtubeDB *sql.DB) []string {
	youtubeIDs := []string{}

	rows, err := youtubeDB.Query(`
		select ch.YoutubeID 
		from ChannelsToDownload ctd
		inner join Channels ch on ch.ChannelID=ctd.ChannelID;`)
	if err != nil {
		log.Println("unable to select youtube ids from ChannelsToDownload table", err)
		return youtubeIDs
	}
	defer rows.Close()

	for rows.Next() {
		var youtubeID string
		err := rows.Scan(&youtubeID)
		if err != nil {
			log.Println("unable to scan youtube ID from ChannelsToDownload", err)
			continue
		}

		youtubeIDs = append(youtubeIDs, strings.TrimSpace(youtubeID))
	}

	return youtubeIDs
}

// SelectSubscriptionChannelYoutubeIDs returns a list of youtube ids from
// the channels in the Subscriptions table.
func SelectSubscriptionChannelYoutubeIDs(youtubeDB *sql.DB) []string {
	youtubeIDs := []string{}

	rows, err := youtubeDB.Query(`
		select ch.YoutubeID 
		from Subscriptions s
		inner join Channels ch on ch.ChannelID=s.BeingSubscribedToChannelID;`)
	if err != nil {
		log.Println("unable to select youtube ids from ChannelsToDownload table", err)
		return youtubeIDs
	}
	defer rows.Close()

	for rows.Next() {
		var youtubeID string
		err := rows.Scan(&youtubeID)
		if err != nil {
			log.Println("unable to scan youtube ID from ChannelsToDownload", err)
			continue
		}

		youtubeIDs = append(youtubeIDs, strings.TrimSpace(youtubeID))
	}

	return youtubeIDs
}

// SelectChannelIDsFromYoutubeIDs does a batch select of channel ids for the given youtube channels
func SelectChannelIDsFromYoutubeIDs(youtubeDB *sql.DB, youtubeIDs []string) map[string]int {
	channelIDMap := make(map[string]int)

	rows, err := youtubeDB.Query(`select ChannelID, YoutubeID from Channels where YoutubeID = ANY($1);`, pq.Array(youtubeIDs))
	if err != nil {
		log.Println("unable to select channel ids from youtube ids", err)
		return channelIDMap
	}
	defer rows.Close()

	for rows.Next() {
		var channelID int
		var youtubeID string

		err := rows.Scan(&channelID, &youtubeID)
		if err != nil {
			log.Println("unable to scan channel id and youtube id from youtube ids", err)
			continue
		}

		channelIDMap[strings.TrimSpace(youtubeID)] = channelID
	}

	return channelIDMap
}

// SelectChannelIDFromYoutubeID selects the corresponding channel id from
// a given youtube id
func SelectChannelIDFromYoutubeID(youtubeDB *sql.DB, youtubeID string) int {
	if youtubeID == "" {
		log.Fatalln("Youtube ID cannot be blank while selecting channel id from Youtube ID")
	}

	channelID := 0

	rows, err := youtubeDB.Query("select ChannelID from Channels where YoutubeID=$1", youtubeID)
	if err != nil {
		log.Println("unable to select channel id from youtube id", err)
		return channelID
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&channelID); err != nil {
			log.Println("unable to scan channel id from youtube id", err)
		}
	}

	return channelID
}

// SelectChannelUploadsPlaylistID selects the uploads playlist id of the channel
// with id of channelID
func SelectChannelUploadsPlaylistID(youtubeDB *sql.DB, channelID int) string {
	uploadsPlaylistID := ""

	rows, err := youtubeDB.Query("select UploadPlaylist from Channels where ChannelID=$1", channelID)
	if err != nil {
		log.Printf("unable to select UploadPlaylist for channel: %d %v\n", channelID, err)
		return uploadsPlaylistID
	}

	if rows.Next() {
		if err := rows.Scan(&uploadsPlaylistID); err != nil {
			log.Printf("unable to scan UploadPlaylist for channel: %d %v\n", channelID, err)
		}
	}

	return uploadsPlaylistID
}

// SelectVideoCountOfChannel gets the count of video uploads for a channel
func SelectVideoCountOfChannel(youtubeDB *sql.DB, channelID int) uint64 {
	var count uint64

	rows, err := youtubeDB.Query(`select VideoCount from Channels where ChannelID=$1;`, channelID)
	if err != nil {
		log.Printf("unable to select video count for channel id = %d, %v\n", channelID, err)
		return count
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Printf("unable to scan video count for channel id = %d, %v\n", channelID, err)
		}
	}

	return count
}

// DeleteChannelFromChannelsToDownload deletes the specified channel id from the
// ChannelsToDownload table
func DeleteChannelFromChannelsToDownload(youtubeDB *sql.DB, channelID int) {
	stmt, err := youtubeDB.Prepare(`delete from ChannelsToDownload where ChannelID=$1;`)
	if err != nil {
		log.Println("unable to prepare delete for channel from ChannelsToDownload", err)
		return
	}

	if _, err := stmt.Exec(channelID); err != nil {
		log.Println("unable to prepare delete for channel from ChannelsToDownload", err)
		return
	}
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
	videoCountMap := make(map[int]uint64)

	// Convert slice of channels to slice of channel ids
	channelIDs := make([]int, len(channels))
	for i := range channels {
		channelIDs[i] = channels[i].ChannelID
	}

	rows, err := youtubeDB.Query("select ChannelID,VideoCount from Channels where ChannelID = ANY($1)", pq.Array(channelIDs))
	if err != nil {
		log.Println("unable to select ChannelID, VideoCount for channel ids", err)
		return videoCountMap
	}
	defer rows.Close()

	for rows.Next() {
		var channelID int
		var videoCount uint64

		err := rows.Scan(&channelID, &videoCount)
		if err != nil {
			log.Println("unable to select ChannelID, VideoCount for channel ids", err)
			continue
		}
		videoCountMap[channelID] = videoCount
	}

	return videoCountMap
}
