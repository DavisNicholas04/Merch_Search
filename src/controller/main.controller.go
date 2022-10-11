package controller

import (
	"encoding/json"
	"fmt"
	"merchSearch/model"
	"merchSearch/service"
	"time"
)

func Run() {

	//if running on a local machine use the .env file
	service.LoadDotEnv(".env")

	for {
		// will be replaced with db check if the users current session still has an active oauth token in the db
		// which will send the users to the ebay sign it page. Upon success, they will be sent to the landing page
		// which will resume the program.
		service.GenerateNewTokenIfNotExist()

		// Instantiate Clients
		animeClient := service.InstantiateClient("anime_search")
		ebayClient := service.InstantiateClient("ebay_search")

		var userCompletedResponse model.MalUserListResponse
		animeBodyBytes := service.RequestMalBytes(animeClient, AnimeListCompletedEndpoint)
		// Convert JSON into go type definition
		animeUnmarshalErr := json.Unmarshal(animeBodyBytes, &userCompletedResponse)
		service.UnmarshalError(animeUnmarshalErr, animeClient)
		ebayItems := service.SearchEbay(ebayClient, userCompletedResponse, "Ascending")
		fmt.Println(ebayItems)
		fmt.Println("---END---")
		time.Sleep(2 * time.Minute)
	}
}
