package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var svc = dynamodb.New(sess)

type Item struct {
	Id       string
	Phone    string
	UserType string
}

func getItem(userId string) Item {
	tableName := "users"

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(userId),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		return Item{}
	}

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return item
}

type SearchEvent struct {
	UserId string `json:"userId"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func HandleRequest(ctx context.Context, inputEvent SearchEvent) (Response, error) {

	res, err := json.Marshal(getItem(inputEvent.UserId))

	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       "ERROR",
		}, nil
	}

	return Response{
		StatusCode: 200,
		Body:       string(res),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
