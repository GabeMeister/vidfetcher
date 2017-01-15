package db

import (
	"database/sql"
	"log"
	"strings"
	// For postgres db
	"github.com/GabeMeister/vidfetcher/youtubedata"
)

// SelectVideoYoutubeIDsMatchingChannelID selects the YoutubeIDs of
// all videos that have channelID
func SelectVideoYoutubeIDsMatchingChannelID(youtubeDB *sql.DB, channelID int) []string {
	rows, err := youtubeDB.Query("select YoutubeID from Videos where ChannelID=$1", channelID)
	if err != nil {
		log.Fatalln(err)
	}

	youtubeIDs := []string{}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatalln(err)
		}
		youtubeIDs = append(youtubeIDs, strings.TrimSpace(id))
	}

	return youtubeIDs
}

// InsertVideos inserts a video into the database
func InsertVideos(youtubeDB *sql.DB, video []youtubedata.Video) {
	// TODO
}
