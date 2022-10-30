package utils

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"ndavis20_server/model"
	"net/http"
	"regexp"
)

func createDynamoDBClient() *dynamodb.DynamoDB {
	dynamodbLoggly := InstantiateClient("server.DynamoDB.create.client")
	ddbSession := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
	dynamodbClient := dynamodb.New(ddbSession)
	err := dynamodbLoggly.EchoSend("info", "client created")
	if err != nil {
		log.Println(err)
	}

	// Create an Amazon S3 service client
	return dynamodbClient
}

func FindItem(tableName string, userId string, itemId string) *model.UserEntry {
	findItemClient := InstantiateClient("service.ddbUtils.findItem")
	dynamodbItem, err := getItem(tableName, userId, itemId)
	if err != nil {
		log.Println(findItemClient.EchoSend("error", err.Error()))
	}
	goObjItem := unmarshallDynamodbObj(dynamodbItem)
	return goObjItem
}

func GetAllItems(tableName string) []*model.UserEntry {
	scannedItems := scan(tableName)
	goObjItem := unmarshallDynamodbObjMulti(scannedItems)

	return goObjItem
}

func scan(tableName string) *dynamodb.ScanOutput {
	scanItemsClient := InstantiateClient("service.ddbUtils.scan")
	dynamodbClient := createDynamoDBClient()

	items, err := dynamodbClient.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Println(scanItemsClient.EchoSend("error", err.Error()))
	}
	return items
}

func getItem(tableName string, userId string, itemId string) (*dynamodb.GetItemOutput, error) {
	dynamodbClient := createDynamoDBClient()
	result, err := dynamodbClient.GetItem(&dynamodb.GetItemInput{

		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(userId),
			},
			"item_id": {
				S: aws.String(itemId),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		msg := "Could not find '" + userId + "'"
		return nil, errors.New(msg)
	}

	return result, nil
}

func unmarshallDynamodbObj(result *dynamodb.GetItemOutput) *model.UserEntry {
	var item *model.UserEntry

	err := dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	return item
}

func unmarshallDynamodbObjMulti(result *dynamodb.ScanOutput) []*model.UserEntry {
	var items []*model.UserEntry
	for count, i := range result.Items {
		items = append(items, &model.UserEntry{})
		err := dynamodbattribute.UnmarshalMap(i, &items[count])
		if err != nil {
			log.Println(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
	}
	return items
}

func CheckAgainstRegSearch(table string, itemId string, userId string, writer http.ResponseWriter) bool {
	tableReg, _ := regexp.Compile("ndavis20-merchSearch")
	itemIdReg, _ := regexp.Compile("^((v[0-9])\\|([0-9]{1,12})\\|([0-9]{1,12}))$")
	userIdReg, _ := regexp.Compile("[a-zA-z0-9_-]{2,16}")

	if !tableReg.MatchString(table) {
		http.Error(
			writer,
			"You either do not have permission to access this table "+
				"or the table does not exist",
			http.StatusUnauthorized,
		)
		return false
	}
	if !itemIdReg.MatchString(itemId) || !userIdReg.MatchString(userId) {
		http.Error(
			writer,
			"Malformed request",
			http.StatusBadRequest,
		)
		return false
	}
	return true
}

func CheckTableRegex(table string, writer http.ResponseWriter) bool {
	tableReg, _ := regexp.Compile("ndavis20-merchSearch")

	if !tableReg.MatchString(table) {
		http.Error(
			writer,
			"You either do not have permission to access this table "+
				"or the table does not exist",
			http.StatusUnauthorized,
		)
		return false
	}
	return true
}
