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

type Event struct {
	Id     string `json:"id"`
	Nombre string `json:"nombre"`
}

type Student struct {
	Identificador string `json:"identificador"`
	Nombre        string `json:"nombre"`
	Apellidos     string `json:"apellidos"`
	Edad          int    `json:"edad"`
}
type Request struct {
	Identificador string `json:"identificador"`
	Nombre        string `json:"nombre"`
}
type Report struct {
	Student  Student   `json:"student"`
	Requests []Request `json:"requests"`
}

func handler(ctx context.Context, event Event) ([]Report, error) {

	tableName := os.Getenv("tableName")
	index := "Identificador-Nombre-index"
	var students []Student
	var reportStudent Report
	var reports []Report

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("region"))},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect session, %v", err))
	}
	svc := dynamodb.New(sess)

	errQ := svc.QueryPages(&dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("Identificador = :id AND Nombre = :name"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(event.Id),
			},
			":name": {
				S: aws.String(event.Nombre),
			},
		},
	}, func(resultQuery *dynamodb.QueryOutput, last bool) bool {
		items := []Student{}
		err := dynamodbattribute.UnmarshalListOfMaps(resultQuery.Items, &items)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshall dynamodb scan items, %v", err))
		}

		students = append(students, items...)

		return true
	})
	if errQ != nil {
		panic(fmt.Sprintf("Got error calling Query: %s", errQ))
	}

	for _, stud := range students {
		requests := []Request{}
		reportStudent.Student = stud

		errQ := svc.QueryPages(&dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			IndexName:              aws.String(index),
			KeyConditionExpression: aws.String("Identificador = :Identificador AND Nombre = :Nombre"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":Identificador": {
					S: aws.String(event.Id),
				},
				":Nombre": {
					S: aws.String(event.Nombre),
				},
				":Apellidos": {
					S: aws.String(stud.Apellidos),
				},
			},
			FilterExpression: aws.String("Apellidos = :Apellidos"),
		}, func(resultQuery *dynamodb.QueryOutput, last bool) bool {
			items := []Request{}
			err := dynamodbattribute.UnmarshalListOfMaps(resultQuery.Items, &items)
			if err != nil {
				panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
			}

			requests = append(requests, items...)

			return true // keep paging
		})

		if errQ != nil {
			panic(fmt.Sprintf("Got error calling Query2: %s", errQ))
		}
		fmt.Println("requests : ", requests)
		reportStudent.Requests = requests
		reports = append(reports, reportStudent)
	}
	return reports, nil
}
func main() {
	lambda.Start(handler)
}
