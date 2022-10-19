// This package is used to test the Loggly package
package main

import (
	"merchSearch/model"
	"merchSearch/utils"
	"os"
	"sync"
)

func main() {
	utils.LoadDotEnv(".env")
	item := utils.DynamoDBRepo{
		TableName: os.Getenv("DYNAMO_DB_TABLE_NAME"),
	}
	md := model.UserEntry{
		UserId: "Curiossity",
		ItemId: "1234565",
		ItemSummaries: model.ItemSummaries{
			ItemId:           "123456",
			Title:            "0987654",
			Categories:       nil,
			Image:            model.Image{},
			Price:            model.Price{},
			Seller:           model.Seller{},
			Condition:        "new",
			ThumbnailImages:  nil,
			ItemWebURL:       "url",
			AdditionalImages: nil,
			AdultOnly:        false,
		},
	}
	var wg sync.WaitGroup
	wg.Add(1)
	_, err := item.Save(&md, &wg)
	wg.Wait()

	if err != nil {
		print(err)
	}
}
