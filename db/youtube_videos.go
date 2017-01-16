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

// InsertVideos inserts a slice of videos into the database
func InsertVideos(youtubeDB *sql.DB, videos []youtubedata.Video) {
	for _, vid := range videos {
		InsertVideo(youtubeDB, vid)
	}
}

// InsertVideo inserts 1 video into the database
func InsertVideo(youtubeDB *sql.DB, video youtubedata.Video) {
	log.Printf("inserting %s %s\n", video.YoutubeID(), video.Title())

	stmt, err := youtubeDB.Prepare("insert into Videos (YoutubeID,ChannelID,Title,Thumbnail,Duration,ViewCount,PublishedAt) values ($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		log.Fatalln(err)
	}

	result, err := stmt.Exec(video.YoutubeID(), video.ChannelID, video.Title(), video.Thumbnail(), video.Duration(), video.ViewCount(), video.PublishedAt())
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}

	if rowsAffected < 1 {
		log.Fatalln("less than 1 row affected for %s %s", video.YoutubeID(), video.Title)
	}

}
