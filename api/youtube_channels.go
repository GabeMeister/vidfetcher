package api

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/GabeMeister/vidfetcher/util"
	. "github.com/ahmetalpbalkan/go-linq"
)

// YoutubeChannel contains data for a YouTube channel
type YoutubeChannel struct {
	Title      string
	YoutubeID  string
	ChannelID  uint64
	VideoCount uint64
}

func (c *YoutubeChannel) String() string {
	return fmt.Sprintf("%s: %d videos", c.Title, c.VideoCount)
}

// JSONString returns the json encoding of the YoutubeChannel object
func (c *YoutubeChannel) JSONString() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Unable to marshal youtube channel: ", c.String())
	}
	return string(bytes)
}

// FetchChannelData - Fetches the number of uploads of a channel
//		channelID - id of the channel to fetch
func FetchChannelData(waitGroup *sync.WaitGroup, youtubeIDs string) chan YoutubeChannel {
	waitGroup.Add(1)
	defer waitGroup.Done()

	ch := make(chan YoutubeChannel)

	go func() {
		service := getYoutubeService()
		call := service.Channels.List("snippet,statistics").Id(youtubeIDs)
		response, err := call.Do()
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range response.Items {
			youtubeChannel := YoutubeChannel{YoutubeID: item.Id, Title: item.Snippet.Title, VideoCount: item.Statistics.VideoCount}
			ch <- youtubeChannel
		}

		close(ch)
	}()

	return ch
}

// FetchChannelUploadPlaylistIDs - Fetches Youtube video data from youtube video id
func FetchChannelUploadPlaylistIDs(waitGroup *sync.WaitGroup, youtubeIDs string) chan string {
	waitGroup.Add(1)
	defer waitGroup.Done()

	ch := make(chan string)

	go func() {
		service := getYoutubeService()
		call := service.Channels.List("contentDetails").Id(youtubeIDs)
		response, err := call.Do()
		if err != nil {
			log.Fatal(err)
		}

		if len(response.Items) > 0 {
			for _, item := range response.Items {
				ch <- item.ContentDetails.RelatedPlaylists.Uploads
			}
		}
		close(ch)
	}()

	return ch
}

// FetchLatestVideoID - Get the latest video id of channel uploads playlist
// 		playlistID - id of playlist to fetch latest video of
func FetchLatestVideoID(playlistID string) string {
	service := getYoutubeService()
	call := service.PlaylistItems.
		List("contentDetails").
		PlaylistId(playlistID).
		MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(len(response.Items))

	latestVideoID := ""
	if len(response.Items) > 0 {
		latestVideoID = response.Items[0].ContentDetails.VideoId
	}

	return latestVideoID

}

// MergeChannels merges several channels of Youtube Channels into one
func MergeChannels(channelsToMerge []chan YoutubeChannel) <-chan YoutubeChannel {
	out := make(chan YoutubeChannel)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	var waitGrp sync.WaitGroup
	waitGrp.Add(len(channelsToMerge))

	for _, currChan := range channelsToMerge {
		go func(chanToRead <-chan YoutubeChannel) {
			for n := range chanToRead {
				out <- n
			}
			waitGrp.Done()
		}(currChan)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		waitGrp.Wait()
		close(out)
	}()

	return out
}

// WriteYoutubeChannelsToFile writes a slice of youtube channel objects to a file
func WriteYoutubeChannelsToFile(channels []YoutubeChannel, filePath string) error {
	var channelDataStrings []string

	From(channels).
		Select(func(x interface{}) interface{} {
			return x.(*YoutubeChannel).String()
		}).
		ToSlice(&channelDataStrings)

	return util.WriteLines(channelDataStrings, filePath)
}

// GetYoutubeChannelsToFetch filters a slice of Youtube Channels to only include channels that
// we want to fetch
func GetYoutubeChannelsToFetch(channels []YoutubeChannel) []YoutubeChannel {
	// Youtube channels to fetch have greater than one video upload
	var sortedChannelData []YoutubeChannel
	From(channels).
		Where(func(x interface{}) bool {
			return x.(YoutubeChannel).VideoCount > 0
		}).
		OrderByDescending(func(x interface{}) interface{} {
			return x.(YoutubeChannel).VideoCount
		}).
		ToSlice(&sortedChannelData)

	return sortedChannelData
}
