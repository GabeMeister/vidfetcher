package api

import (
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

// MaxAPIGoRoutines is the maximum allowable number of go routines
// that are allowed while fetching data from the youtube api
const MaxAPIGoRoutines = 50

// MaxAPIResults is the maximum amount of results allowed to ask for in
// a call to the Youtube API
const MaxAPIResults = 50

// DeveloperKey is the id used to make youtube api calls
const DeveloperKey = "AIzaSyC9uXxwF4PxYilaOvPTDLdXAnToBwFvXcs"

// BreakYoutubeIDsIntoBatches breaks a list of ids into smaller batches
func BreakYoutubeIDsIntoBatches(youtubeIDs []string, batchSize int) (batchArr []string) {
	for idIndex := 0; idIndex < len(youtubeIDs); idIndex += batchSize {
		batchSize := GetBatchSize(len(youtubeIDs), idIndex, batchSize)
		batchArr = append(batchArr, strings.Join(youtubeIDs[idIndex:idIndex+batchSize], ","))
	}
	return batchArr
}

// GetBatchSize determines how big a batch can be without going over the bounds of the array
func GetBatchSize(arrSize int, index int, maxBatchSize int) int {
	batchSize := maxBatchSize
	if index+maxBatchSize > arrSize {
		batchSize = arrSize - index
	}
	return batchSize
}

// GetYoutubeService returns a youtube service to make api calls to
func GetYoutubeService() *youtube.Service {
	client := getYoutubeClient()
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalln("error creating new YouTube client", err)
	}

	return service
}

func getYoutubeClient() *http.Client {
	return &http.Client{
		Transport: &transport.APIKey{Key: DeveloperKey},
	}
}
