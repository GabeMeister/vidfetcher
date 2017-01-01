package db

import (
	"database/sql"
	"log"

	// For postgres db

	_ "github.com/lib/pq"
)

// SelectAllChannelYoutubeIDs fetches all channel
func SelectAllChannelYoutubeIDs(youtubeDB *sql.DB) []string {
	// TODO: Remove the 250 video count limit
	return SelectColumnFromTable(youtubeDB, "YoutubeID", "Channels", 50)
}

// SelectChannelIdsFromYoutubeIDs selects the corresponding channel ids from a
// slice of youtube ids
func SelectChannelIdsFromYoutubeIDs(youtubeDB *sql.DB, youtubeIDs []string) []int {
	var channelIDs []int
	for _, youtubeID := range youtubeIDs {
		channelIDs = append(channelIDs, SelectChannelIDFromYoutubeID(youtubeDB, youtubeID))
	}

	return channelIDs
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
	rows, err := youtubeDB.Query(`select count(*) from videos where channelid=$1;`, channelID)
	if err != nil {
		log.Println("Incorrect sql")
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
