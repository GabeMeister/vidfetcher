package youtubedata

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	youtube "google.golang.org/api/youtube/v3"

	"github.com/GabeMeister/vidfetcher/api"
	"github.com/GabeMeister/vidfetcher/util"
)

// Channel contains YouTube channel data
type Channel struct {
	apiChannel *youtube.Channel
	ChannelID  int
}

func (c *Channel) String() string {
	return fmt.Sprintf("%d %s %s: %d videos", c.ChannelID, c.YoutubeID(), c.Title(), c.VideoCount())
}

// JSONString returns the json encoding of the youtubedata.Channel object
func (c *Channel) JSONString() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Fatalln("Unable to marshal youtube channel:", c.String())
	}
	return string(bytes)
}

// Title is the title of the channel
func (c *Channel) Title() string {
	if c.apiChannel == nil || c.apiChannel.Snippet == nil {
		log.Fatalln("apiChannel is nil")
	}
	return c.apiChannel.Snippet.Title
}

// YoutubeID is the channels 32 character id string recognized by Youtube
func (c *Channel) YoutubeID() string {
	if c.apiChannel == nil {
		log.Fatalln("apiChannel is nil")
	}
	return c.apiChannel.Id
}

// VideoCount is the video upload count of the channel
func (c *Channel) VideoCount() uint64 {
	if c.apiChannel == nil || c.apiChannel.Statistics == nil {
		log.Fatalln("apiChannel is nil")
	}
	return c.apiChannel.Statistics.VideoCount
}

// UploadsPlaylistID returns the youtube id of the uploads playlist for the channel
func (c *Channel) UploadsPlaylistID() string {
	if c.apiChannel == nil || c.apiChannel.ContentDetails == nil {
		log.Fatalln("channel does not have content details to fetch uploads playlist id")
	}

	return c.apiChannel.ContentDetails.RelatedPlaylists.Uploads
}

// IsChannelIDPopulated checks if this Youtube channel has had
// its channel id populated with what is in the database
func (c *Channel) IsChannelIDPopulated() bool {
	return c.ChannelID > 0
}

// ChannelsByDescendingVideoCount represents a slice of Youtube Channels
// sorted in descending order by video counts
type ChannelsByDescendingVideoCount []Channel

func (c ChannelsByDescendingVideoCount) Len() int {
	return len(c)
}

func (c ChannelsByDescendingVideoCount) Less(i, j int) bool {
	return c[i].VideoCount() > c[j].VideoCount()
}

func (c ChannelsByDescendingVideoCount) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// FetchChannelDataFromAPI - Fetches the number of uploads of a channel
// TODO: move this into api package, and change go routine strategy from fan in to bounded parallelism
func FetchChannelDataFromAPI(waitGroup *sync.WaitGroup, youtubeIDCommaText string) chan Channel {
	waitGroup.Add(1)
	defer waitGroup.Done()

	ch := make(chan Channel)

	go func() {
		service := api.GetYoutubeService()
		call := service.Channels.List("snippet,statistics,contentDetails").Id(youtubeIDCommaText)
		response, err := call.Do()
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range response.Items {
			youtubeChannel := Channel{apiChannel: item}
			ch <- youtubeChannel
		}

		close(ch)
	}()

	return ch
}

// MergeChannels merges several channels of Youtube Channels into one
func MergeChannels(channelsToMerge []chan Channel) <-chan Channel {
	out := make(chan Channel)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	var waitGrp sync.WaitGroup
	waitGrp.Add(len(channelsToMerge))

	for _, currChan := range channelsToMerge {
		go func(chanToRead <-chan Channel) {
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
func WriteYoutubeChannelsToFile(channels []Channel, filePath string) error {
	channelDataStrings := make([]string, len(channels))
	for i, channel := range channels {
		channelDataStrings[i] = channel.String()
	}

	return util.WriteLines(channelDataStrings, filePath)
}

// GetOnlyChannelsWithVideos filters out all channels in the slice that
// don't have any videos, and returns a new channel slice that contains
// only channels with videos
func GetOnlyChannelsWithVideos(channels []Channel) []Channel {
	var channelsWithVideos []Channel

	for i := range channels {
		if channels[i].VideoCount() > 0 {
			channelsWithVideos = append(channelsWithVideos, channels[i])
		}
	}

	return channelsWithVideos
}
