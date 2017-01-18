package youtubedata

import (
	"encoding/json"
	"fmt"
	"log"

	"time"

	"strings"

	youtube "google.golang.org/api/youtube/v3"
)

// Video contains Youtube video data
type Video struct {
	APIVideo  *youtube.Video
	VideoID   int
	ChannelID int
}

func (v *Video) String() string {
	return fmt.Sprintf("YoutubeID:%s ChannelID:%d Title:%s Thumbnail:%s Duration:%s ViewCount:%d PublishedAt:%s",
		v.YoutubeID(), v.ChannelID, v.Title(), v.Thumbnail(), v.Duration(), v.ViewCount(), v.PublishedAt())
}

// JSONString returns the json encoding of the youtubedata.Video object
func (v *Video) JSONString() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Println("Unable to marshal youtube video", v.String(), err)
		return ""
	}
	return string(bytes)
}

// Title is the title of the video
func (v *Video) Title() string {
	if !v.snippetExists() {
		log.Fatalln("APIVideo's snippet is nil, cannot access video title")
	}
	title := v.APIVideo.Snippet.Title
	if len(title) > 150 {
		title = title[:150]
	}
	return title
}

// YoutubeID is the video's 32 character id string recognized by Youtube
func (v *Video) YoutubeID() string {
	if !v.snippetExists() {
		log.Fatalln("APIVideo's snippet is nil, cannot access YoutubeID")
	}
	return v.APIVideo.Id
}

// Thumbnail is the medium thumbnail url of the video
func (v *Video) Thumbnail() string {
	if !v.snippetExists() {
		log.Fatalln("APIVideo's snippet is nil, cannot access Thumbnail")
	}
	return v.APIVideo.Snippet.Thumbnails.Medium.Url
}

// Duration is the time duration of the youtube video
// It is retrieved by a string and converted
func (v *Video) Duration() string {
	if !v.contentDetailsExist() {
		log.Fatalln("APIVideo's content details is nil, cannot access Duration")
	}

	durationStr := strings.TrimLeft(v.APIVideo.ContentDetails.Duration, "PT")
	durationStr = strings.ToLower(durationStr)
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Println("unable to parse duration string", err)
		// If unable to parse duration, we just pretend the video is 2 minutes long
		duration = time.Duration(2) * time.Minute
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - (hours * 60)
	seconds := int(duration.Seconds()) - (hours * 60) - (minutes * 60)

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// ViewCount is the amount of views the Youtube Video has
func (v *Video) ViewCount() uint64 {
	if !v.statisticsExist() {
		log.Fatalln("APIVideo's statistics is nil, cannot access view count")
	}
	return v.APIVideo.Statistics.ViewCount
}

// PublishedAt is the time when the video was published
func (v *Video) PublishedAt() string {
	if !v.snippetExists() {
		log.Fatalln("APIVideo's snippet is nil, cannot access published at")
	}

	vidTime, err := time.Parse(time.RFC3339, v.APIVideo.Snippet.PublishedAt)
	if err != nil {
		log.Println("could not parse time", err)
		// If unable to put in time, just use today's time
		vidTime = time.Now()
	}

	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		vidTime.Year(), vidTime.Month(), vidTime.Day(), vidTime.Hour(), vidTime.Minute(), vidTime.Second())
}

func (v *Video) snippetExists() bool {
	return v.APIVideo != nil && v.APIVideo.Snippet != nil
}

func (v *Video) contentDetailsExist() bool {
	return v.APIVideo != nil && v.APIVideo.ContentDetails != nil
}

func (v *Video) statisticsExist() bool {
	return v.APIVideo != nil && v.APIVideo.Statistics != nil
}
