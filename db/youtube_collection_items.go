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
		log.Fatal(err)
	}
	defer rows.Close()

	var channelIDs []string
	for rows.Next() {
		var channelIDStr string
		err := rows.Scan(&channelIDStr)
		if err != nil {
			log.Fatal(err)
		}

		channelIDs = append(channelIDs, strings.TrimSpace(channelIDStr))
	}

	return channelIDs
}
