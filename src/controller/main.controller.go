package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"io"
	"log"
	"merchSearch/model"
	"net/http"
	"runtime"
)

const (
	MalRoot                    = "https://api.myanimelist.net/v2/"
	AnimeListCompletedEndpoint = "users/%s/animelist?status=completed"
	//MangaListCompletedEndpoint = "users/%s/mangalist?status=completed"
)

func Run() {
	LoadDotEnv(".env")

	// Instantiate Clients
	animeTag := "anime_search"
	//ebayTag := "ebay_search"
	animeClient := loggly.New(animeTag)
	//ebayClient := loggly.New(ebayTag)

	var userCompletedResponse model.MalUserListResponse

	bodyBytes := GetAnimeListBytes(animeClient)

	// Convert JSON into go type definition
	unmarshalErr := json.Unmarshal(bodyBytes, &userCompletedResponse)
	unmarshalError(unmarshalErr, animeClient)

}

func GetAnimeListBytes(client *loggly.ClientType) []byte {
	response, httpErr := GetCompleted(AnimeListCompletedEndpoint, "Curiossity")
	httpErrorCheck(httpErr, client)

	bodyBytes, readBytesErr := io.ReadAll(response.Body)
	readBytesErrorCheck(readBytesErr, client)

	deferResponseBodyClose(response)

	// send loggly msg
	clientErr := client.EchoSend(
		"info", fmt.Sprintf("statusCode: %v\nresponseSize: %v", response.StatusCode, len(bodyBytes)),
	)
	clientErrorCheck(clientErr)

	return bodyBytes
}

func httpErrorCheck(httpErr error, client *loggly.ClientType) {
	if httpErr != nil {
		clientErr := client.EchoSend("error", fmt.Sprintf("Was not able to pull users Completed list\nerror: %s", httpErr))
		if clientErr != nil {
			log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
		}
	}
}

func deferResponseBodyClose(response *http.Response) {
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

func readBytesErrorCheck(readBytesErr error, client *loggly.ClientType) {
	// Get the name of the file that called this function
	_, file, _, _ := runtime.Caller(1)

	// Handle error by sending a log to loggly or log.Fatalln if not possible
	if readBytesErr != nil {
		clientErr := client.EchoSend(
			"error", fmt.Sprintf("Location: %v\n io.Read all could not readBytes \nerror: %s", file, readBytesErr),
		)
		if clientErr != nil {
			log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
		}
	}
}

func clientErrorCheck(clientErr error) {
	if clientErr != nil {
		log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
	}
}

func unmarshalError(unmarshalErr error, client *loggly.ClientType) {
	if unmarshalErr != nil {
		clientErr := client.EchoSend(
			"error", fmt.Sprintf("Was not able to unmarshal the bytes of the response body\nerror: %s", unmarshalErr),
		)
		if clientErr != nil {
			log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
		}
	}
}

func LoadDotEnv(filenames ...string) {
	err1 := godotenv.Load(filenames...)
	if err1 != nil {
		log.Fatalln("Error loading .env file")
	}
}
