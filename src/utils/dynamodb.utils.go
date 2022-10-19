package utils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"merchSearch/model"
	"sync"
)

type DynamoDBRepo struct {
	TableName string
}

func CreateDynamoDBClient() *dynamodb.DynamoDB {
	session := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

	// Create an Amazon S3 service client
	return dynamodb.New(session)
}

func (repo *DynamoDBRepo) Save(post *model.UserEntry, wg *sync.WaitGroup) (*model.UserEntry, error) {

	// creation of a new DynamoDb client
	dynamodbClient := CreateDynamoDBClient()

	attributeVal, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}
	item := &dynamodb.PutItemInput{
		Item:      attributeVal,
		TableName: aws.String(repo.TableName),
	}

	_, err = dynamodbClient.PutItem(item)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	wg.Done()

	return post, nil
}

func (repo *DynamoDBRepo) SaveMultiple(posts []model.UserEntry) ([]*model.UserEntry, error) {
	var saved []*model.UserEntry
	for _, post := range posts {
		var wg sync.WaitGroup
		wg.Add(1)
		currentPost, err := repo.Save(&post, &wg)
		wg.Wait()
		if err != nil {
			return nil, err
		}
		saved = append(saved, currentPost)
	}
	return saved, nil
}

func EbayToUserEntry(ebayItems *model.ItemSummaries, userId string) model.UserEntry {
	return model.UserEntry{
		UserId:        userId,
		ItemId:        ebayItems.ItemId,
		ItemSummaries: *ebayItems,
	}
}

func EbayListToUserEntryList(ebayItems *[]model.ItemSummaries, userId string) []model.UserEntry {
	var userEntry []model.UserEntry
	for _, item := range *ebayItems {
		userEntry = append(userEntry, EbayToUserEntry(&item, userId))
	}
	return userEntry
}
