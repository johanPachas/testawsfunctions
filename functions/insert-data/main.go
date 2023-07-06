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

type Student struct {
	Identificador string `json:"Identificador"`
	Nombre        string `json:"Nombre"`
	Apellidos     string `json:"Apellidos"`
	Edad          int    `json:"Edad"`
}

type CreateStudent struct {
	Identificador string `json:"Identificador"`
	Nombre        string `json:"Nombre"`
}

func handler(ctx context.Context, createStudent CreateStudent) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		panic(err)
	}
	svc := dynamodb.New(sess)
	inputPut := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TableName")),
		Item: map[string]*dynamodb.AttributeValue{
			"Identificador": {
				S: aws.String(createStudent.Identificador),
			},
			"Nombre": {
				S: aws.String(createStudent.Nombre),
			},
			"Apellidos": {
				S: aws.String("apellido"),
			},
			"Edad": {
				N: aws.String("20"),
			},
			"Booleano": {
				BOOL: aws.Bool(true),
			},
		},
	}
	_, err = svc.PutItem(inputPut)
	if err != nil {
		fmt.Println("Error llamando putItem")
	}

	return "Works", nil
}
func main() {
	lambda.Start(handler)
}
