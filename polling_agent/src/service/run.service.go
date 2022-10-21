package service

import (
	"encoding/json"
	"fmt"
	"os"
	"polling_agent/handler"
	"polling_agent/model"
	"polling_agent/utils"
	"time"
)

var F = fmt.Sprintf
var userName = "curiossity"

func Run() {

	//if running on a local machine use the .env file
	utils.LoadDotEnv(".env")
	// Instantiate Clients
	animeClient := utils.InstantiateClient("anime_search")
	ebayClient := utils.InstantiateClient("ebay_search")
	dynamodbTable := utils.DynamoDBRepo{TableName: os.Getenv("DYNAMO_DB_TABLE_NAME")}
	// compatible
	for {
		// will be replaced with db check if the users current session still has an active oauth token in the db
		// which will send the users to the ebay sign it page. Upon success, they will be sent to the landing page
		// which will resume the program.
		GenerateNewTokenIfNotExist()

		var userCompletedResponse model.MalUserListResponse

		animeBodyBytes := GetAnime(animeClient, AnimeListCompletedEndpoint, userName, 5)
		// Convert JSON into go type definition
		animeUnmarshalErr := json.Unmarshal(animeBodyBytes, &userCompletedResponse)
		handler.UnmarshalError(animeUnmarshalErr, animeClient, string(animeBodyBytes))
		ebayItems := SearchItemsOnEbay(ebayClient, userCompletedResponse, "Ascending")
		DynamoItems := utils.EbayListToUserEntryList(&ebayItems, userName)
		_, _ = dynamodbTable.SaveMultiple(DynamoItems)

		fmt.Println(ebayItems)
		fmt.Println("---END---")
		time.Sleep(10 * time.Second)
	}
}
