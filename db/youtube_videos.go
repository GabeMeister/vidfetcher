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
	youtubeIDs := []string{}

	rows, err := youtubeDB.Query("select YoutubeID from Videos where ChannelID=$1", channelID)
	if err != nil {
		log.Println("unable to select video youtube ids matching channel id", err)
		return youtubeIDs
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Println("unable to scan video youtube id", err)
			continue
		}
		youtubeIDs = append(youtubeIDs, strings.TrimSpace(id))
	}

	return youtubeIDs
}

// InsertVideos inserts a slice of videos into the database
func InsertVideos(youtubeDB *sql.DB, videos []youtubedata.Video) {
	for _, vid := range videos {
		InsertVideo(youtubeDB, vid)
	}
}

// InsertVideo inserts 1 video into the database
func InsertVideo(youtubeDB *sql.DB, video youtubedata.Video) {
	// log.Printf("inserting %s %s\n", video.YoutubeID(), video.Title())

	stmt, err := youtubeDB.Prepare("insert into Videos (YoutubeID,ChannelID,Title,Thumbnail,Duration,ViewCount,PublishedAt) values ($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		log.Printf("Skipped inserting %s, %v\n", video.Title(), err)
		return
	}

	result, err := stmt.Exec(video.YoutubeID(), video.ChannelID, video.Title(), video.Thumbnail(), video.Duration(), video.ViewCount(), video.PublishedAt())
	if err != nil {
		log.Printf("Skipped inserting %s, %v\n", video.Title(), err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Skipped inserting %s, %v\n", video.Title(), err)
		return
	}

	if rowsAffected < 1 {
		log.Printf("Rows affected < 1 for %s\n", video.Title())
		return
	}

}
