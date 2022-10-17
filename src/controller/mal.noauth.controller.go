package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"math"
	"math/rand"
	"merchSearch/model"
	"merchSearch/service"
	"merchSearch/utils"
	"time"
)

const (
	MalRoot                    = "https://api.myanimelist.net/v2/"
	AnimeListCompletedEndpoint = "users/%s/animelist?status=completed&limit=%d"
)

func GetAnime(malClient *loggly.ClientType, endPoint string, user string, limit int) []byte {
	offset := getRandOffset(malClient, endPoint, user, limit)
	response := utils.HttpGet(MalRoot+F(endPoint, user, limit)+F("&offset=%d", offset), malClient)
	return service.GetBytes(response, malClient)
}

func getRandOffset(malClient *loggly.ClientType, endPoint string, user string, limit int) int {
	numOfCompletedAnime := GetCompletedCount(malClient, endPoint, user)
	numOfPages := int(
		math.Ceil(
			float64(numOfCompletedAnime) / float64(limit),
		))
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(numOfPages)
}

func GetCompletedCount(malClient *loggly.ClientType, endPoint string, user string) int {
	var numOfCompletedAnime int
	var userCompletedResponse model.MalUserListResponse
	hasNext := true
	for hasNext {
		response := utils.HttpGet(MalRoot+fmt.Sprintf(endPoint, user, 1000), malClient)
		animeBodyBytes := service.GetBytes(response, malClient)
		animeUnmarshalErr := json.Unmarshal(animeBodyBytes, &userCompletedResponse)
		service.UnmarshalError(animeUnmarshalErr, malClient)
		numOfCompletedAnime += len(userCompletedResponse.Data)
		if userCompletedResponse.Paging.Next == "" {
			hasNext = false
		}
	}
	return numOfCompletedAnime
}
