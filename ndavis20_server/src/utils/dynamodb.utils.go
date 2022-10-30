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
	"os"
	"regexp"
)

// createDynamoDBClient: Creates and returns a new dynamodb session client
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

/*
FindItem : Uses getItem() to find and return a go object representation of a single item from
the specified dynamodb table.
*/
func FindItem(tableName string, userId string, itemId string) *model.UserEntry {
	findItemClient := InstantiateClient("service.ddbUtils.findItem")
	dynamodbItem, err := getItem(tableName, userId, itemId)
	if err != nil {
		log.Println(findItemClient.EchoSend("error", err.Error()))
	}
	goObjItem := unmarshallDynamodbObj(dynamodbItem)
	return goObjItem
}

/*
GetAllItems : uses dynamodb.DynamoDB.Scan() to return all items from the specified table.
*/
func GetAllItems(tableName string) []*model.UserEntry {
	scannedItems := scan(tableName)
	goObjItem := unmarshallDynamodbObjMulti(scannedItems)

	return goObjItem
}

/*
GetLiveItemCount : Returns the number of items in the specified dynamodb table

WARNING:

Utilizes dynamodb.DynamoDB.Scan() which retrieves every item from the specified dynamodb table.
If you do not need a live item count use GetItemCount which is guaranteed to be updated every 6 hours.
*/
func GetLiveItemCount(tableName string) *int64 {
	LiveItemCountClient := InstantiateClient("service.ddbUtils.scan")
	dynamodbClient := createDynamoDBClient()

	items, err := dynamodbClient.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Println(LiveItemCountClient.EchoSend("error", err.Error()))
	}
	return items.Count
}

/*
GetItemCount : Returns the number of items in the specified dynamodb table

WARNING:

Utilizes dynamodb.DynamoDB.DescribeTable() which is only updated every six hours. If you need a live
item count use GetLiveItemCount
*/
func GetItemCount(tableName string) *int64 {
	ItemCountClient := InstantiateClient("service.ddbUtils.scan")
	dynamodbClient := createDynamoDBClient()

	tableDescription, err := dynamodbClient.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Println(ItemCountClient.EchoSend("error", err.Error()))
	}

	return tableDescription.Table.ItemCount
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

// getItem: this and that
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

func CheckItemIdReg(itemId string, writer http.ResponseWriter) bool {
	itemIdReg, _ := regexp.Compile("^((v[0-9])\\|([0-9]{1,12})\\|([0-9]{1,12}))$")
	if !itemIdReg.MatchString(itemId) {
		http.Error(
			writer,
			"Malformed request, check your itemId",
			http.StatusBadRequest,
		)
		return false
	}
	return true
}

func CheckUserIdReg(userId string, writer http.ResponseWriter) bool {
	userIdReg, _ := regexp.Compile("[a-zA-z0-9_-]{2,16}")

	if !userIdReg.MatchString(userId) {
		http.Error(
			writer,
			"Malformed request, check your userId",
			http.StatusBadRequest,
		)
		return false
	}
	return true
}

func CheckTableRegex(table string, writer http.ResponseWriter) bool {
	tableReg, _ := regexp.Compile(os.Getenv("DYNAMO_DB_TABLE_NAME"))

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
