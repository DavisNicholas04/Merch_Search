package utils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"merchSearch/model"
)

type DynamoDBRepo struct {
	tableName string
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

func (repo *DynamoDBRepo) save(post *model.UserEntry) (*model.UserEntry, error) {

	// creation of a new DynamoDb client
	dynamodbClient := CreateDynamoDBClient()

	attributeVal, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		return nil, err
	}
	item := &dynamodb.PutItemInput{
		Item:      attributeVal,
		TableName: aws.String(repo.tableName),
	}

	_, err = dynamodbClient.PutItem(item)
	if err != nil {
		return nil, err
	}
	return post, nil
}
