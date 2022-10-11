package service

import (
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"merchSearch/controller"
	"merchSearch/model"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// to be replaced with sending token info to a db
func SetTokenEnvs(tokenModel model.TokenInfo) {
	// Instantiate Client
	tokenGeneratorTag := fmt.Sprintf("%ssetTokenEnvs()", controller.TagRoot)
	tokenGenClient := loggly.New(tokenGeneratorTag)

	setBearerErr := os.Setenv("EBAY_BEARER_TOKEN", fmt.Sprintf("Bearer %s", tokenModel.AccessToken))
	if setBearerErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay Bearer token environment Variable")
		ClientErrorCheck(clientErr)
	}
	// adds the time till expiration to the current time and converts it to a string in preparation of env/db storage
	dateOfTokenExpiration := strconv.FormatInt(
		time.Now().Add(time.Second*time.Duration(tokenModel.ExpiresIn)).Unix(),
		10,
	)
	setExpiryErr := os.Setenv("EBAY_BEARER_TOKEN_EXPIRATION", dateOfTokenExpiration)
	if setExpiryErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay token timer token environment Variable")
		ClientErrorCheck(clientErr)
	}
	setTokenTypeErr := os.Setenv("EBAY_TOKEN_TYPE", tokenModel.TokenType)
	if setTokenTypeErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay token type environment Variable")
		ClientErrorCheck(clientErr)
	}
}

func SearchEbay(ebayClient *loggly.ClientType, malRes model.MalUserListResponse, sort string) []model.ItemSummaries {
	var ebayResponseModel model.EbaySearchResponse
	var itemSummaries []model.ItemSummaries
	for _, data := range malRes.Data {
		title := strings.ReplaceAll(data.Node.Title, " ", "+")
		ebayBytes := RequestEbayBytes(ebayClient, title, sort)
		ebayUnmarshalErr := json.Unmarshal(ebayBytes, &ebayResponseModel)
		UnmarshalError(ebayUnmarshalErr, ebayClient)
		itemSummaries = append(itemSummaries, ebayResponseModel.ItemSummaries...)
		time.Sleep(1 * time.Second)
	}
	return itemSummaries
}

func RequestEbayBytes(ebayClient *loggly.ClientType, title string, sort string) []byte {

	response, httpErr := controller.Search(title, 0, 3, sort)
	HttpErrorCheck(httpErr, ebayClient)

	return GetBytes(response, ebayClient)
}

func GenerateNewTokenIfNotExist() {
	var wg sync.WaitGroup

	if TokenExpiredOrDoesntExist() {
		fmt.Println("authToken expired. generating new token. . .")
		wg.Add(1)
		controller.TokenGenerator(&wg)
		wg.Wait()
		fmt.Println("authToken generated")
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
