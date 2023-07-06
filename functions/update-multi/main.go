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
	Direccion     string `json:"direccion"`
	DNI           string `json:"dni"`
}

type Event struct {
	Student []Student `json:"students"`
}
type Report struct {
	Student Student `json:"students"`
}

func handler(ctx context.Context, event Event) (string, error) {
	var students []Student
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		panic(err)
	}

	i := 0
	svc := dynamodb.New(sess)

	students = append(students, event.Student...)
	var Direct string
	var DNI string
	for _, stud := range students {
		fmt.Println("Direccion1:", students[0].Direccion)
		fmt.Println("Direccion2:", students[1].Direccion)
		fmt.Println("I:", i)
		if i == 0 {
			Direct = students[1].Direccion
			DNI = students[1].DNI
		} else {
			Direct = students[0].Direccion
			DNI = students[0].DNI
		}
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
					S: aws.String("UpdateApellido"),
				},
				"Edad": {
					S: aws.String(stud.Edad),
				},
				"Booleano": {
					BOOL: aws.Bool(true),
				},
				"Direccion": {
					S: aws.String(Direct),
				},
				"DNI": {
					S: aws.String(DNI),
				},
			},
		}
		i++
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
