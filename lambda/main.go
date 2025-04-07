package main

import (
	"fmt"

	"lambda-func/app"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// Take in payload and process it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	return fmt.Sprintf("Successfuly called by - %s\n", event.Username), nil 
}

func main() {
	myApp := app.NewApp()

	lambda.Start(myApp.ApiHandler.RegisterUserHandler)
}