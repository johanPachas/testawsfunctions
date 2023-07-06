package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Student struct {
	Identificador string `json:"Identificador"`
	Nombre        string `json:"Nombre"`
	Apellidos     string `json:"Apellidos"`
	Edad          string `json:"Edad"`
}

type Event struct {
	Identificador string `json:"Identificador"`
	Nombre        string `json:"Nombre"`
	Apellidos     string `json:"Apellidos"`
	Edad          string `json:"Edad"`
	Booleano      bool   `json:"Booleano"`
}

func handler(ctx context.Context, event Event) (string, error) {
	tableName := os.Getenv("TableName")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		return "ERROR: ", err
	}

	svc := dynamodb.New(sess)
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Identificador": {
				S: aws.String(event.Identificador),
			},
			"Nombre": {
				S: aws.String(event.Nombre),
			},
		},
		UpdateExpression: aws.String("set Edad = :edad, Booleano = :Bol, Apellidos = :apellido"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":edad": {
				N: aws.String(event.Edad),
			},
			":Bol": {
				BOOL: aws.Bool(false),
			},
			":apellido": {
				S: aws.String("MULTIUPDATEDATOS"),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	result, err := svc.UpdateItem(input)
	if err != nil {
		panic(fmt.Sprintf("failed to Dynamodb Update Items, %v", err))
	}

	user := Student{}

	dynamodbattribute.UnmarshalMap(result.Attributes, &user)
	return "WORKS", nil
}
func main() {
	lambda.Start(handler)
}
