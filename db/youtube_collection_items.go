package db

import (
	"database/sql"
	"log"
	"strings"

	// For postgres db
	_ "github.com/lib/pq"
)

// SelectAllCollectionItemYoutubeIDs selects all channel ids from collections
func SelectAllCollectionItemYoutubeIDs(youtubeDB *sql.DB) []string {
	channelIDs := []string{}

	rows, err := youtubeDB.Query(`select ch.YoutubeID 
							from Channels ch
							where ch.ChannelID in
							(
								select ci.ItemChannelID 
								from CollectionItems ci 
								group by ci.ItemChannelID
							)
							order by ch.YoutubeID;`)
	if err != nil {
		log.Println("unable to select collection item youtube ids", err)
		return channelIDs
	}
	defer rows.Close()

	for rows.Next() {
		var youtubeID string
		err := rows.Scan(&youtubeID)
		if err != nil {
			log.Println("unable to scan collection item youtube ids", err)
			continue
		}

		channelIDs = append(channelIDs, strings.TrimSpace(youtubeID))
	}

	return channelIDs
}
