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
	Identificador string `json:"identificador"`
	Nombre        string `json:"nombre"`
	Apellidos     string `json:"apellidos"`
	Edad          string `json:"edad"`
	DNI           string `json:"dni"`
	Direccion     string `json:"direccion"`
	Correo        string `json:"correo"`
}

type CreateStudent struct {
	Student []Student `json:"student"`
}

func handler(ctx context.Context, createStudent CreateStudent) (string, error) {
	var students []Student
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		panic(err)
	}
	svc := dynamodb.New(sess)
	students = append(students, createStudent.Student...)
	fmt.Println(students)
	for _, stud := range students {
		inputPut := &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("TableName")),
			Item: map[string]*dynamodb.AttributeValue{
				"Identificador": {
					S: aws.String(stud.Identificador),
				},
				"Nombre": {
					S: aws.String(stud.Nombre),
				},
				"Apellidos": {
					S: aws.String("apellido"),
				},
				"Edad": {
					S: aws.String(stud.Edad),
				},
				"Booleano": {
					BOOL: aws.Bool(true),
				},
				"Direccion": {
					S: aws.String("direccion1"),
				},
				"DNI": {
					S: aws.String("74001292"),
				},
			},
		}
		fmt.Println(inputPut)
		_, err = svc.PutItem(inputPut)
		if err != nil {
			fmt.Println("Error llamando putItem")
		}
	}
	return "Works", nil
}
func main() {
	lambda.Start(handler)
}
