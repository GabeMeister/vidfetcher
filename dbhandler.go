package main

import (
	"database/sql"
	"log"
	"strings"

	"fmt"

	_ "github.com/lib/pq"
)

// SelectColumnFromChannels fetches all the channel ids in the channels table
func SelectColumnFromChannels(db *sql.DB, column string, table string) (results []string) {
	sql := fmt.Sprintf("select %s from %s;", column, table)
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var colValue string
		err := rows.Scan(&colValue)
		if err != nil {
			panic(err)
		}
		results = append(results, strings.TrimSpace(colValue))
	}

	return results
}

// SelectChannelIdsFromYoutubeIDs selects the corresponding channel ids from a
// slice of youtube ids
func SelectChannelIdsFromYoutubeIDs(db *sql.DB, youtubeIDs []string) []int {
	var channelIDs []int
	for _, youtubeID := range youtubeIDs {
		channelIDs = append(channelIDs, SelectChannelIDFromYoutubeID(db, youtubeID))
	}

	return channelIDs
}

// SelectChannelIDFromYoutubeID selects the corresponding channel id from
// a given youtube id
func SelectChannelIDFromYoutubeID(db *sql.DB, youtubeID string) int {
	rows, err := db.Query("select channelid from channels where youtubeid = $1", youtubeID)
	if err != nil {
		log.Fatal(err)
	}

	rows.Next()

	var channelID int
	err = rows.Scan(&channelID)
	if err != nil {
		log.Fatal(err)
	}

	return channelID
}

// SelectUniqueCollectionItemYoutubeIDs selects all channel ids from collections
func SelectUniqueCollectionItemYoutubeIDs(db *sql.DB) []string {

	// Joined Tables:
	// - Channels
	// - CollectionItems
	rows, err := db.Query(`select ch.YoutubeID 
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

// SelectVideoCountOfChannel gets the count of video uploads for a channel
func SelectVideoCountOfChannel(db *sql.DB, channelID int) int {
	rows, err := db.Query(`select count(*) from videos where channelid=?;`, channelID)
	if err != nil {
		log.Println("Incorrect sql")
		log.Fatal(err)
	}

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}

	return count
}

// SelectChannelsFromDb selects all channels from database
func SelectChannelsFromDb() {
	dbinfo := dbConnectionStr()
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select youtubeid from channels")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var channelYoutubeID string
	for rows.Next() {
		err := rows.Scan(&channelYoutubeID)
		if err != nil {
			log.Fatal()
		}
		log.Println(channelYoutubeID)
	}
}

// func DoesVideoExist(youtubeID string) {

// }
