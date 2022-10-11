package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"io"
	"log"
	"merchSearch/model"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Run() {

	var wg sync.WaitGroup
	//if running on a local machine use the .env file
	if FileExist(".env") {
		LoadDotEnv(".env")
	}
	// will be replaced with db check if the users current session still has an active oauth token in the db
	// which will send the users to the ebay sign it page. Upon success, they will be sent to the landing page
	// which will resume the program.
	for {

		if TokenExpiredOrDoesntExist() {
			fmt.Println("authToken expired. generating new token. . .")
			wg.Add(1)
			tokenGenerator(&wg)
			wg.Wait()
			fmt.Println("authToken generated")
		}
		// Instantiate Clients
		animeTag := "anime_search"
		ebayTag := "ebay_search"
		animeClient := loggly.New(animeTag)
		ebayClient := loggly.New(ebayTag)

		var userCompletedResponse model.MalUserListResponse
		animeBodyBytes := RequestMalBytes(animeClient, AnimeListCompletedEndpoint)
		// Convert JSON into go type definition
		animeUnmarshalErr := json.Unmarshal(animeBodyBytes, &userCompletedResponse)
		unmarshalError(animeUnmarshalErr, animeClient)
		ebayItems := searchEbay(ebayClient, userCompletedResponse, "Ascending")
		fmt.Println(ebayItems)
		fmt.Println("---END---")
		time.Sleep(2 * time.Minute)
	}
}

func RequestEbayBytes(ebayClient *loggly.ClientType, title string, sort string) []byte {

	response, httpErr := search(title, 0, 3, sort)
	httpErrorCheck(httpErr, ebayClient)

	return GetBytes(response, ebayClient)
}

func RequestMalBytes(malClient *loggly.ClientType, endpoint string) []byte {
	response, httpErr := GetCompleted(endpoint, "Curiossity")
	httpErrorCheck(httpErr, malClient)

	return GetBytes(response, malClient)
}

func GetBytes(response *http.Response, client *loggly.ClientType) []byte {
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

func searchEbay(ebayClient *loggly.ClientType, malRes model.MalUserListResponse, sort string) []model.ItemSummaries {
	var ebayResponseModel model.EbaySearchResponse
	var itemSummaries []model.ItemSummaries
	for _, data := range malRes.Data {
		title := strings.ReplaceAll(data.Node.Title, " ", "+")
		ebayBytes := RequestEbayBytes(ebayClient, title, sort)
		ebayUnmarshalErr := json.Unmarshal(ebayBytes, &ebayResponseModel)
		unmarshalError(ebayUnmarshalErr, ebayClient)
		itemSummaries = append(itemSummaries, ebayResponseModel.ItemSummaries...)
		time.Sleep(1 * time.Second)
	}
	return itemSummaries
}

func httpErrorCheck(httpErr error, client *loggly.ClientType) {
	if httpErr != nil {
		clientErr := client.EchoSend("error", fmt.Sprintf("Was not able to connect to endpoint list\nerror: %s", httpErr))
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
			"error", fmt.Sprintf(
				"Was not able to unmarshal the bytes of the response body\nerror: %s\n"+
					"",
				unmarshalErr),
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

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func TokenExpiredOrDoesntExist() bool {
	expirationTime, _ := strconv.ParseInt(os.Getenv("EBAY_BEARER_TOKEN_EXPIRATION"), 10, 64)
	EbayTokenExpired := expirationTime < time.Now().Unix()
	TokenNotExist := !(first(os.LookupEnv("EBAY_BEARER_TOKEN")) &&
		first(os.LookupEnv("EBAY_BEARER_TOKEN_EXPIRATION")) &&
		first(os.LookupEnv("EBAY_TOKEN_TYPE")))
	return EbayTokenExpired || TokenNotExist
}

func first(key string, isExist bool) bool {
	return isExist
}
