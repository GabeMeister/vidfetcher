package tasks

import (
	"database/sql"
	"fmt"
	"sort"
	"sync"

	youtube "google.golang.org/api/youtube/v3"

	"strings"

	"github.com/GabeMeister/vidfetcher/api"
	"github.com/GabeMeister/vidfetcher/db"
	"github.com/GabeMeister/vidfetcher/youtubedata"
)

const maxConcurrentGoRoutines = 50
const maxAPIResults = 50

// FetchYoutubeChannelInfoFromAPI fetches info from the Youtube API for the given youtube ids
func FetchYoutubeChannelInfoFromAPI(youtubeIDs []string) []youtubedata.Channel {
	var waitGroup sync.WaitGroup
	youtubeIDBatches := api.BreakYoutubeIDsIntoBatches(youtubeIDs, maxAPIResults)

	fmt.Println("api calls to make:", len(youtubeIDBatches))

	var channelsInBatch []chan youtubedata.Channel
	var count int
	var youtubeChannelData []youtubedata.Channel

	for batchStart := 0; batchStart < len(youtubeIDBatches); batchStart += maxConcurrentGoRoutines {
		channelsInBatch = nil
		batchSize := api.GetBatchSize(len(youtubeIDBatches), batchStart, maxConcurrentGoRoutines)

		for batchIndex := batchStart; batchIndex < batchStart+batchSize; batchIndex++ {
			ch := youtubedata.FetchChannelDataFromAPI(&waitGroup, youtubeIDBatches[batchIndex])
			channelsInBatch = append(channelsInBatch, ch)
		}
		mergedChannel := youtubedata.MergeChannels(channelsInBatch)

		for item := range mergedChannel {
			count++
			// fmt.Println(count, item)
			youtubeChannelData = append(youtubeChannelData, item)
		}
		waitGroup.Wait()
		fmt.Println("api calls made:", batchStart+batchSize)
	}

	return youtubeChannelData
}

// FetchNewVideosForChannels fetches any new videos for youtubeChannels, and stores them
// in youtubeDB
func FetchNewVideosForChannels(youtubeDB *sql.DB, youtubeChannels []youtubedata.Channel) {
	numGoRoutines := getChannelGoRoutineCount(youtubeChannels)

	var wg sync.WaitGroup
	wg.Add(numGoRoutines)

	// Loop to create all the go routines, all reading off of 1 channel
	ch := make(chan youtubedata.Channel)
	for i := 0; i < numGoRoutines; i++ {
		go func() {
			for {
				youtubeChannel, ok := <-ch
				if !ok {
					wg.Done()
					return
				}

				playlistItemsToFetch := FetchNewVideosForChannel(youtubeDB, &youtubeChannel)
				fmt.Printf("%s: %d vids to fetch\n", youtubeChannel.Title(), len(playlistItemsToFetch))

				// Fetch video info from playlist item
				// videosToInsert := youtubedata.FetchVideoInfo(playlistItemsToFetch)

				// Insert videos into database
				// db.InsertVideos(youtubeDB, videosToInsert)
			}
		}()
	}

	// Loop to send youtube channels down the channel
	for i := range youtubeChannels {
		ch <- youtubeChannels[i]
	}

	close(ch)
	wg.Wait()
}

// FetchNewVideosForChannel fetches all the new video uploads for the specified youtube channel
// It will stop fetching when it has fetched a video whose YoutubeID already exists in the database
func FetchNewVideosForChannel(youtubeDB *sql.DB, youtubeChannel *youtubedata.Channel) []youtubedata.PlaylistItem {
	channelID := db.SelectChannelIDFromYoutubeID(youtubeDB, youtubeChannel.YoutubeID())

	// Get a map of all youtube video ids that exist for the channel in the database
	videoMap := db.GetVideoMapForChannel(youtubeDB, channelID)

	playlistItemsToFetch := []youtubedata.PlaylistItem{}

	// Fetch videos up to 50 at a time from api
	var response *youtube.PlaylistItemListResponse

	pageToken := " "
	doneFetching := false
	for pageToken != "" && !doneFetching {
		response = youtubedata.FetchChannelUploads(youtubeChannel, strings.TrimSpace(pageToken))
		pageToken = response.NextPageToken

		for _, playlistItem := range response.Items {
			videoID := playlistItem.Snippet.ResourceId.VideoId
			_, exists := videoMap[videoID]
			if exists {
				doneFetching = true
				break
			}

			playlistItemsToFetch = append(playlistItemsToFetch, youtubedata.PlaylistItem{APIPlaylistItem: playlistItem})
		}
	}

	return playlistItemsToFetch
}

// GetOutOfDateChannels returns a slice that contains only youtube channels that are
// out of date in the database when compared to the latest video uploads
func GetOutOfDateChannels(youtubeDB *sql.DB, channels []youtubedata.Channel) []youtubedata.Channel {
	// Only get channels with videos
	channels = youtubedata.GetOnlyChannelsWithVideos(channels)

	// If channel ids haven't been initialized with database yet, then populate the channel id
	db.PopulateChannelIDsFromYoutubeIDs(youtubeDB, channels)

	// Only get channels that don't have matching video counts
	channels = db.GetOutOfDateChannels(youtubeDB, channels)

	// Sort channels by video count descending
	sortedChannels := make(youtubedata.ChannelsByDescendingVideoCount, len(channels))
	copy(sortedChannels, channels)
	sort.Sort(sortedChannels)

	return sortedChannels
}

// AreVideosOutOfDate determines if there needs to be new videos fetched for a particular channel
func AreVideosOutOfDate(youtubeDB *sql.DB, channel *youtubedata.Channel) bool {
	// Get channel id from youtube id
	if !channel.IsChannelIDPopulated() {
		db.PopulateChannelIDFromYoutubeID(youtubeDB, channel)
	}

	// Get count from database
	dbVideoCount := db.SelectVideoCountOfChannel(youtubeDB, channel.ChannelID)
	apiVideoCount := channel.VideoCount()

	// Video counts are out of date if the count from the database
	// doesn't match the count from the api
	isOutOfDate := (dbVideoCount != apiVideoCount)
	if isOutOfDate {
		fmt.Printf("'%s' out of date. DB: %d API: %d\n", channel.Title(), dbVideoCount, apiVideoCount)
	}

	return isOutOfDate
}

func getChannelGoRoutineCount(youtubeChannels []youtubedata.Channel) int {
	if len(youtubeChannels) < maxConcurrentGoRoutines {
		return len(youtubeChannels)
	}

	return maxConcurrentGoRoutines
}
