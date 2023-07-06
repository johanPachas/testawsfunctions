package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Event struct {
	Id   string `json:"Identificador"`
	Name string `json:"Nombre"`
}
type Student struct {
	Identificador string `json:"Identificador"`
	Nombre        string `json:"Nombre"`
	Apellidos     string `json:"Apellidos"`
	Edad          int    `json:"Edad"`
}

func handler(ctc context.Context, event Event) (string, error) {
	tableName := os.Getenv("TableName")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		panic(err)
	}
	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Identificador": {
				S: aws.String(event.Id),
			},
			"Nombre": {
				S: aws.String(event.Name),
			},
		},
	}
	_, error := svc.DeleteItem(input)
	if err != nil {
		panic(fmt.Sprintf("Got error calling DeleteItem: %s", error))

	}
	return "SUCCESS", nil
}
func main() {
	lambda.Start(handler)
}
