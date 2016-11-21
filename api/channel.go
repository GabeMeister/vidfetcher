package api

import (
	"fmt"
)

// Channel contains data for a YouTube channel
type Channel struct {
	Title      string
	YoutubeID  string
	ChannelID  uint64
	VideoCount uint64
}

func (channel Channel) String() string {
	return fmt.Sprintf("%s, VideoCount: %d", channel.Title, channel.VideoCount)
}

// ChannelsByVideoCount is a collections of channels that can be sorted by the video count
type ChannelsByVideoCount []Channel

// Len gets the length of the channels array
func (channels ChannelsByVideoCount) Len() int {
	return len(channels)
}

func (channels ChannelsByVideoCount) Less(i, j int) bool {
	result := false
	if channels[i].VideoCount < channels[j].VideoCount {
		result = true
	}
	return result
}

func (channels ChannelsByVideoCount) Swap(i, j int) {
	channels[i], channels[j] = channels[j], channels[i]
}
